package resources

import (
	"fmt"
	"net/http"
	"strconv"

	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"

	"github.com/google/jsonapi"

	"gitlab.com/tokend/horizon/txsub/v2"
	regources "gitlab.com/tokend/regources/generated"
)

//NewTxKey - creates new Tx Key for specified ID
func NewTxKey(txID int64) regources.Key {
	return regources.Key{
		ID:   strconv.FormatInt(txID, 10),
		Type: regources.TRANSACTIONS,
	}
}

func NewTxFailure(env txsub.EnvelopeInfo, txSubErr txsub.Error) error {
	meta := make(map[string]interface{})
	meta["envelope"] = env.RawBlob
	meta["result_xdr"] = txSubErr.ResultXDR()
	var parsedResult xdr.TransactionResult
	err := xdr.SafeUnmarshalBase64(txSubErr.ResultXDR(), &parsedResult)
	if err != nil {
		return errors.Wrap(err, "Failed to get parse tx result")
	}

	var parsedEnvelope xdr.TransactionEnvelope
	err = xdr.SafeUnmarshalBase64(env.RawBlob, &parsedEnvelope)
	if err != nil {
		return errors.Wrap(err, "Failed to unmarshal tx envelope")
	}

	meta["parsed_envelope"] = parsedEnvelope
	meta["parsed_result"] = parsedResult
	return &jsonapi.ErrorObject{
		Status: fmt.Sprintf("%d", txSubErr.Code()),
		Detail: txSubErr.Details(),
		Title:  http.StatusText(txSubErr.Code()),
		Meta:   &meta,
	}
}
