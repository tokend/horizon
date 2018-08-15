package reviewablerequest

import (
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/tokend/regources"
	"gitlab.com/tokend/go/amount"
)

func PopulatePreIssuanceRequest(histRequest history.PreIssuanceRequest) (
	r *regources.PreIssuanceRequest, err error,
) {
	r = &regources.PreIssuanceRequest{}
	r.Asset = histRequest.Asset
	r.Amount = regources.Amount(amount.MustParse(histRequest.Amount))
	r.Signature = histRequest.Signature
	r.Reference = histRequest.Reference
	return
}
