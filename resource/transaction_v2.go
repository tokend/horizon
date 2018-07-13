package resource

import (
	"time"
	"gitlab.com/swarmfund/horizon/db2/history"
)

// TransactionV2 represents a single, successful transaction with ledger changes
type TransactionV2 struct {
	ID              string                `json:"id"`
	PT              string                `json:"paging_token"`
	Hash            string                `json:"hash"`
	LedgerCloseTime time.Time             `json:"created_at"`
	EnvelopeXdr     string                `json:"envelope_xdr"`
	ResultXdr       string                `json:"result_xdr"`
	Changes         []LedgerEntryChangeV2 `json:"changes"`
}

type LedgerEntryChangeV2 struct {
	Effect    int32  `json:"effect"`
	EntryType int32  `json:"entry_type"`
	Payload   string `json:"payload"`
}

// Populate fills out the details
func (t *TransactionV2) Populate(transactionRow history.Transaction, ledgerChangesRow []history.LedgerChanges) error {
	t.ID = transactionRow.TransactionHash
	t.PT = transactionRow.PagingToken()
	t.Hash = transactionRow.TransactionHash
	t.LedgerCloseTime = transactionRow.LedgerCloseTime
	t.EnvelopeXdr = transactionRow.TxEnvelope
	t.ResultXdr = transactionRow.TxResult
	for _, change := range ledgerChangesRow {
		ledgerEntryChangeV2 := LedgerEntryChangeV2{
			Effect:    int32(change.Effect),
			EntryType: int32(change.EntryType),
			Payload: change.Payload,
		}
		t.Changes = append(t.Changes, ledgerEntryChangeV2)
	}
	return nil
}

// PagingToken implementation for hal.Pageable
func (t TransactionV2) PagingToken() string {
	return t.PT
}