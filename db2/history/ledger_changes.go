package history

// LedgerChanges is a row of data from the `history_ledgers_changes` table
type LedgerChanges struct {
	ID uint64 `json:"id"`
	TransactionID uint64 `json:"tx_id"`
	OperationID   uint64 `json:"op_id"`
	OrderNumber    int    `json:"order_number"`
	Effect        int    `json:"effect"`
	EntryType     int    `json:"entry_type"`
	Payload     int    `json:"xdr"`
}