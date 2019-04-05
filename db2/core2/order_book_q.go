package core2

import (
	"fmt"
	sq "github.com/lann/squirrel"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2"
)

// OrderBooksQ is a helper struct to aid in configuring queries that loads order book entries
type OrderBooksQ struct {
	repo     *db2.Repo
	selector sq.SelectBuilder
}

// NewOrderBooksQ - creates new instance of OrderBooksQ with no filters
func NewOrderBooksQ(repo *db2.Repo) OrderBooksQ {
	subQuery := sq.Select(
		"format('%s:%s:%s:%s:%s', offers.quote_asset_code, offers.base_asset_code, "+
			"offers.order_book_id, offers.is_buy, offers.price) id",
		"offers.base_asset_code",
		"offers.quote_asset_code",
		"offers.order_book_id",
		"offers.is_buy",
		"offers.price",
		"sum(offers.base_amount) base_amount",
		"sum(offers.quote_amount) quote_amount",
		"to_timestamp(min(offers.created_at)) at time zone 'UTC' created_at",
		"to_timestamp(max(offers.created_at)) at time zone 'UTC' updated_at",
	).From("offer offers").GroupBy(
		"offers.base_asset_code",
		"offers.quote_asset_code",
		"offers.order_book_id",
		"offers.is_buy",
		"offers.price",
	)

	return OrderBooksQ{
		repo: repo,
		selector: sq.Select(
			"order_book_entries.id",
			"order_book_entries.base_asset_code",
			"order_book_entries.quote_asset_code",
			"order_book_entries.order_book_id",
			"order_book_entries.is_buy",
			"order_book_entries.price",
			"order_book_entries.base_amount",
			"order_book_entries.quote_amount",
			"order_book_entries.created_at",
			"order_book_entries.updated_at",
		).FromSelect(subQuery, "order_book_entries"),
	}
}

// WithBaseAsset - joins base asset
func (q OrderBooksQ) WithBaseAsset() OrderBooksQ {
	q.selector = q.selector.
		Columns(db2.GetColumnsForJoin(assetColumns, "base_assets")...).
		LeftJoin("asset base_assets ON order_book_entries.base_asset_code = base_assets.code")

	return q
}

// WithQuoteAsset - joins quote asset
func (q OrderBooksQ) WithQuoteAsset() OrderBooksQ {
	q.selector = q.selector.
		Columns(db2.GetColumnsForJoin(assetColumns, "quote_assets")...).
		LeftJoin("asset quote_assets ON order_book_entries.quote_asset_code = quote_assets.code")

	return q
}

// FilterByBaseAssetCode - returns q with filter by base asset
func (q OrderBooksQ) FilterByBaseAssetCode(code string) OrderBooksQ {
	q.selector = q.selector.Where("order_book_entries.base_asset_code = ?", code)
	return q
}

// FilterByQuoteAssetCode - returns q with filter by quote asset
func (q OrderBooksQ) FilterByQuoteAssetCode(code string) OrderBooksQ {
	q.selector = q.selector.Where("order_book_entries.quote_asset_code = ?", code)
	return q
}

// FilterByIsBuy - returns q with filter by is_buy
func (q OrderBooksQ) FilterByIsBuy(isBuy bool) OrderBooksQ {
	q.selector = q.selector.Where("order_book_entries.is_buy = ?", isBuy)
	return q
}

// FilterByOrderBookID - returns q with order by price
func (q OrderBooksQ) OrderByPrice(order string) OrderBooksQ {
	q.selector = q.selector.OrderBy(fmt.Sprintf("%s %s", "order_book_entries.price", order))
	return q
}

// FilterByOrderBookID - returns q with filter by order book ID
func (q OrderBooksQ) FilterByOrderBookID(id uint64) OrderBooksQ {
	q.selector = q.selector.Where("order_book_entries.order_book_id = ?", id)
	return q
}

// Page - returns Q with specified limit and offset params
func (q OrderBooksQ) Page(params db2.OffsetPageParams) OrderBooksQ {
	q.selector = params.ApplyTo(q.selector,
		"order_book_entries.price",
		"order_book_entries.base_asset_code",
		"order_book_entries.quote_asset_code",
		"order_book_entries.is_buy",
	)
	return q
}

// Select - selects slice from the db, if no order book entries found - returns nil, nil
func (q OrderBooksQ) Select() ([]OrderBookEntry, error) {
	var result []OrderBookEntry
	err := q.repo.Select(&result, q.selector)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load order book entries")
	}

	return result, nil
}
