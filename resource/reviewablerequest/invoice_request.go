package reviewablerequest

import (
	"strconv"

	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/tokend/go/amount"
)

type InvoiceRequest struct {
	Amount     string                 `json:"amount"`
	Asset      string                 `json:"asset"`
	ContractID string                 `json:"contract_id,omitempty"`
	Details    map[string]interface{} `json:"details"`
}

func (r *InvoiceRequest) Populate(histRequest history.InvoiceRequest) error {
	r.Amount = amount.StringU(histRequest.Amount)
	r.Asset = histRequest.Asset
	if histRequest.ContractID != nil {
		r.ContractID = strconv.FormatInt(*histRequest.ContractID, 10)
	}
	r.Details = histRequest.Details
	return nil
}
