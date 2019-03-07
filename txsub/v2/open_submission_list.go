package txsub

import (
	"regexp"
	"sync"
	"time"

	"github.com/go-errors/errors"
	"golang.org/x/net/context"
)

// NewDefaultSubmissionList returns a list that manages open submissions purely
// in memory.
func NewDefaultSubmissionList(retry time.Duration) openSubmissionList {
	return &SubmissionList{
		submissions:  map[string]*openSubmission{},
		retryTimeout: retry,
	}
}

// openSubmission tracks a slice of channels that should be emitted to when we
// know the result for the transactions with the provided hash
type openSubmission struct {
	Envelope    *EnvelopeInfo
	SubmittedAt time.Time
	Listeners   []listener
}

type SubmissionList struct {
	sync.Mutex
	submissions  map[string]*openSubmission
	retryTimeout time.Duration
}

func (s *SubmissionList) Envelope(hash string) *EnvelopeInfo {
	if os, ok := s.submissions[hash]; ok {
		return os.Envelope
	}
	return nil
}

func (s *SubmissionList) Add(ctx context.Context, envelope *EnvelopeInfo, l listener) error {
	s.Lock()
	defer s.Unlock()

	if cap(l) == 0 {
		panic("Unbuffered listener cannot be added to openSubmissionList")
	}
	regex := regexp.MustCompile("[[:xdigit:]]+")
	if !regex.MatchString(envelope.ContentHash) {
		return errors.New("Unexpected character sequence in hash: must be hexadecimal")
	}

	if len(envelope.ContentHash) != 64 {
		return errors.New("Unexpected transaction hash length: must be 64 hex characters")
	}

	os, ok := s.submissions[envelope.ContentHash]

	if !ok {
		os = &openSubmission{
			Envelope:    envelope,
			SubmittedAt: time.Now(),
			Listeners:   []listener{},
		}
		s.submissions[envelope.ContentHash] = os
	}

	os.Listeners = append(os.Listeners, l)

	return nil
}

func (s *SubmissionList) Finish(ctx context.Context, r fullResult) error {
	s.Lock()
	defer s.Unlock()

	os, ok := s.submissions[r.Hash]
	if !ok {
		return nil
	}

	for _, l := range os.Listeners {
		l <- r
		close(l)
	}

	delete(s.submissions, r.Hash)
	return nil
}

func (s *SubmissionList) Clean(ctx context.Context, maxAge time.Duration) (int, error) {
	s.Lock()
	defer s.Unlock()

	for _, os := range s.submissions {
		if time.Since(os.SubmittedAt) > maxAge {
			delete(s.submissions, os.Envelope.ContentHash)
			for _, l := range os.Listeners {
				l <- fullResult{Err: timeoutError}
				close(l)
			}
		}
	}

	return len(s.submissions), nil
}

func (s *SubmissionList) Pending(ctx context.Context) []string {
	s.Lock()
	defer s.Unlock()
	pendingHashes := make([]string, 0, len(s.submissions))

	for hash := range s.submissions {
		pendingHashes = append(pendingHashes, hash)
	}

	return pendingHashes
}

func (s *SubmissionList) ShouldRetry(ctx context.Context, hash string, t time.Time) bool {
	s.Lock()
	defer s.Unlock()
	submission := s.submissions[hash]

	return t.After(submission.SubmittedAt.Add(s.retryTimeout))
}
