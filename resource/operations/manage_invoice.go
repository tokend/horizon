package operations

type ManageInvoice struct {
	Base
	Amount          string  `json:"amount"`
	ReceiverBalance string  `json:"receiver_balance,omitempty"`
	Sender          string  `json:"sender,omitempty"`
	InvoiceID       int64   `json:"invoice_id"`
	RejectReason    *string `json:"reject_reason,omitempty"`
	Asset           string  `json:"asset"`
}
