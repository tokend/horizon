package core

import (
	sq "github.com/lann/squirrel"
)

type OrderBookQ struct {
	Err    error
	parent *Q
	sql    sq.SelectBuilder
}

func (q *Q) OrderBook() *OrderBookQ {
	return &OrderBookQ{
		parent: q,
		sql: sq.Select(
			"o.base_amount",
			"o.quote_amount",
			"o.price",
			"o.created_at",
		).From("offer o"),
	}
}

func (q *OrderBookQ) ForAssets(base, quote string) *OrderBookQ {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("base_asset_code = ? AND quote_asset_code = ?", base, quote)
	return q
}

func (q *OrderBookQ) Direction(isBuy bool) *OrderBookQ {
	if q.Err != nil {
		return q
	}

	orderDirection := "ASC"
	if isBuy {
		orderDirection = "DESC"
	}

	q.sql = q.sql.Where("is_buy = ?", isBuy).OrderBy("price " + orderDirection).OrderBy("created_at ASC")
	return q
}

// Select loads the results of the query specified by `q` into `dest`.
func (q *OrderBookQ) Select(dest interface{}) error {
	if q.Err != nil {
		return q.Err
	}

	q.Err = q.parent.Select(dest, q.sql)
	return q.Err
}
