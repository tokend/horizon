package core2

import (
	sq "github.com/lann/squirrel"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2"
)

// AccountsQ is a helper struct to aid in configuring queries that loads
// account structs.
type AccountsQ struct {
	repo     *db2.Repo
	selector sq.SelectBuilder
}

// NewAccountsQ - creates new instance of AccountsQ
func NewAccountsQ(repo *db2.Repo) AccountsQ {
	return AccountsQ{
		repo: repo,
		selector: sq.Select("accounts.account_id",
			"accounts.sequential_id",
			"accounts.referrer",
			"accounts.sequential_id",
			"accounts.role_id",
		).From("accounts accounts"),
	}
}

// GetByAddress loads a row from `accounts`, by address
// returns nil, nil - if account does not exists
func (q AccountsQ) GetByAddress(address string) (*Account, error) {
	return q.FilterByAddress(address).Get()
}

//FilterByAddress - returns q with filter by address
func (q AccountsQ) FilterByAddress(address string) AccountsQ {
	q.selector = q.selector.Where("accounts.account_id = ?", address)
	return q
}

//FilterByReferrer - returns q with filter by referrer
func (q AccountsQ) FilterByReferrer(address string) AccountsQ {
	q.selector = q.selector.Where("accounts.referrer = ?", address)
	return q
}

// Get - loads a row from `accounts`
// returns nil, nil - if account does not exists
// returns error if more than one Account found
func (q AccountsQ) Get() (*Account, error) {
	var result Account
	err := q.repo.Get(&result, q.selector)
	if err != nil {
		if q.repo.NoRows(err) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load account")
	}

	return &result, nil
}

//Select - selects slice from the db, if no accounts found - returns nil, nil
func (q AccountsQ) Select() ([]Account, error) {
	var result []Account
	err := q.repo.Select(&result, q.selector)
	if err != nil {
		if q.repo.NoRows(err) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load account accounts")
	}

	return result, nil
}
