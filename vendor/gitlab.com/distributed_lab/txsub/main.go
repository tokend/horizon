package txsub

import (
	"time"

	"golang.org/x/net/context"
)

// ResultProvider represents an abstract store that can lookup Result objects
// by transaction hash or by [address,sequence] pairs.  A ResultProvider is
// used within the transaction submission system to decide whether a submission should
// be submitted to the backing stellar-core process, as well as looking up the status
// of each transaction in the open submission list at each tick (i.e. ledger close)
type ResultProvider interface {
	// Look up a result by transaction hash. If result is not found returns nil
	ResultByHash(context.Context, string) *Result
}

// Listener represents some client who is interested in retrieving the result
// of a specific transaction.
type Listener chan<- Result

// OpenSubmissionList represents the structure that tracks pending transactions
// and forwards Result structs on to listeners as they become available.
//
// NOTE:  An implementation of this interface will be called from multiple go-routines
// concurrently.
//
// NOTE:  A Listener must be a buffered channel.  A panic will trigger if you
// provide an unbuffered channel
type OpenSubmissionList interface {
	// Add registers the provided listener as interested in being notified when a
	// result is available for the provided transaction hash.
	Add(context.Context, string, Listener) error

	// Finish forwards the provided result on to any listeners and cleans up any
	// resources associated with the transaction that this result is for
	Finish(context.Context, Result) error

	// Clean removes any open submissions over the provided age.
	Clean(context.Context, time.Duration) (int, error)

	// Pending return a list of transaction hashes that have at least one
	// listener registered to them in this list.
	Pending(context.Context) []string
}

// Submitter represents the low-level "submit a transaction to stellar-core"
// provider.
type Submitter interface {
	// Submit sends the provided transaction envelope to stellar-core
	Submit(ctx context.Context, env string) SubmissionResult
}

// SubmissionResult gets returned in response to a call to Submitter.Submit.
// It represents a single discrete submission of a transaction envelope to
// the stellar network.
type SubmissionResult struct {
	// Any error that occurred during the attempted submission.  A nil value
	// indicates that the submission will or already is being considered for
	// inclusion in the ledger (i.e. A successful submission).
	Err error

	// Duration records the time it took to submit a transaction
	// to stellar-core
	Duration time.Duration
}
