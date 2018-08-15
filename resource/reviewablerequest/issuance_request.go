package reviewablerequest

import (
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/tokend/regources"
)

func PopulateIssuanceRequest(histRequest history.IssuanceRequest) (
	r *regources.IssuanceRequest, err error,
) {
	r = &regources.IssuanceRequest{}
	r.Asset = histRequest.Asset
	r.Amount = histRequest.Amount
	r.Receiver = histRequest.Receiver
	r.ExternalDetails = histRequest.ExternalDetails
	return
}
