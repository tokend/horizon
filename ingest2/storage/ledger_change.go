package storage

import (
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

func convertLedgerChangesToParams(ledgerChange history2.LedgerChanges) []interface{} {
	return []interface{}{
		ledgerChange.TransactionID, ledgerChange.OperationID, ledgerChange.OrderNumber, ledgerChange.Effect,
		ledgerChange.EntryType, ledgerChange.Payload,
	}
}

// Insert - inserts new sale
func (q *LedgerChange) Insert(ledgerChanges []history2.LedgerChanges) error {
	columns := []string{"tx_id", "op_id", "order_number", "effect", "entry_type", "payload"}
	err := history2LedgerChangesBatchInsert(q.repo, ledgerChanges, "ledger_changes", columns, convertLedgerChangesToParams)
	if err != nil {
		return errors.Wrap(err, "failed to insert ledger changes")
	}
	return nil
}
