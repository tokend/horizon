package db2

import (
	sq "github.com/Masterminds/squirrel"
	"gitlab.com/distributed_lab/kit/pgdb"
)

type OffsetPageParams struct {
	Limit      uint64
	Order      OrderType
	PageNumber uint64
}

type CursorPageParams struct {
	Limit  uint64
	Order  OrderType
	Cursor uint64
}

// ApplyTo returns a new SelectBuilder after applying the paging effects of
// `p` to `sql`.
func (p *OffsetPageParams) ApplyTo(sql sq.SelectBuilder, cols ...string) sq.SelectBuilder {
	offsetPageParams := pgdb.OffsetPageParams{
		Limit:      p.Limit,
		PageNumber: p.PageNumber,
		Order:      string(p.Order),
	}

	return offsetPageParams.ApplyTo(sql, cols...)
}

// ApplyTo returns a new SelectBuilder after applying the paging effects of
// `p` to `sql`.  This method provides the default case for paging: int64
// cursor-based paging by an id column.
func (p *CursorPageParams) ApplyTo(sql sq.SelectBuilder, col string) sq.SelectBuilder {
	cursorPageParams := pgdb.CursorPageParams{
		Limit:  p.Limit,
		Cursor: p.Cursor,
		Order:  string(p.Order),
	}
	return cursorPageParams.ApplyTo(sql, col)
}
