package db2

import (
	"fmt"
	sq "github.com/lann/squirrel"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

// OffsetPageParams - page params of the db query
type OffsetPageParams struct {
	Limit      uint64
	PageNumber uint64
	Order      OrderType
}

// ApplyTo returns a new SelectBuilder after applying the paging effects of
// `p` to `sql`.
func (p *OffsetPageParams) ApplyTo(sql sq.SelectBuilder, cols ...string) sq.SelectBuilder {
	offset := p.Limit * p.PageNumber

	sql = sql.Limit(p.Limit).Offset(offset)

	switch p.Order {
	case OrderAsc:
		for _, col := range cols {
			sql = sql.OrderBy(fmt.Sprintf("%s %s", col, "asc"))
		}
	case OrderDesc:
		for _, col := range cols {
			sql = sql.OrderBy(fmt.Sprintf("%s %s", col, "desc"))
		}
	default:
		panic(errors.From(errors.New("unexpected order type"), logan.F{
			"order_type": p.Order,
		}))
	}

	return sql
}
