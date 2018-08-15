package reviewablerequest

import (
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/tokend/regources"
)

func PopulatePreIssuanceRequest(histRequest history.PreIssuanceRequest) (
	r *regources.PreIssuanceRequest, err error,
) {
	r = &regources.PreIssuanceRequest{}
	r.Asset = histRequest.Asset
	r.Amount = histRequest.Amount
	r.Signature = histRequest.Signature
	r.Reference = histRequest.Reference
	return
}
