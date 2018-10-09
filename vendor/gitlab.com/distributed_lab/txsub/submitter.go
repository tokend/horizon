package txsub

import (
	"fmt"
	"gitlab.com/distributed_lab/corer"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"golang.org/x/net/context"
	"time"
)

// NewDefaultSubmitter returns a new, simple Submitter implementation
// that submits directly to the stellar-core at `url` using the http client
// `h`.
func NewDefaultSubmitter(connector corer.Connector) Submitter {
	return &submitter{
		connector: connector,
	}
}

// submitter is the default implementation for the Submitter interface.  It
// submits directly to the configured stellar-core instance using the
// configured http client.
type submitter struct {
	connector corer.Connector
}

// Submit sends the provided envelope to stellar-core and parses the response into
// a SubmissionResult
func (sub *submitter) Submit(ctx context.Context, env string) (result SubmissionResult) {
	start := time.Now()
	defer func() { result.Duration = time.Since(start) }()

	coreResponse, err := sub.connector.SubmitTx(env)
	if err != nil {
		result.Err = errors.Wrap(err, "Failed to submit tx to core")
		return
	}

	// interpet response
	if coreResponse.Exception != "" {
		result.Err = fmt.Errorf("stellar-core exception: %s", coreResponse.Exception)
		return
	}

	switch coreResponse.Status {
	case corer.TxStatusError:
		result.Err = NewRejectedTxError(coreResponse.Error)
	case corer.TxStatusPending, corer.TxStatusDuplicate:
		//noop.  A nil Err indicates success
	default:
		result.Err = fmt.Errorf("Unrecognized stellar-core status response: %s", coreResponse.Status)
	}

	return
}
