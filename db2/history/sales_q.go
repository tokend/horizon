package history

import (
	"time"

	"fmt"

	sq "github.com/lann/squirrel"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/horizon/db2"
)

type SalesQ interface {
	// ByID - selects sale by specified id. Returns nil, nil if not found
	ByID(saleID uint64) (*Sale, error)
	// ForOwner filters sale by owner
	ForOwner(ownerID string) SalesQ
	// ForBaseAsset - filters by base asset
	ForBaseAsset(baseAsset string) SalesQ
	// ForName - filters by `name` field in the `details` json.
	ForName(baseAsset string) SalesQ
	// Open - selects only open sales
	Open(now time.Time) SalesQ
	// Upcoming - selects only upcoming sales.
	Upcoming(now time.Time) SalesQ
	// CollectedValueBound - selects all sales in which the `current_cap` is above bound.
	CollectedValueBound(bound int64) SalesQ
	// CurrentSoftCapsRatio is selects all sales in which the `current_cap`
	// is filled by more than a percentBound of the `soft_cap`.
	CurrentSoftCapsRatio(percentBound int64) SalesQ
	// OrderByEndTime is set ordering by `end_time`.
	OrderByEndTime() SalesQ
	// OrderByCurrentCap is set ordering by `current_cap`.
	OrderByCurrentCap(desc bool) SalesQ
	// OrderByPopularity is merge with quantity of the
	// unique investors for each sale, and sort sales by quantity.
	OrderByPopularity(values db2.OrderBooksInvestors) SalesQ
	// Insert - inserts new sale
	Insert(sale Sale) error
	// Update - updates existing sale
	Update(sale Sale) error
	// SetState - sets state
	SetState(saleID uint64, state SaleState) error
	// Select - selects slice of Sales using specified filters
	Select() ([]Sale, error)
	// Page specifies the paging constraints for the query being built by `q`.
	Page(page db2.PageQuery) SalesQ
}

type saleQ struct {
	Err    error
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

	q.sql = q.sql.Where("start_time > ?", now)
	return q
}

// CollectedValueBound - selects all sales in which the `current_cap` is above bound.
func (q *saleQ) CollectedValueBound(bound int64) SalesQ {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("current_cap >= ?", bound)
	return q
}

// ReachedSoftCap - selects all sales in which the `current_cap` is above `soft_cap`.
func (q *saleQ) ReachedSoftCap() SalesQ {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("current_cap > soft_cap")
	return q
}

// CurrentSoftCapsRatio is selects all sales in which the `current_cap` is filled by more than a percentBound of the `soft_cap`.
func (q *saleQ) CurrentSoftCapsRatio(percentBound int64) SalesQ {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.
		Where("div((current_cap * 100 ), soft_cap) > ?", percentBound)

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

// Insert - inserts new sale
func (q *saleQ) Insert(sale Sale) error {
	sql := sq.Insert("sale").
		Columns(
			"id", "owner_id", "base_asset", "quote_asset", "start_time", "end_time",
			"price", "soft_cap", "hard_cap", "current_cap", "details", "state",
		).
		Values(
			sale.ID, sale.OwnerID, sale.BaseAsset, sale.QuoteAsset, sale.StartTime, sale.EndTime,
			sale.Price, sale.SoftCap, sale.HardCap, sale.CurrentCap, sale.Details, sale.State,
		)

	_, err := q.parent.Exec(sql)
	if err != nil {
		return errors.Wrap(err, "failed to insert sale")
	}

	return nil
}

// Update - updates existing sale
func (q *saleQ) Update(sale Sale) error {
	sql := sq.Update("sale").SetMap(map[string]interface{}{
		"owner_id":    sale.OwnerID,
		"base_asset":  sale.BaseAsset,
		"quote_asset": sale.QuoteAsset,
		"start_time":  sale.StartTime,
		"end_time":    sale.EndTime,
		"price":       sale.Price,
		"soft_cap":    sale.SoftCap,
		"hard_cap":    sale.HardCap,
		"current_cap": sale.CurrentCap,
		"details":     sale.Details,
		"state":       sale.State,
	}).Where("id = ?", sale.ID)

	_, err := q.parent.Exec(sql)
	if err != nil {
		return errors.Wrap(err, "failed to update sale", logan.F{"sale_id": sale.ID})
	}

	return nil
}

// SetState - sets state
func (q *saleQ) SetState(saleID uint64, state SaleState) error {
	sql := sq.Update("sale").Set("state", state).Where("id = ?", saleID)
	_, err := q.parent.Exec(sql)
	if err != nil {
		return errors.Wrap(err, "failed to set state", logan.F{"sale_id": saleID})
	}

	return nil
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
	q.sql = q.sql.Join(
		fmt.Sprintf(
			"(values %s) as investors_count(order_book_id, quantity) on id = investors_count.order_book_id",
			values.String()),
	).OrderBy("investors_count.quantity DESC")

	return q
}

var selectSales = sq.Select(
	"id", "owner_id", "base_asset", "quote_asset", "start_time", "end_time", "price", "soft_cap", "hard_cap",
	"current_cap", "details", "state").From("sale")
