package requests

import (
	"gitlab.com/tokend/horizon/db2"
	"net/http"
)

const (
	// IncludeTypeOrderBookBaseAssets - defines if the base assets should be included in the response
	IncludeTypeOrderBookBaseAssets = "base_asset"
	// IncludeTypeOrderBookQuoteAssets = "quote_asset" - defines if the quote assets should be included in the response
	IncludeTypeOrderBookQuoteAssets = "quote_asset"

	// FilterTypeOrderBookBaseAsset - defines if we need to filter the list by base asset
	FilterTypeOrderBookBaseAsset = "base_asset"
	// FilterTypeOrderBookQuoteAsset - defines if we need to filter the list by quote asset
	FilterTypeOrderBookQuoteAsset = "quote_asset"
	// FilterTypeOrderBookIsBuy - defines if we need to filter the list by is buy
	FilterTypeOrderBookIsBuy = "is_buy"
)

var includeTypeOrderBookAll = map[string]struct{}{
	IncludeTypeOrderBookBaseAssets:  {},
	IncludeTypeOrderBookQuoteAssets: {},
}

var filterTypeOrderBookAll = map[string]struct{}{
	FilterTypeOrderBookBaseAsset:  {},
	FilterTypeOrderBookQuoteAsset: {},
	FilterTypeOrderBookIsBuy:      {},
}

// GetOrderBook represents params to be specified by user for getOfferList handler
type GetOrderBook struct {
	*base
	ID      uint64
	Filters struct {
		BaseAsset  string `fig:"base_asset"`
		QuoteAsset string `fig:"quote_asset"`
		IsBuy      bool   `fig:"is_buy,required"`
	}

	PageParams *db2.OffsetPageParams
}

// NewGetOrderBook - returns new instance of GetOrderBook
func NewGetOrderBook(r *http.Request) (*GetOrderBook, error) {
	b, err := newBase(r, baseOpts{
		supportedIncludes: includeTypeOrderBookAll,
		supportedFilters:  filterTypeOrderBookAll,
	})
	if err != nil {
		return nil, err
	}

	pageParams, err := b.getOffsetBasedPageParams()
	if err != nil {
		return nil, err
	}

	id, err := b.getUint64("id")
	if err != nil {
		return nil, err
	}

	request := GetOrderBook{
		base:       b,
		PageParams: pageParams,
		ID: id,
	}

	err = b.populateFilters(&request.Filters)
	if err != nil {
		return nil, err
	}

	return &request, nil
}
