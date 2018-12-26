package storage

import (
	sq "github.com/lann/squirrel"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2"
	"gitlab.com/tokend/horizon/db2/history2"
)

type Tx struct {
	repo *db2.Repo
}

func NewTx(repo *db2.Repo) *Tx {
	return &Tx{
		repo: repo,
	}
}

func (s *Tx) Insert(txs []history2.Transaction) error {
	sql := sq.Insert("tx").
		Columns(
			"id", "tx_hash", "ledger_sequence", "ledger_close_time", "application_order", "account",
			"operation_count", "envelope", "result", "meta", "valid_after", "valid_before",
		)

	for _, tx := range txs {
		sql = sql.Values(tx.ID, tx.TxHash, tx.LedgerSequence, tx.LedgerCloseTime, tx.ApplicationOrder, tx.Account,
			tx.OperationCount, tx.Envelope, tx.Result, tx.Meta,
			tx.ValidAfter, tx.ValidBefore)
	}

	_, err := s.repo.Exec(sql)
	if err != nil {
		return errors.Wrap(err, "failed to insert txs", logan.F{"txs_len": len(txs)})
	}

	return nil
}
