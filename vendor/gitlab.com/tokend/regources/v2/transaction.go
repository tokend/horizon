package regources

import "time"

type TransactionsResponse struct {
	Links    *Links        `json:"links"`
	Data     []Transaction `json:"data"`
	Included Included      `json:"included"`
}

type Transaction struct {
	Key
	Attributes    TransactionAttrs     `json:"attributes"`
	Relationships TransactionRelations `json:"relationships"`
}

type TransactionAttrs struct {
	Hash           string    `json:"hash"`
	LedgerSequence int32     `json:"ledger_sequence"`
	CreatedAt      time.Time `json:"created_at"`
	EnvelopeXdr    string    `json:"envelope_xdr"`
	ResultXdr      string    `json:"result_xdr"`
	ResultMetaXdr  string    `json:"result_meta_xdr"`
}

type TransactionRelations struct {
	Source             *Relation           `json:"source"`
	Operations         *RelationCollection `json:"operations"`
	LedgerEntryChanges *RelationCollection `json:"ledger_entry_changes"`
}
