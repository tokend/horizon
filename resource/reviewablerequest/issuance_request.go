package reviewablerequest

import (
	"gitlab.com/swarmfund/horizon/db2/history"
)

type IssuanceRequest struct {
	Asset    string `json:"asset"`
	Amount   string `json:"amount"`
	Receiver string `json:"receiver"`
	ExternalDetails map[string]interface{} `json:"external_details"`
}

func (r *IssuanceRequest) Populate(histRequest history.IssuanceRequest) error {
	r.Asset = histRequest.Asset
	r.Amount = histRequest.Amount
	r.Receiver = histRequest.Receiver
	r.ExternalDetails = histRequest.ExternalDetails
	return nil
}
