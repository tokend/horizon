package requests

import (
	"encoding/json"
	"net/http"

	"github.com/go-ozzo/ozzo-validation"

	"gitlab.com/tokend/regources/generated"

	"gitlab.com/distributed_lab/logan/v3/errors"

	"gitlab.com/tokend/horizon/txsub/v2"

	"gitlab.com/tokend/horizon/web_v2/ctx"
)

// GetAccountRule - represents params to be specified by user for Get account rule handler
type CreateTransaction struct {
	*base
	Env           *txsub.EnvelopeInfo
	WaitForIngest bool
}

// NewCreateTransactionRequest returns new instance of NewCreateTransation request
func NewCreateTransactionRequest(r *http.Request) (*CreateTransaction, error) {
	b, err := newBase(r, baseOpts{})
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(b.request.Body)
	var body regources.SubmitTransactionBody
	err = decoder.Decode(&body)
	if err != nil {
		return nil, err
	}

	if body.Tx == "" {
		return nil, validation.Errors{
			"tx": errors.New("missing in the body of request"),
		}
	}
	info := ctx.CoreInfo(r)

	envelopeInfo, err := txsub.ExtractEnvelopeInfo(body.Tx, info.NetworkPassphrase)
	if err != nil {
		return nil, validation.Errors{
			"tx": errors.Wrap(err, "failed to extract envelope info"),
		}
	}

	if body.WaitForIngest && !ctx.Config(r).Ingest {
		return nil, validation.Errors{
			"wait_for_ingest": errors.New("wait for ingest is not allowed as this horizon does not perform ingest"),
		}
	}

	return &CreateTransaction{
		base:          b,
		Env:           envelopeInfo,
		WaitForIngest: body.WaitForIngest,
	}, nil
}
