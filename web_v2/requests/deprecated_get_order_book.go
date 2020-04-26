package requests

import (
	"gitlab.com/tokend/horizon/bridge"
	"net/http"
)

const (
	// DeprecatedIncludeTypeOrderBookBaseAssets - defines if the base assets should be included in the response
	DeprecatedIncludeTypeOrderBookBaseAssets = "base_asset"
	// DeprecatedIncludeTypeOrderBookQuoteAssets = "quote_asset" - defines if the quote assets should be included in the response
	DeprecatedIncludeTypeOrderBookQuoteAssets = "quote_asset"

	// DeprecatedFilterTypeOrderBookBaseAsset - defines if we need to filter the list by base asset
	DeprecatedFilterTypeOrderBookBaseAsset = "base_asset"
	// DeprecatedFilterTypeOrderBookQuoteAsset - defines if we need to filter the list by quote asset
	DeprecatedFilterTypeOrderBookQuoteAsset = "quote_asset"
	// DeprecatedFilterTypeOrderBookIsBuy - defines if we need to filter the list by is buy
	DeprecatedFilterTypeOrderBookIsBuy = "is_buy"
)

var deprecatedIncludeTypeOrderBookAll = map[string]struct{}{
	DeprecatedIncludeTypeOrderBookBaseAssets:  {},
	DeprecatedIncludeTypeOrderBookQuoteAssets: {},
}

var deprecatedFilterTypeOrderBookAll = map[string]struct{}{
	DeprecatedFilterTypeOrderBookBaseAsset:  {},
	DeprecatedFilterTypeOrderBookQuoteAsset: {},
	DeprecatedFilterTypeOrderBookIsBuy:      {},
}

// DeprecatedGetOrderBook represents params to be specified by user for getOfferList handler
type DeprecatedGetOrderBook struct {
	*base
	ID      uint64
	Filters struct {
		BaseAsset  string `fig:"base_asset"`
		QuoteAsset string `fig:"quote_asset"`
		IsBuy      bool   `fig:"is_buy"`
	}

	PageParams *bridge.OffsetPageParams
}

// NewDeprecatedGetOrderBook - returns new instance of DeprecatedGetOrderBook
func NewDeprecatedGetOrderBook(r *http.Request) (*DeprecatedGetOrderBook, error) {
	b, err := newBase(r, baseOpts{
		supportedIncludes: deprecatedIncludeTypeOrderBookAll,
		supportedFilters:  deprecatedFilterTypeOrderBookAll,
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

	request := DeprecatedGetOrderBook{
		base:       b,
		PageParams: pageParams,
		ID:         id,
	}

	err = b.populateFilters(&request.Filters)
	if err != nil {
		return nil, err
	}

	return &request, nil
}
