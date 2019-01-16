package core2

import (
	"fmt"

	sq "github.com/lann/squirrel"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2"
)

// BalancesQ is a helper struct to aid in configuring queries that loads balances
type BalancesQ struct {
	repo     *db2.Repo
	selector sq.SelectBuilder
}

// NewBalancesQ - creates new instance of BalanceQ with no filters
func NewBalancesQ(repo *db2.Repo) BalancesQ {
	return BalancesQ{
		repo: repo,
		selector: sq.Select("balances.balance_id", "balances.sequential_id", "balances.asset", "balances.account_id",
			"balances.amount", "balances.locked").From("balance balances"),
	}
}

//GetByAddress  - gets balance by address, if one does not exists - returns nil, nil
func (q BalancesQ) GetByAddress(address string) (*Balance, error) {
	return q.FilterByAddress(address).Get()
}

// FilterByAddress - returns new Q with added filter for balance address
func (q BalancesQ) FilterByAddress(address string) BalancesQ {
	q.selector = q.selector.Where("balances.balance_id = ?", address)
	return q
}

//FilterByAccount - returns new Q with added filter for account address
func (q BalancesQ) FilterByAccount(accountAddress string) BalancesQ {
	q.selector = q.selector.Where("balances.account_id = ?", accountAddress)
	return q
}

//WithAsset - joins asset
func (q BalancesQ) WithAsset() BalancesQ {
	q.selector = q.selector.Columns(getAssetColumns()...).LeftJoin("asset assets ON balances.asset = assets.code")
	return q
}

func getAssetColumns() []string {
	result := make([]string, 0, len(assetColumns))
	for _, column := range assetColumns {
		result = append(result, fmt.Sprintf(`%s "%s"`, column, column))
	}

	return result
}

// Get - selects balance from db using specified filters. Returns nil, nil - if one does not exists
// Returns error if more than one exists
func (q BalancesQ) Get() (*Balance, error) {
	var result Balance
	err := q.repo.Get(&result, q.selector)
	if err != nil {
		if q.repo.NoRows(err) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to get balance")
	}

	return &result, nil
}

// Select - selects balances from fb using specified filters. Returns nil, nil - if one does not exists
func (q BalancesQ) Select() ([]Balance, error) {
	var result []Balance
	err := q.repo.Select(&result, q.selector)
	if err != nil {
		if q.repo.NoRows(err) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to select balances")
	}

	return result, nil
}
