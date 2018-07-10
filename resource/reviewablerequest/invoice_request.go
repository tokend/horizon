package reviewablerequest

import (
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/tokend/go/amount"
)

type InvoiceRequest struct {
	Amount          string                 `json:"amount"`
	ReceiverBalance string                 `json:"receiver_balance"`
	Sender          string                 `json:"sender"`
	Details         map[string]interface{} `json:"details"`
}

func (r *InvoiceRequest) Populate(histRequest history.InvoiceRequest) error {
	r.Amount = amount.StringU(histRequest.Amount)
	r.ReceiverBalance = histRequest.ReceiverBalanceID
	r.Sender = histRequest.SenderAccountID
	r.Details = histRequest.Details
	return nil
}