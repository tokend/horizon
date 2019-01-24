package db2

import (
	"fmt"

	sq "github.com/lann/squirrel"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

//OrderType - represents sorting order of the query
type OrderType string

//Invert - inverts order by
func (o OrderType) Invert() OrderType {
	switch o {
	case OrderDesc:
		return OrderAsc
	case OrderAsc:
		return OrderDesc
	default:
		panic(errors.From(errors.New("unexpected order type"), logan.F{
			"order_type": o,
		}))
	}
}

const (
	// OrderAsc - ascending order
	OrderAsc OrderType = "asc"
	// OrderDesc - descending order
	OrderDesc OrderType = "desc"
)

//CursorPageParams - page params of the db query
type CursorPageParams struct {
	Cursor uint64
	Order  OrderType
	Limit  uint64
}

// ApplyTo returns a new SelectBuilder after applying the paging effects of
// `p` to `sql`.  This method provides the default case for paging: int64
// cursor-based paging by an id column.
func (p CursorPageParams) ApplyTo(sql sq.SelectBuilder, col string) sq.SelectBuilder {
	sql = sql.Limit(p.Limit)

	switch p.Order {
	case OrderAsc:
		sql = sql.
			Where(fmt.Sprintf("%s > ?", col), p.Cursor).
			OrderBy(fmt.Sprintf("%s asc", col))
	case OrderDesc:
		sql = sql.
			Where(fmt.Sprintf("%s < ?", col), p.Cursor).
			OrderBy(fmt.Sprintf("%s desc", col))
	default:
		panic(errors.From(errors.New("unexpected order type"), logan.F{
			"order_type": p.Order,
		}))
	}

	return sql
}
