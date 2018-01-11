package core

import sq "github.com/lann/squirrel"

var selectBalance = sq.Select(
	"ba.balance_id",
	"ba.account_id",
	"ba.asset",
	"ba.amount",
	"ba.locked",
).From("balance ba")

type BalancesQ struct {
	Err    error
	parent *Q
	sql    sq.SelectBuilder
}

type BalancesQI interface {
	// ByAsset filters `balances` by asset code.
	ByAsset(asset string) BalancesQI
	// ByAsset filters `balances` by account address.
	ByAddress(address string) BalancesQI
	Select() ([]Balance, error)
}

func (q *Q) Balances() BalancesQI {
	return &BalancesQ{
		parent: q,
		sql:    selectBalance,
	}
}

func (q *BalancesQ) ByAddress(address string) BalancesQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("ba.account_id = ?", address)
	return q
}

func (q *BalancesQ) ByAsset(asset string) BalancesQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("ba.asset = ?", asset)
	return q
}

// Select loads the results of the query specified by `q` into `dest`.
func (q *BalancesQ) Select() ([]Balance, error) {
	if q.Err != nil {
		return nil, q.Err
	}
	dest := make([]Balance, 0)
	q.Err = q.parent.Select(&dest, q.sql)
	return dest, q.Err
}
