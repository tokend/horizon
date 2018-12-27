package storage

import (
	sq "github.com/lann/squirrel"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2"
	"gitlab.com/tokend/horizon/db2/history2"
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

//Insert - stores slice of the operations via batch insert.
func (s *Operation) Insert(ops []history2.Operation) error {
	//TODO: might have issue due to the limit of params supported by the Postgres, so it might be good idea to refactor
	// it to split slice into smaller batches
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
