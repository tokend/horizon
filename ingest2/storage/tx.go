package storage

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/bridge"
	"gitlab.com/tokend/horizon/db2/history2"
)

//Tx - handles write operations in DB for txs
type Tx struct {
	repo *bridge.Mediator
}

//NewTx - creates new instance of Tx
func NewTx(repo *bridge.Mediator) *Tx {
	return &Tx{
		repo: repo,
	}
}

func convertTxToParams(tx history2.Transaction) []interface{} {
	return []interface{}{
		tx.ID, tx.Hash, tx.LedgerSequence, tx.LedgerCloseTime, tx.ApplicationOrder, tx.Account,
		tx.OperationCount, tx.Envelope, tx.Result, tx.Meta,
		tx.ValidAfter, tx.ValidBefore,
	}
}

//Insert - inserts slice of txs into DB in one slice
func (s *Tx) Insert(txs []history2.Transaction) error {
	columns := []string{"id", "hash", "ledger_sequence", "ledger_close_time", "application_order", "account",
		"operation_count", "envelope", "result", "meta", "valid_after", "valid_before"}
	err := history2TransactionBatchInsert(s.repo, txs, "transactions", columns, convertTxToParams)
	if err != nil {
		return errors.Wrap(err, "failed to insert txs")
	}
	return nil
}
