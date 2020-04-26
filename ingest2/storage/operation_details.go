package storage

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/bridge"
	"gitlab.com/tokend/horizon/db2/history2"
)

// Operation - helper struct to store operation
type Operation struct {
	repo *bridge.Mediator
}

// NewOperationDetails - creates new instance of `Operation`
func NewOperationDetails(repo *bridge.Mediator) *Operation {
	return &Operation{
		repo: repo,
	}
}

func convertOpDetails(op history2.Operation) []interface{} {
	return []interface{}{
		op.ID, op.Type, op.Details, op.LedgerCloseTime, op.Source, op.TxID,
	}
}

//Insert - stores slice of the operations via batch insert.
func (s *Operation) Insert(ops []history2.Operation) error {
	columns := []string{"id", "type", "details", "ledger_close_time", "source", "tx_id"}
	err := history2OperationBatchInsert(s.repo, ops, "operations", columns, convertOpDetails)
	if err != nil {
		return errors.Wrap(err, "failed to insert operation details")
	}
	return nil
}
