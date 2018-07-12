package resource

import (
	"time"
	"gitlab.com/swarmfund/horizon/db2/history"
)

// TransactionV2 represents a single, successful transaction
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
	Effect    string `json:"effect"`
	EntryType string `json:"entry_type"`
	Payload   string `json:"payload"`
}

func (t *TransactionV2) Populate(transactionRow history.Transaction, ledgerChangesRow history.LedgerChanges) error {
	t.ID = transactionRow.TransactionHash
	t.PT = transactionRow.PagingToken()
	t.Hash = transactionRow.TransactionHash
	t.LedgerCloseTime = transactionRow.LedgerCloseTime
	t.EnvelopeXdr = transactionRow.TxEnvelope
	t.ResultXdr = transactionRow.TxResult
	return nil
}