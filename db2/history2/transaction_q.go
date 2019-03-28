package history2

import (
	"fmt"
	sq "github.com/lann/squirrel"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2"
)

// TransactionsQ is a helper struct to aid in configuring queries that loads
// transactions structures.
type TransactionsQ struct {
	repo     *db2.Repo
	selector sq.SelectBuilder
}

// NewTransactionsQ - creates new instance of TransactionsQ
func NewTransactionsQ(repo *db2.Repo) TransactionsQ {
	return TransactionsQ{
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
		).
			From("transactions").
			// To apply filters on ledger_changes properties:
			Distinct().
			LeftJoin("ledger_changes ON ledger_changes.tx_id = transactions.id"),
	}
}

// FilterByLedgerEntryTypes - returns q with filter by entry types
func (q TransactionsQ) FilterByLedgerEntryTypes(types ...int) TransactionsQ {
	q.selector = q.selector.Where(sq.Eq{"ledger_changes.entry_type": types})
	return q
}

// FilterByLedgerEntryTypes - returns q with filter by effect(ledger entry change) types
func (q TransactionsQ) FilterByEffectTypes(types ...int) TransactionsQ {
	q.selector = q.selector.Where(sq.Eq{"ledger_changes.effect": types})
	return q
}

// FilterByID - returns q with filter by transaction ID
func (q TransactionsQ) FilterByID(id uint64) TransactionsQ {
	q.selector = q.selector.Where("transactions.id = ?", id)
	return q
}

// GetByID loads a row from `transactions`, by ID
// returns nil, nil - if transaction does not exists
func (q TransactionsQ) GetByID(id uint64) (*Transaction, error) {
	return q.FilterByID(id).Get()
}

// Page - returns Q with specified limit and cursor params
func (q TransactionsQ) Page(params db2.CursorPageParams) TransactionsQ {
	q.selector = params.ApplyTo(q.selector, "transactions.id")
	return q
}

// Get - loads a row from `transactions`
// returns nil, nil - if transaction does not exists
// returns error if more than one Transaction found
func (q TransactionsQ) Get() (*Transaction, error) {
	var result Transaction
	err := q.repo.Get(&result, q.selector)
	if err != nil {
		if q.repo.NoRows(err) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load transaction")
	}

	return &result, nil
}

//Select - selects slice from the db, if no sales found - returns nil, nil
func (q TransactionsQ) Select() ([]Transaction, error) {
	var result []Transaction

	fmt.Println("Result:")
	fmt.Println(q.selector.ToSql())

	err := q.repo.Select(&result, q.selector)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load transactions")
	}

	return result, nil
}
