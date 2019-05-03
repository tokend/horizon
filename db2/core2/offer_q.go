package core2

import (
	sq "github.com/lann/squirrel"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2"
)

// OffersQ is a helper struct to aid in configuring queries that loads offers
type OffersQ struct {
	repo     *db2.Repo
	selector sq.SelectBuilder
}

// NewOffersQ - creates new instance of OffersQ with no filters
func NewOffersQ(repo *db2.Repo) OffersQ {
	return OffersQ{
		repo: repo,
		selector: sq.Select(
			"offers.offer_id",
			"offers.owner_id",
			"offers.order_book_id",
			"offers.base_asset_code",
			"offers.quote_asset_code",
			"offers.base_balance_id",
			"offers.quote_balance_id",
			"offers.fee",
			"offers.is_buy",
			"offers.created_at",
			"offers.base_amount",
			"offers.quote_amount",
			"offers.price",
		).From("offer offers"),
	}
}

// FilterByBaseBalanceID - returns q with filter by base balance
func (q OffersQ) FilterByBaseBalanceID(id string) OffersQ {
	q.selector = q.selector.Where("offers.base_balance_id = ?", id)
	return q
}

// FilterByQuoteBalanceID - returns q with filter by quote balance
func (q OffersQ) FilterByQuoteBalanceID(id string) OffersQ {
	q.selector = q.selector.Where("offers.quote_balance_id = ?", id)
	return q
}

// FilterByBaseAssetCode - returns q with filter by base asset
func (q OffersQ) FilterByBaseAssetCode(code string) OffersQ {
	q.selector = q.selector.Where("offers.base_asset_code = ?", code)
	return q
}

// FilterByQuoteAssetCode - returns q with filter by quote asset
func (q OffersQ) FilterByQuoteAssetCode(code string) OffersQ {
	q.selector = q.selector.Where("offers.quote_asset_code = ?", code)
	return q
}

// FilterByOwnerID - returns q with filter by owner ID
func (q OffersQ) FilterByOwnerID(id string) OffersQ {
	q.selector = q.selector.Where("offers.owner_id = ?", id)
	return q
}

// FilterByOrderBookID - returns q with filter by order book ID
// use 0 - to get all offers from secondary market
// use -1 - to get all offers from primary market (all non zero order books)
// use saleID - to get offers of specified sale
func (q OffersQ) FilterByOrderBookID(id int64) OffersQ {
	if id < 0 {
		q.selector = q.selector.Where("offers.order_book_id <> ?", 0)
		return q
	}

	q.selector = q.selector.Where("offers.order_book_id = ?", id)
	return q
}

// FilterByIsBuy - returns q with filter by is_buy
func (q OffersQ) FilterByIsBuy(isBuy bool) OffersQ {
	q.selector = q.selector.Where("offers.is_buy = ?", isBuy)
	return q
}

// FilterByOfferID - returns q with filter by offer ID
func (q OffersQ) FilterByOfferID(id uint64) OffersQ {
	q.selector = q.selector.Where("offers.offer_id = ?", id)
	return q
}

// Page - returns Q with specified limit and offset params
func (q OffersQ) Page(params db2.OffsetPageParams) OffersQ {
	q.selector = params.ApplyTo(q.selector, "offers.offer_id")
	return q
}

// CursorPage - returns Q with specified limit and offset params
func (q OffersQ) CursorPage(params db2.CursorPageParams) OffersQ {
	q.selector = params.ApplyTo(q.selector, "offers.offer_id")
	return q
}

// GetByOfferID - loads a row from `offers` found with offer ID
// returns nil, nil - if such offer doesn't exist
func (q OffersQ) GetByOfferID(id uint64) (*Offer, error) {
	return q.FilterByOfferID(id).Get()
}

// WithBaseAsset - joins base asset
func (q OffersQ) WithBaseAsset() OffersQ {
	q.selector = q.selector.
		Columns(db2.GetColumnsForJoin(assetColumns, "base_assets")...).
		LeftJoin("asset base_assets ON offers.base_asset_code = base_assets.code")

	return q
}

// WithQuoteAsset - joins quote asset
func (q OffersQ) WithQuoteAsset() OffersQ {
	q.selector = q.selector.
		Columns(db2.GetColumnsForJoin(assetColumns, "quote_assets")...).
		LeftJoin("asset quote_assets ON offers.quote_asset_code = quote_assets.code")

	return q
}

// Get - loads a row from `asset_pairs`
// returns nil, nil - if asset pair does not exists
// returns error if more than one asset pair found
func (q OffersQ) Get() (*Offer, error) {
	var result Offer
	err := q.repo.Get(&result, q.selector)
	if err != nil {
		if q.repo.NoRows(err) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load offer")
	}

	return &result, nil
}

// Select - selects slice from the db, if no pairs found - returns nil, nil
func (q OffersQ) Select() ([]Offer, error) {
	var result []Offer
	err := q.repo.Select(&result, q.selector)
	if err != nil {
		if q.repo.NoRows(err) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load offers")
	}

	return result, nil
}
