package regources

type TransactionRelationships struct {
	LedgerEntryChanges *RelationCollection `json:"ledger_entry_changes,omitempty"`
	Operations         *RelationCollection `json:"operations,omitempty"`
	Source             *Relation           `json:"source,omitempty"`
}
