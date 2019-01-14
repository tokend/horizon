package history2

import (
	sq "github.com/lann/squirrel"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2"
)

// AccountsQ is a helper struct to aid in configuring queries that loads accounts
type AccountsQ struct {
	repo *db2.Repo
}

// NewAccountsQ - creates new instance of AccountsQ
func NewAccountsQ(repo *db2.Repo) *AccountsQ {
	return &AccountsQ{
		repo: repo,
	}
}

// GetByAddress loads a row from `accounts`, by address
// returns nil, nil - if account does not exists
func (q *AccountsQ) ByAddress(address string) (*Account, error) {
	var result Account
	err := q.repo.Get(&result, sq.Select("a.id, a.address").From("accounts a").
		Where("a.address = ?", address))
	if err != nil {
		if q.repo.NoRows(err) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load account by address", logan.F{"address": address})
	}

	return &result, nil
}
