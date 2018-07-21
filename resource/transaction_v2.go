package resource

import (
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/tokend/regources"
)

// Populate fills out the details
func PopulateTransactionV2(transactionRow history.Transaction, ledgerChangesRows []history.LedgerChanges,
) (t regources.TransactionV2) {
	t.ID = transactionRow.TransactionHash
	t.PT = transactionRow.PagingToken()
	t.Hash = transactionRow.TransactionHash
	t.LedgerCloseTime = transactionRow.LedgerCloseTime
	t.EnvelopeXDR = transactionRow.TxEnvelope
	t.ResultXDR = transactionRow.TxResult
	t.LedgerSequence = transactionRow.LedgerSequence
	for _, change := range ledgerChangesRows {
		ledgerEntryChangeV2 := regources.LedgerEntryChangeV2{
			Effect:    int32(change.Effect),
			EntryType: int32(change.EntryType),
			Payload: change.Payload,
		}
		t.Changes = append(t.Changes, ledgerEntryChangeV2)
	}

	return
}
