package storage

import (
	"gitlab.com/tokend/horizon/db2"
	"gitlab.com/tokend/horizon/db2/history2"
	sq "github.com/lann/squirrel"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

// Operation - helper struct to store operation
type Operation struct {
	repo *db2.Repo
}

// NewOperationDetails - creates new instance of `Operation`
func NewOperationDetails(repo *db2.Repo) *Operation {
	return &Operation{
		repo: repo,
	}
}

func (s *Operation) Insert(ops []history2.Operation) error {
	query := sq.Insert("operations").Columns("id, op_type, details, ledger_close_time, source")
	for _, op := range ops {
		query = query.Values(op.ID, op.Type, op.Details, op.LedgerCloseTime, op.Source)
	}

	_, err := s.repo.Exec(query)
	if err != nil {
		return errors.Wrap(err, "failed to insert operations")
	}

	return nil
}
