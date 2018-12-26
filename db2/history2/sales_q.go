package history2
/*
import (
	"time"

	"fmt"

	sq "github.com/lann/squirrel"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2"
	"gitlab.com/tokend/horizon/db2/sqx"
)


type saleQ struct {
	parent *Q
	sql    sq.SelectBuilder
}

func (q *Q) Sales() SalesQ {
	return &saleQ{
		parent: q,
		sql:    selectSales,
	}
}

// ForOwner filters sale by owner
func (q *saleQ) ForOwner(ownerID string) SalesQ {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("owner_id = ?", ownerID)
	return q
}

// ForBaseAsset - filters by base asset
func (q *saleQ) ForBaseAsset(baseAsset string) SalesQ {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("base_asset = ?", baseAsset)
	return q
}

// ForName - filters by `name` field in the `details` json.
func (q *saleQ) ForName(name string) SalesQ {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("details ->> 'name' ilike ?", fmt.Sprint("%", name, "%"))
	return q
}

// Open - selects only open sales
func (q *saleQ) Open(now time.Time) SalesQ {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("state = ? AND end_time >= ?", SaleStateOpen, now)
	return q
}

// Upcoming - selects only upcoming sales.
func (q *saleQ) Upcoming(now time.Time) SalesQ {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("state = ? AND start_time > ?", SaleStateOpen, now)
	return q
}

// Promotions - selects only promotion sales.
func (q *saleQ) Promotions() SalesQ {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("state = ?", SaleStatePromotion)
	return q
}

// Voting - selects only voting sales.
func (q *saleQ) Voting() SalesQ {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("state = ?", SaleStateVoting)
	return q
}

// ByID - selects sale by specified id. Returns nil, nil if not found
func (q *saleQ) ByID(saleID uint64) (*Sale, error) {
	if q.Err != nil {
		return nil, errors.Wrap(q.Err, "error for q builder")
	}

	q.sql = q.sql.Where("id = ?", saleID)
	var result Sale
	err := q.parent.Get(&result, q.sql)
	if q.parent.NoRows(err) {
		return nil, nil
	}

	if err != nil {
		return nil, errors.Wrap(err, "failed to select sale by id")
	}

	return &result, nil
}

// Select - selects slice of Sales using specified filters
func (q *saleQ) Select() ([]Sale, error) {
	if q.Err != nil {
		return nil, errors.Wrap(q.Err, "error from query builder")
	}

	var result []Sale
	err := q.parent.Select(&result, q.sql)
	if q.parent.NoRows(err) {
		return nil, nil
	}

	if err != nil {
		return nil, errors.Wrap(err, "failed to select sales")
	}

	return result, nil
}

// Page specifies the paging constraints for the query being built by `q`.
func (q *saleQ) Page(page db2.PageQuery) SalesQ {
	if q.Err != nil {
		return q
	}

	q.sql, q.Err = page.ApplyTo(q.sql, "sale.id")
	return q
}

func (q *saleQ) PageV2(page db2.PageQueryV2) SalesQ {
	if q.Err != nil {
		return q
	}

	q.sql, q.Err = page.ApplyTo(q.sql)

	return q
}

// OrderByEndTime is set ordering by `start_time`.
func (q *saleQ) OrderByStartTime() SalesQ {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.OrderBy("start_time ASC")
	return q
}

// OrderById assumes order were validate beforehand and it's value is one of db accepts
func (q *saleQ) OrderById(order string) SalesQ {
	if q.Err != nil {
		return q
	}
	q.sql = q.sql.OrderBy(fmt.Sprintf("sale.id %s", order))
	return q
}

// OrderByEndTime is set ordering by `end_time`.
func (q *saleQ) OrderByEndTime() SalesQ {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.OrderBy("end_time ASC")
	return q
}

// OrderByCurrentCap is set ordering by `current_cap`.
func (q *saleQ) OrderByCurrentCap(desc bool) SalesQ {
	if q.Err != nil {
		return q
	}
	order := "ASC"
	if desc {
		order = "DESC"
	}
	q.sql = q.sql.OrderBy(fmt.Sprintf("current_cap %s", order))

	return q
}

// OrderByPopularity is merge with quantity of the unique investors for each sale,
// and sort sales by quantity.
func (q *saleQ) OrderByPopularity(values db2.OrderBooksInvestors) SalesQ {
	if q.Err != nil {
		return q
	}
	q.sql = q.sql.LeftJoin(
		fmt.Sprintf(
			"(values %s) as investors_count(order_book_id, quantity) on id = investors_count.order_book_id",
			values.String()),
	).OrderBy("investors_count.quantity DESC NULLS LAST")

	return q
}

// ForBaseAssets - filters by base assets
func (q *saleQ) ForBaseAssets(baseAssets ...string) SalesQ {
	if q.Err != nil {
		return q
	}

	query, values := sqx.InForString("base_asset", baseAssets...)
	q.sql = q.sql.Where(query, values...)
	return q
}

var selectSales = sq.Select(
	"id", "owner_id", "base_asset", "default_quote_asset", "start_time", "end_time", "quote_assets", "soft_cap", "hard_cap",
	"details", "state", "base_hard_cap", "base_current_cap, sale_type").From("sale")
*/