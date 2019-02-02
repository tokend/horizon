package history2

import (
	sq "github.com/lann/squirrel"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2"
)

// BalancesQ is a helper struct to aid in configuring queries that loads balances
type BalancesQ struct {
	repo *db2.Repo
}

// NewBalancesQ - creates new instance of BalanceQ
func NewBalancesQ(repo *db2.Repo) BalancesQ {
	return BalancesQ{
		repo: repo,
	}
}

// GetByAddress loads a row from `balances`, by address
// returns nil, nil - if balance does not exists
func (q BalancesQ) GetByAddress(address string) (*Balance, error) {
	var result Balance
	err := q.repo.Get(&result, sq.Select(balanceColumns...).
		From("balances balances").Where("balances.address = ?", address))
	if err != nil {
		if q.repo.NoRows(err) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load balance by address", logan.F{"balance_address": address})
	}

	return &result, nil
}

var balanceColumns = []string{"balances.id", "balances.account_id", "balances.address", "balances.asset_code"}
