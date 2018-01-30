package core

import (
	sq "github.com/lann/squirrel"
)

type SaleQ struct {
	Err error
	parent *Q
	sql sq.SelectBuilder
}

func (q *Q) Sales() *SaleQ {
	return &SaleQ{
		parent:q,
		sql: selectSale,
	}
}

func (q *SaleQ) Select(dest interface{}) error {
	if q.Err != nil {
		return q.Err
	}

	q.Err = q.parent.Select(dest, q.sql)
	return q.Err
}


var selectSale = sq.Select(
	"s.id",
	"s.owner_id",
	"s.base_asset",
	"s.default_quote_asset",
	"s.start_time",
	"s.end_time",
	"s.soft_cap",
	"s.hard_cap",
	"s.hard_cap_in_base",
	"s.current_cap_in_base",
	"s.details",
	"s.base_balance",
).From("sale s")
