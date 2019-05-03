package history2

import (
	"time"

	"gitlab.com/tokend/go/xdr"

	sq "github.com/lann/squirrel"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2"
)

// SalesQ is a helper struct to aid in configuring queries that loads
// sales structures.
type SalesQ struct {
	repo     *db2.Repo
	selector sq.SelectBuilder
}

// NewSalesQ - creates new instance of SalesQ
func NewSalesQ(repo *db2.Repo) SalesQ {
	return SalesQ{
		repo: repo,
		selector: sq.Select(
			"sales.id",
			"sales.soft_cap",
			"sales.hard_cap",
			"sales.base_current_cap",
			"sales.base_hard_cap",
			"sales.sale_type",
			"sales.owner_address",
			"sales.base_asset",
			"sales.default_quote_asset",
			"sales.start_time",
			"sales.end_time",
			"sales.details",
			"sales.state",
			"sales.quote_assets",
			"sales.version",
		).From("sales"),
	}
}

func (q SalesQ) Whitelisted(address string) SalesQ {
	subQuery := sq.StatementBuilder.
		Select("(sr.key#>>'{sale,saleID}')::int").
		From("account_specific_rules sr").
		Where(sq.Or{sq.Expr("sr.address = ?", address), sq.Expr("sr.address is null")}).
		GroupBy("sr.key#>>'{sale,saleID}'").Having("bool_and(sr.forbids)")
	q.selector = q.selector.Where(
		subQuery.
			Prefix("sales.id not in (").
			Suffix(") or sales.version < ?",
				int32(xdr.LedgerVersionAddSaleWhitelists)))
	return q
}

// FilterByID - returns q with filter by sale ID
func (q SalesQ) FilterByID(id uint64) SalesQ {
	q.selector = q.selector.Where("sales.id = ?", id)
	return q
}

// GetByID loads a row from `sales`, by ID
// returns nil, nil - if sale does not exists
func (q SalesQ) GetByID(id uint64) (*Sale, error) {
	return q.FilterByID(id).Get()
}

// FilterByOwner - returns q with filter by Owner
func (q SalesQ) FilterByOwner(ownerAddress string) SalesQ {
	q.selector = q.selector.Where("sales.owner_address = ?", ownerAddress)
	return q
}

// FilterByBaseAsset - returns q with filter by BaseAsset
func (q SalesQ) FilterByBaseAsset(baseAssetCode string) SalesQ {
	q.selector = q.selector.Where("sales.base_asset = ?", baseAssetCode)
	return q
}

// FilterByMaxEndTime - returns q with filter by max end time
func (q SalesQ) FilterByMaxEndTime(time time.Time) SalesQ {
	q.selector = q.selector.Where("sales.end_time <= ?", time)
	return q
}

// FilterByMaxStartTime - returns q with filter by start_time
func (q SalesQ) FilterByMaxStartTime(time time.Time) SalesQ {
	q.selector = q.selector.Where("sales.start_time <= ?", time)
	return q
}

// FilterByMinStartTime - returns q with filter by start_time
func (q SalesQ) FilterByMinStartTime(time time.Time) SalesQ {
	q.selector = q.selector.Where("sales.start_time >= ?", time)
	return q
}

// FilterByMinEndTime - returns q with filter by end_time
func (q SalesQ) FilterByMinEndTime(time time.Time) SalesQ {
	q.selector = q.selector.Where("sales.end_time >= ?", time)
	return q
}

// FilterByState - returns q with filter by state
func (q SalesQ) FilterByState(state uint64) SalesQ {
	q.selector = q.selector.Where("sales.state = ?", state)
	return q
}

// FilterBySaleType - returns q with filter by type
func (q SalesQ) FilterBySaleType(saleType uint64) SalesQ {
	q.selector = q.selector.Where("sales.sale_type = ?", saleType)
	return q
}

// FilterByMinHardCap - returns q with filter by min hard cap
func (q SalesQ) FilterByMinHardCap(value uint64) SalesQ {
	q.selector = q.selector.Where("sales.hard_cap >= ?", value)
	return q
}

// FilterByMinSoftCap - returns q with filter by min soft cap
func (q SalesQ) FilterByMinSoftCap(value uint64) SalesQ {
	q.selector = q.selector.Where("sales.soft_cap >= ?", value)
	return q
}

// FilterByMaxHardCap - returns q with filter by max hard cap
func (q SalesQ) FilterByMaxHardCap(value uint64) SalesQ {
	q.selector = q.selector.Where("sales.hard_cap <= ?", value)
	return q
}

// FilterByMaxSoftCap - returns q with filter by max sof cap
func (q SalesQ) FilterByMaxSoftCap(value uint64) SalesQ {
	q.selector = q.selector.Where("sales.soft_cap <= ?", value)
	return q
}

// Page - returns Q with specified limit and offset params
func (q SalesQ) Page(params db2.OffsetPageParams) SalesQ {
	q.selector = params.ApplyTo(q.selector, "sales.id")
	return q
}

// CursorPage - returns Q with specified limit and offset params
func (q SalesQ) CursorPage(params db2.CursorPageParams) SalesQ {
	q.selector = params.ApplyTo(q.selector, "sales.id")
	return q
}

// Get - loads a row from `sales`
// returns nil, nil - if sale does not exists
// returns error if more than one Sale found
func (q SalesQ) Get() (*Sale, error) {
	var result Sale
	err := q.repo.Get(&result, q.selector)
	if err != nil {
		if q.repo.NoRows(err) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load sale")
	}

	return &result, nil
}

//Select - selects slice from the db, if no sales found - returns nil, nil
func (q SalesQ) Select() ([]Sale, error) {
	var result []Sale
	err := q.repo.Select(&result, q.selector)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load sales")
	}

	return result, nil
}
