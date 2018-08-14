package reviewablerequest

import (
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/tokend/regources/reviewablerequest2"
)

func PopulatePreIssuanceRequest(histRequest history.PreIssuanceRequest) (
	r *reviewablerequest2.PreIssuanceRequest, err error,
) {
	r = &reviewablerequest2.PreIssuanceRequest{}
	r.Asset = histRequest.Asset
	r.Amount = histRequest.Amount
	r.Signature = histRequest.Signature
	r.Reference = histRequest.Reference
	return
}
