package history

import (
	"gitlab.com/distributed_lab/tokend/horizon/db2"
	sq "github.com/lann/squirrel"
)

var selectBalance = sq.Select("hb.*").From("history_balances hb")

type BalancesQ struct {
	Err    error
	parent *Q
	sql    sq.SelectBuilder
}

type BalancesQI interface {
	ForAccount(aid string) BalancesQI
	ForExchange(eid string) BalancesQI
	ForAsset(asset string) BalancesQI
	Page(page db2.PageQuery) BalancesQI
	Select(dest interface{}) error
}

func (q *Q) Balances() BalancesQI {
	return &BalancesQ{
		parent: q,
		sql:    selectBalance,
	}
}

func (q *Q) BalanceByID(dest interface{}, id string) error {
	sql := selectBalance.Limit(1).Where("hb.balance_id = ?", id)
	return q.Get(dest, sql)
}

func (q *BalancesQ) ForAccount(aid string) BalancesQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("hb.account_id = ?", aid)

	return q
}

func (q *BalancesQ) ForExchange(eid string) BalancesQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("hb.exchange_id = ?", eid)

	return q
}

func (q *BalancesQ) ForAsset(asset string) BalancesQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("hb.asset = ?", asset)

	return q
}

func (q *BalancesQ) Page(page db2.PageQuery) BalancesQI {
	if q.Err != nil {
		return q
	}

	q.sql, q.Err = page.ApplyTo(q.sql, "hb.id")
	return q
}

// Select loads the results of the query specified by `q` into `dest`.
func (q *BalancesQ) Select(dest interface{}) error {
	if q.Err != nil {
		return q.Err
	}

	q.Err = q.parent.Select(dest, q.sql)
	return q.Err
}
