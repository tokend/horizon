package operations

type ManageInvoice struct {
	Base
	Amount          string `json:"amount"`
	ReceiverBalance string `json:"receiver_balance,omitempty"`
	Sender          string `json:"sender,omitempty"`
	InvoiceID       uint64 `json:"invoice_id"`
	Asset           string `json:"asset"`
}
