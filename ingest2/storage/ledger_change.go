package storage

import (
	sq "github.com/lann/squirrel"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2"
	"gitlab.com/tokend/horizon/db2/history2"
)

// LedgerChange is helper struct to operate with `ledger changes`
type LedgerChange struct {
	repo *db2.Repo
}

// NewLedgerChange - creates new instance of the `ledger change`
func NewLedgerChange(repo *db2.Repo) *LedgerChange {
	return &LedgerChange{
		repo: repo,
	}
}

// Insert - inserts new sale
func (q *LedgerChange) Insert(ledgerChanges []history2.LedgerChanges) error {
	sql := sq.Insert("ledger_changes").
		Columns(
			"tx_id", "op_id", "order_number", "effect", "entry_type", "payload",
		)

	for _, ledgerChange := range ledgerChanges {
		sql.Values(ledgerChange.TransactionID, ledgerChange.OperationID, ledgerChange.OrderNumber, ledgerChange.Effect,
			ledgerChange.EntryType, ledgerChange.Payload)
	}

	_, err := q.repo.Exec(sql)
	if err != nil {
		return errors.Wrap(err, "failed to insert ledger changes", logan.F{"ledger_changes_len": len(ledgerChanges)})
	}

	return nil
}
