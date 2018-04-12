package reviewablerequest

import (
	"gitlab.com/swarmfund/horizon/db2/history"
)

type PreIssuanceRequest struct {
	Asset     string `json:"asset"`
	Amount    string `json:"amount"`
	Signature string `json:"signature"`
	Reference string `json:"reference"`
}

func (r *PreIssuanceRequest) Populate(histRequest history.PreIssuanceRequest) error {
	r.Asset = histRequest.Asset
	r.Amount = histRequest.Amount
	r.Signature = histRequest.Signature
	r.Reference = histRequest.Reference
	return nil
}
