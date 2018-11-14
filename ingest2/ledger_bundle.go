package ingest2

import "gitlab.com/tokend/horizon/db2/core"

// ledgerBundle represents a single ledger's worth of novelty created by one
// ledger close
type ledgerBundle struct {
	Sequence        int32
	Header          core.LedgerHeader
	Transactions    []core.Transaction
}

func newLedgerBundle(seq int32, header core.LedgerHeader, txs []core.Transaction) ledgerBundle{
	return ledgerBundle{
		Sequence: seq,
		Header: header,
		Transactions: txs,
	}
}