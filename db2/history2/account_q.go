package history2

import (
	sq "github.com/lann/squirrel"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/bridge"
)

// AccountsQ is a helper struct to aid in configuring queries that loads accounts
type AccountsQ struct {
	repo     *bridge.Mediator
	selector sq.SelectBuilder
}

// NewAccountsQ - creates new instance of AccountsQ
func NewAccountsQ(repo *bridge.Mediator) AccountsQ {
	return AccountsQ{
		repo:     repo,
		selector: sq.Select(accountColumns...).From("accounts accounts"),
	}
}

// ByAddress loads a row from `accounts`, by address
// returns nil, nil - if account does not exists
func (q AccountsQ) ByAddress(address string) (*Account, error) {
	q.selector = q.selector.Where("accounts.address = ?", address)
	return q.Get()
}

//ByID - gets account by ID, returns nil, nil if one does not exist
func (q AccountsQ) ByID(id uint64) (*Account, error) {
	q.selector = q.selector.Where("accounts.id = ?", id)
	return q.Get()
}

//Get - selects account from db, returns nil, nil if one does not exists
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

var accountColumns = []string{"accounts.id", "accounts.address", "accounts.kyc_recovery_status"}
