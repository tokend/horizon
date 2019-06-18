package txsub

import (
	"context"
	"time"

	"gitlab.com/distributed_lab/running"

	"gitlab.com/tokend/horizon/log"

	"github.com/lib/pq"

	"gitlab.com/distributed_lab/logan/v3/errors"

	"gitlab.com/distributed_lab/logan/v3"
)

// System represents a completely configured transaction submission system.
// Its methods tie together the various pieces used to reliably submit transactions
// to a core instance.
type System struct {
	List              openSubmissionList
	Results           resultProvider
	Submitter         Submitter
	NetworkPassphrase string
	SubmissionTimeout time.Duration
	Log               *log.Entry
	CoreListener      *pq.Listener
	HistoryListener   *pq.Listener
}

// Submit submits the provided base64 encoded transaction envelope to the
// network using this submission system.
func (s *System) Submit(ctx context.Context, envelope EnvelopeInfo, waitIngest bool) (*Result, error) {
	submission := s.trySubmit(ctx, envelope, waitIngest)
	select {
	case result := <-submission:
		return result.unwrap()
	case <-ctx.Done():
		return fullResult{
			Err: timeoutError,
		}.unwrap()
	}
}

func (s *System) checkForDuplicates(hash string) (*Result, error) {
	// check for tx duplication
	result, err := s.Results.FromCore(hash)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *System) trySubmit(ctx context.Context, info EnvelopeInfo, waitIngest bool) <-chan fullResult {
	listener := make(chan fullResult, 1)

	result, err := s.checkForDuplicates(info.ContentHash)
	if err != nil {
		return send(listener, fullResult{Err: errors.Wrap(err,
			"failed to check for tx duplication",
			info.GetLoganFields()),
		})
	}
	if result == nil {
		return s.submit(ctx, info, listener, waitIngest)
	}

	return send(listener, fullResult{
		Result: *result,
	})
}

func (s *System) submit(ctx context.Context, info EnvelopeInfo, l chan fullResult, waitIngest bool) <-chan fullResult {
	_, err := s.Submitter.Submit(ctx, &info)
	if err != nil {
		return send(l,
			fullResult{
				Err: errors.Wrap(err, "failed to submit transaction",
					info.GetLoganFields()),
			})
	}

	err = s.List.Add(&info, waitIngest, l)
	if err != nil {
		return send(l,
			fullResult{
				Err: errors.Wrap(err, "failed to add tx to pending list",
					info.GetLoganFields(),
				),
			})
	}

	return l
}

func (s *System) tryResubmit(ctx context.Context, hash string) error {
	if !s.List.ShouldRetry(hash, time.Now()) {
		return nil
	}

	res, err := s.checkForDuplicates(hash)
	if err != nil {
		return errors.Wrap(err,
			"failed to check for tx duplication", logan.F{
				"tx_hash": hash,
			})
	}
	if res != nil {
		return nil
	}

	env := s.List.Envelope(hash)
	if env == nil {
		return errors.New("trying to resubmit tx which is not in pending list")
	}
	_, err = s.Submitter.Submit(ctx, env)

	return err
}

func (s *System) tickCore(ctx context.Context) {
	for _, hash := range s.List.PendingCore() {
		res, err := s.Results.FromCore(hash)
		if IsInternalError(errors.Cause(err)) {
			s.List.Finish(fullResult{Result: Result{Hash: hash}, Err: err})
			continue
		}
		if err != nil {
			s.Log.
				WithError(err).
				WithFields(logan.F{
					"tx_hash": hash,
				}).
				Error("failed to get result from core")
			continue
		}

		if res == nil {
			err := s.tryResubmit(ctx, hash)
			if IsInternalError(errors.Cause(err)) {
				s.List.Finish(fullResult{Result: Result{Hash: hash}, Err: err})
				continue
			}
			if err != nil {
				s.Log.
					WithError(err).
					WithFields(logan.F{
						"tx_hash": hash,
					}).
					Error("failed to resubmit tx")
			}
			continue
		}

		s.Log.WithFields(log.F{
			"tx_hash": hash,
		}).Debug("Transaction successfully submitted")

		s.List.Finish(fullResult{Result: *res})
	}
}

func (s *System) tickHistory(ctx context.Context) {
	for _, hash := range s.List.PendingHistory() {
		res, err := s.Results.FromHistory(hash)
		if err != nil {
			s.Log.
				WithError(err).
				WithFields(logan.F{
					"tx_hash": hash,
				}).
				Error("failed to get result from history")
		}

		if res == nil {
			err := s.tryResubmit(ctx, hash)
			if IsInternalError(errors.Cause(err)) {
				s.List.Finish(fullResult{Result: Result{Hash: hash}, Err: err})
				continue
			}
			if err != nil {
				s.Log.
					WithError(err).
					WithFields(logan.F{
						"tx_hash": hash,
					}).
					Error("failed to resubmit tx")
			}
			continue
		}

		s.Log.WithFields(log.F{
			"tx_hash": hash,
		}).Debug("Transaction successfully submitted")

		s.List.Finish(fullResult{Result: *res})
	}
}

func (s *System) history(ctx context.Context) {
	defer func() {
		if rvr := recover(); rvr != nil {
			s.Log.WithRecover(rvr).Error("submitter_v2 panicked")
		}
	}()
	for {
		select {
		case <-s.HistoryListener.Notify:
			s.tickHistory(ctx)
		case <-ctx.Done():
			return
		}
	}
}

func (s *System) core(ctx context.Context) {
	defer func() {
		if rvr := recover(); rvr != nil {
			s.Log.WithRecover(rvr).Error("submitter_v2 panicked")
		}
	}()

	for {
		select {
		case <-s.CoreListener.Notify:
			s.tickCore(ctx)
		case <-ctx.Done():
			return
		}
	}
}

func (s *System) cleaner(ctx context.Context) {
	defer func() {
		if rvr := recover(); rvr != nil {
			s.Log.WithRecover(rvr).Error("submitter_v2 panicked")
		}
	}()

	running.WithBackOff(ctx, s.Log, "submitter_v2_cleaner", func(ctx context.Context) error {
		s.List.Clean(s.SubmissionTimeout)
		return nil
	}, s.SubmissionTimeout, time.Second, 2*s.SubmissionTimeout)
}

func (s *System) Start(ctx context.Context) {
	go s.cleaner(ctx)
	go s.core(ctx)
	go s.history(ctx)
}
