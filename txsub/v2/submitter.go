package txsub

import (
	"time"

	"golang.org/x/net/context"
)

// resultProvider represents an abstract store that can lookup Result objects
// by transaction hash or by [address,sequence] pairs.  A resultProvider is
// used within the transaction submission system to decide whether a submission should
// be submitted to the backing core process, as well as looking up the status
// of each transaction in the open submission list at each tick (i.e. ledger close)
type resultProvider interface {
	// Look up a result by transaction hash. If result is not found returns nil
	ResultByHash(context.Context, string) *fullResult
}

// openSubmissionList represents the structure that tracks pending transactions
// and forwards Result structs on to listeners as they become available.
//
// NOTE:  An implementation of this interface will be called from multiple go-routines
// concurrently.
//
// NOTE:  A Listener must be a buffered channel.  A panic will trigger if you
// provide an unbuffered channel
type openSubmissionList interface {
	// Add registers the provided listener as interested in being notified when a
	// result is available for the provided transaction hash.
	Add(context.Context, *EnvelopeInfo, listener) error

	// Finish forwards the provided result on to any listeners and cleans up any
	// resources associated with the transaction that this result is for
	Finish(context.Context, fullResult) error

	// Clean removes any open submissions over the provided age.
	Clean(context.Context, time.Duration) (int, error)

	// Pending return a list of transaction hashes that have at least one
	// listener registered to them in this list.
	Pending(context.Context) []string

	//ShouldRetry checks whether or not submission should be resubmitted
	ShouldRetry(context.Context, string, time.Time) bool

	//Envelope returns EnvelopeInfo of open submission by hash
	// If there is no submission with such tx hash, returns nil
	Envelope(string) *EnvelopeInfo
}

// Submitter represents the low-level "submit a transaction to core"
// provider.
type Submitter interface {
	// Submit sends the provided transaction envelope to core
	Submit(ctx context.Context, env *EnvelopeInfo) (time.Duration, error)
}
