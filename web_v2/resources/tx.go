package resources

import (
	"fmt"
	"net/http"
	"strconv"

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

func NewTxFailure(env txsub.EnvelopeInfo, err txsub.Error) error {
	meta := make(map[string]interface{})
	meta["envelope"] = env.RawBlob
	meta["result_xdr"] = err.ResultXDR()
	return &jsonapi.ErrorObject{
		Status: http.StatusText(err.Code()),
		Detail: err.Details(),
		Code:   fmt.Sprintf("%d", err.Code()),
		Title:  "Transaction submit failed",
		Meta:   &meta,
	}
}
