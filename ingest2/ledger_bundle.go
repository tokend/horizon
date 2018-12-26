package ingest2

import core "gitlab.com/tokend/horizon/db2/core2"

// ledgerBundle represents a single ledger's worth of novelty created by one
// ledger close
type ledgerBundle struct {
	Header       core.LedgerHeader
	Transactions []core.Transaction
}

func newLedgerBundle(header core.LedgerHeader, txs []core.Transaction) ledgerBundle {
	return ledgerBundle{
		Header:       header,
		Transactions: txs,
	}
}
