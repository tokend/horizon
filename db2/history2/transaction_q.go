package history2

import (
	sq "github.com/lann/squirrel"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2"
)

// TransactionQ is a helper struct to aid in configuring queries that loads
// tx structures.
type TransactionQ struct {
	repo     *db2.Repo
	selector sq.SelectBuilder
}

// NewTransactionQ - creates new instance of TransactionQ
func NewTransactionQ(repo *db2.Repo) *TransactionQ {
	return &TransactionQ{
		repo: repo,
		selector: sq.Select(
			"transactions.id",
			"transactions.hash",
			"transactions.ledger_sequence",
			"transactions.ledger_close_time",
			"transactions.application_order",
			"transactions.account",
			"transactions.operation_count",
			"transactions.envelope",
			"transactions.result",
			"transactions.meta",
			"transactions.valid_after",
			"transactions.valid_before",
		).From("transactions"),
	}
}

// Get - loads a row from `transactions`
// returns nil, nil - if tx does not exists
// returns error if more than one transaction found
func (q TransactionQ) Get() (*Transaction, error) {
	var result Transaction
	err := q.repo.Get(&result, q.selector)
	if err != nil {
		if q.repo.NoRows(err) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load tx")
	}

	return &result, nil
}

//Select - selects slice from the db, if no transactions found - returns nil, nil
func (q TransactionQ) Select() ([]Transaction, error) {
	var result []Transaction
	err := q.repo.Select(&result, q.selector)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load transactions")
	}

	return result, nil
}

func (q *TransactionQ) GetByHash(hash string) (*Transaction, error) {
	q.selector = q.selector.Where("transactions.hash = ?", hash)
	return q.Get()
}
