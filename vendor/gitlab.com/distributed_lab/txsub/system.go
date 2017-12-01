package txsub

import (
	"context"
	"sync"
	"time"

	"github.com/rcrowley/go-metrics"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

// System represents a completely configured transaction submission system.
// Its methods tie together the various pieces used to reliably submit transactions
// to a stellar-core instance.
type System struct {
	initializer sync.Once

	Pending           OpenSubmissionList
	Results           ResultProvider
	Submitter         Submitter
	NetworkPassphrase string
	SubmissionTimeout time.Duration

	Log *logan.Entry

	Metrics struct {
		// SubmissionTimer exposes timing metrics about the rate and latency of
		// submissions to stellar-core
		SubmissionTimer metrics.Timer

		// BufferedSubmissionGauge tracks the count of submissions buffered
		// behind this system's SubmissionQueue
		BufferedSubmissionsGauge metrics.Gauge

		// OpenSubmissionsGauge tracks the count of "open" submissions (i.e.
		// submissions whose transactions haven't been confirmed successful or failed
		OpenSubmissionsGauge metrics.Gauge

		// FailedSubmissionsMeter tracks the rate of failed transactions that have
		// been submitted to this process
		FailedSubmissionsMeter metrics.Meter

		// SuccessfulSubmissionsMeter tracks the rate of successful transactions that
		// have been submitted to this process
		SuccessfulSubmissionsMeter metrics.Meter
	}
}

// Submit submits the provided base64 encoded transaction envelope to the
// network using this submission system.
func (sys *System) Submit(ctx context.Context, envelope *EnvelopeInfo) Result {
	submission := sys.submit(ctx, envelope)
	select {
	case result := <-submission:
		return result
	case <-ctx.Done():
		return Result{
			Err: timeoutError,
		}
	}
}

func (sys *System) submit(ctx context.Context, info *EnvelopeInfo) <-chan Result {
	sys.Init()
	resultListener := make(chan Result, 1)

	result := sys.submitOnce(ctx, info)
	if result != nil {
		return sendResult(resultListener, *result)
	}

	err := sys.Pending.Add(ctx, info.ContentHash, resultListener)
	if err != nil {
		return sendResult(resultListener, Result{Err: errors.Wrap(err, "Failed to add tx to pending queue")})
	}

	return resultListener
}

func sendResult(resultListener chan Result, result Result) <-chan Result {
	resultListener <- result
	close(resultListener)
	return resultListener
}

// returns nil, if we are not sure about tx status
func (sys *System) submitOnce(ctx context.Context, info *EnvelopeInfo) *Result {
	// check the configured result provider for an existing result
	result := sys.Results.ResultByHash(ctx, info.ContentHash)
	if result != nil {
		return result
	}

	submissionResult := sys.submitToCore(ctx, info)
	// if submission succeeded
	if submissionResult.Err == nil {
		// notify that we are still not sure if transaction succeeded
		return nil
	}

	// if we've received error from core, that might be because tx was already applied -- double check results
	result = sys.Results.ResultByHash(ctx, info.ContentHash)
	if result == nil {
		// return the error from submission
		return &Result{Err: submissionResult.Err, EnvelopeXDR: info.RawBlob, Hash: info.ContentHash}
	}

	return result
}

// Submit submits the provided base64 encoded transaction envelope to the
// network using this submission system.
func (sys *System) submitToCore(ctx context.Context, env *EnvelopeInfo) SubmissionResult {
	// submit to stellar-core
	sr := sys.Submitter.Submit(ctx, env.RawBlob)
	sys.Metrics.SubmissionTimer.Update(sr.Duration)

	// if received or duplicate, add to the open submissions list
	if sr.Err == nil {
		sys.Metrics.SuccessfulSubmissionsMeter.Mark(1)
	} else {
		sys.Metrics.FailedSubmissionsMeter.Mark(1)
	}

	return sr
}

// Tick triggers the system to update itself with any new data available.
func (sys *System) Tick(ctx context.Context) {
	sys.Init()

	defer func() {
		if rec := recover(); rec != nil {
			err := errors.FromPanic(rec)
			sys.Log.WithStack(err).WithError(err).Error("tx sub tick failed")
		}
	}()

	for _, hash := range sys.Pending.Pending(ctx) {
		result := sys.Results.ResultByHash(ctx, hash)
		// no results for transaction -- nothing we can do with it
		if result == nil {
			continue
		}

		if result.HasInternalError() {
			// looks like we have internal error, so lets log it and try again to get tx result later
			sys.Log.WithError(result.Err).Error("Failed to submit tx")
			continue
		}

		sys.Log.WithField("hash", hash).Debug("finishing open submission")
		err := sys.Pending.Finish(ctx, *result)
		if err != nil {
			sys.Log.WithError(err).Error("Failed to remove tx from pending queue")
		}
	}

	stillOpen, err := sys.Pending.Clean(ctx, sys.SubmissionTimeout)
	if err != nil {
		sys.Log.WithError(err).Error("Failed to clean expired pending txs")
		return
	}

	sys.Metrics.OpenSubmissionsGauge.Update(int64(stillOpen))
}

// Init initializes `System` with required fields with default values
func (sys *System) Init() {
	sys.initializer.Do(func() {
		sys.Metrics.FailedSubmissionsMeter = metrics.NewMeter()
		sys.Metrics.SuccessfulSubmissionsMeter = metrics.NewMeter()
		sys.Metrics.SubmissionTimer = metrics.NewTimer()
		sys.Metrics.OpenSubmissionsGauge = metrics.NewGauge()

		if sys.SubmissionTimeout == 0 {
			sys.SubmissionTimeout = 1 * time.Minute
		}

		if sys.Log == nil {
			sys.Log = logan.New().Level(logan.ErrorLevel).WithField("service", "txsub.system")
		}
	})
}
