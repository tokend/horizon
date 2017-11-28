package history

import (
	"gitlab.com/swarmfund/horizon/db2"
	sq "github.com/lann/squirrel"
)

type TradesQI interface {
	ForPair(base, quote string) TradesQI
	Page(page db2.PageQuery) TradesQI
	Select(dest interface{}) error
}

type TradesQ struct {
	Err    error
	parent *Q
	sql    sq.SelectBuilder
}

func (q *Q) Trades() TradesQI {
	return &TradesQ{
		parent: q,
		sql:    selectTrades,
	}
}

func (q *TradesQ) ForPair(base, quote string) TradesQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("ht.base_asset = ? AND ht.quote_asset = ?", base, quote)
	return q
}

func (q *TradesQ) Page(page db2.PageQuery) TradesQI {
	if q.Err != nil {
		return q
	}

	q.sql, q.Err = page.ApplyTo(q.sql, "ht.id")
	return q
}

// Select loads the results of the query specified by `q` into `dest`.
func (q *TradesQ) Select(dest interface{}) error {
	if q.Err != nil {
		return q.Err
	}

	q.Err = q.parent.Select(dest, q.sql)
	return q.Err
}

var selectTrades = sq.Select(
	"ht.id",
	"ht.base_asset",
	"ht.quote_asset",
	"ht.base_amount",
	"ht.quote_amount",
	"ht.price",
	"ht.created_at").
	From("history_trades ht")
