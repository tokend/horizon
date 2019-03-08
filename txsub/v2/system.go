package txsub

import (
	"context"
	"time"

	"github.com/lib/pq"

	"gitlab.com/distributed_lab/logan/v3/errors"

	"gitlab.com/distributed_lab/logan/v3"
)

// System represents a completely configured transaction submission system.
// Its methods tie together the various pieces used to reliably submit transactions
// to a core instance.
type System struct {
	Pending           openSubmissionList
	Results           resultProvider
	Submitter         Submitter
	NetworkPassphrase string
	SubmissionTimeout time.Duration
	Log               *logan.Entry
	Listener          *pq.Listener
}

// Submit submits the provided base64 encoded transaction envelope to the
// network using this submission system.
func (s *System) Submit(ctx context.Context, envelope EnvelopeInfo) (*Result, error) {
	submission := s.trySubmit(ctx, envelope)
	select {
	case result := <-submission:
		return result.unwrap()
	case <-ctx.Done():
		return fullResult{
			Err: timeoutError,
		}.unwrap()
	}
}

func (s *System) trySubmit(ctx context.Context, info EnvelopeInfo) <-chan fullResult {
	listener := make(chan fullResult, 1)

	// check for tx duplication
	result := s.Results.ResultByHash(ctx, info.ContentHash)
	//No duplication
	if result == nil {
		return s.submit(ctx, info, listener)
	}

	if result.Err != nil {
		err := errors.Wrap(result.Err,
			"failed to check for tx duplication",
			info.GetLoganFields())
		result.Err = err
		return send(listener, *result)
	}

	return send(listener, *result)
}

func (s *System) submit(ctx context.Context, info EnvelopeInfo, l chan fullResult) <-chan fullResult {
	_, err := s.Submitter.Submit(ctx, &info)
	if err != nil {
		return send(l,
			fullResult{
				Err: errors.Wrap(err, "failed to submit transaction",
					info.GetLoganFields()),
			})
	}

	err = s.Pending.Add(ctx, &info, l)
	if err != nil {
		return send(l,
			fullResult{
				Err: errors.Wrap(err, "failed to add tx to pending list",
					info.GetLoganFields()),
			})
	}

	return l
}

func (s *System) tryResubmit(ctx context.Context, hash string) error {
	if !s.Pending.ShouldRetry(ctx, hash, time.Now()) {
		return nil
	}

	env := s.Pending.Envelope(hash)
	if env == nil {
		return errors.New("trying to resubmit tx which is not in pending list")
	}
	_, err := s.Submitter.Submit(ctx, env)

	return err
}

func (s *System) tick(ctx context.Context) {

	for _, hash := range s.Pending.Pending(ctx) {
		res := s.Results.ResultByHash(ctx, hash)
		if res == nil {
			err := s.tryResubmit(ctx, hash)
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

		if res.Err != nil {
			s.Log.
				WithError(res.Err).
				WithFields(logan.F{
					"tx_hash": hash,
				}).
				Error("failed to submit tx")
			continue
		}

		s.Log.WithFields(logan.F{
			"tx_hash": hash,
		}).Debug("Transaction successfully submitted")

		if err := s.Pending.Finish(ctx, *res); err != nil {
			s.Log.
				WithError(res.Err).
				WithFields(logan.F{
					"tx_hash": hash,
				}).
				Error("failed to remove tx from pending list")
		}
	}

	_, err := s.Pending.Clean(ctx, s.SubmissionTimeout)
	if err != nil {
		s.Log.WithError(err).Error("failed to clean expired pending txs")
		return
	}
}

func (s *System) run(context context.Context) {
	ticker := time.NewTicker(2 * s.SubmissionTimeout)
	defer func() {
		if rvr := recover(); rvr != nil {
			s.Log.WithRecover(rvr).Error("txsub2 panicked")
		}
		ticker.Stop()
	}()

	for {
		s.wait(ticker)
		s.tick(context)
	}
}

func (s *System) Start(ctx context.Context) {
	go s.run(ctx)
}

func (s *System) wait(ticker *time.Ticker) {
	select {
	case <-s.Listener.Notify:
		return
	case <-ticker.C:
		return
	}
}
