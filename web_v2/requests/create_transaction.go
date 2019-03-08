package requests

import (
	"encoding/json"
	"net/http"

	"gitlab.com/distributed_lab/logan/v3/errors"

	"gitlab.com/tokend/horizon/txsub/v2"

	"gitlab.com/tokend/horizon/web_v2/ctx"
)

// GetAccountRule - represents params to be specified by user for Get account rule handler
type CreateTransaction struct {
	*base
	Env *txsub.EnvelopeInfo
}

// NewGetAccountRule returns new instance of GetAsset request
func NewCreateTransactionRequest(r *http.Request) (*CreateTransaction, error) {
	b, err := newBase(r, baseOpts{})
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(b.request.Body)
	jsonBody := make(map[string]string)
	err = decoder.Decode(&jsonBody)
	if err != nil {
		jsonBody = map[string]string{}
	}
	tx, ok := jsonBody["tx"]
	if !ok {
		return nil, errors.New("Empty tx")
	}

	info := ctx.CoreInfo(r)

	envelopeInfo, err := txsub.ExtractEnvelopeInfo(tx, info.NetworkPassphrase)
	if err != nil {
		return nil, errors.Wrap(err, "failed to extract envelope info")
	}

	return &CreateTransaction{
		base: b,
		Env:  envelopeInfo,
	}, nil
}
