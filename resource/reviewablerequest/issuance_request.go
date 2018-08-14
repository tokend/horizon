package reviewablerequest

import (
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/tokend/regources/reviewablerequest2"
)

func PopulateIssuanceRequest(histRequest history.IssuanceRequest) (
	r *reviewablerequest2.IssuanceRequest, err error,
) {
	r = &reviewablerequest2.IssuanceRequest{}
	r.Asset = histRequest.Asset
	r.Amount = histRequest.Amount
	r.Receiver = histRequest.Receiver
	r.ExternalDetails = histRequest.ExternalDetails
	return
}
