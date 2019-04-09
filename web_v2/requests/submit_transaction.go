package requests

import (
	"encoding/json"
	"net/http"

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

// NewGetAccountRule returns new instance of GetAsset request
func NewCreateTransactionRequest(r *http.Request) (*CreateTransaction, error) {
	b, err := newBase(r, baseOpts{})
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(b.request.Body)
	var body regources.SubmitTransactionBody
	err = decoder.Decode(&body)
	if body.Tx == "" {
		return nil, errors.New("Envelope missing in the body of request")
	}
	info := ctx.CoreInfo(r)

	envelopeInfo, err := txsub.ExtractEnvelopeInfo(body.Tx, info.NetworkPassphrase)
	if err != nil {
		return nil, errors.Wrap(err, "failed to extract envelope info")
	}

	var waitForIngest bool
	if body.WaitForIngest != nil {
		waitForIngest = *body.WaitForIngest && ctx.Config(r).Ingest
	}

	return &CreateTransaction{
		base:          b,
		Env:           envelopeInfo,
		WaitForIngest: waitForIngest,
	}, nil
}
