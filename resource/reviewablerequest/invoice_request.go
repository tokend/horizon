package reviewablerequest

import (
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/tokend/go/amount"
)

type InvoiceRequest struct {
	Amount          string  `json:"amount"`
	ReceiverBalance string  `json:"receiver_balance,omitempty"`
	ReceiverAccount string	`json:"receiver_account,omitempty"`
	Sender          string  `json:"sender,omitempty"`
	RejectReason    *string `json:"reject_reason,omitempty"`
	Asset           string  `json:"asset"`
}

func (r *InvoiceRequest) Populate(histRequest history.InvoiceRequest) error {
	r.Amount = amount.StringU(histRequest.Amount)
	r.ReceiverBalance = histRequest.ReceiverBalanceID
	r.Sender = histRequest.SenderAccountID
	r.ReceiverAccount = histRequest.ReceiverAccountID
	//r.RejectReason = histRequest.
	return nil
}