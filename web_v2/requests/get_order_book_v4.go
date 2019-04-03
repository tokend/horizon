package requests

import (
	"github.com/go-ozzo/ozzo-validation"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"net/http"
)

const (
	// IncludeTypeOrderBookV4BuyEntries - defines if buy entries should be included in the response
	IncludeTypeOrderBookV4BuyEntries = "buy_entries"
	// IncludeTypeOrderBookV4BuyEntriesBaseAssets - defines if base asset of buy entries should be included in the response
	IncludeTypeOrderBookV4BuyEntriesBaseAssets = "buy_entries.base_asset"
	// IncludeTypeOrderBookV4BuyEntriesQuoteAssets - defines if quote asset of buy entries should be included in the response
	IncludeTypeOrderBookV4BuyEntriesQuoteAssets = "buy_entries.quote_asset"

	// IncludeTypeOrderBookV4SellEntries - defines if sell entries should be included in the response
	IncludeTypeOrderBookV4SellEntries = "sell_entries"
	// IncludeTypeOrderBookV4SellEntriesBaseAsset - defines if base asset of sell entries should be included in the response
	IncludeTypeOrderBookV4SellEntriesBaseAssets = "sell_entries.base_asset"
	// IncludeTypeOrderBookV4SellEntriesQuoteAsset - defines if quote asset of sell entries should be included in the response
	IncludeTypeOrderBookV4SellEntriesQuoteAssets = "sell_entries.quote_asset"
)

var includeTypeOrderBookV4All = map[string]struct{}{
	IncludeTypeOrderBookV4BuyEntries:             {},
	IncludeTypeOrderBookV4BuyEntriesBaseAssets:   {},
	IncludeTypeOrderBookV4BuyEntriesQuoteAssets:  {},
	IncludeTypeOrderBookV4SellEntries:            {},
	IncludeTypeOrderBookV4SellEntriesBaseAssets:  {},
	IncludeTypeOrderBookV4SellEntriesQuoteAssets: {},
}

// GetOrderBookV4 represents params to be specified by user for getOrderBookV4 handler
type GetOrderBookV4 struct {
	*base

	BaseAsset   string
	QuoteAsset  string
	OrderBookID uint64
}

// NewGetOrderBookV4 - returns new instance of GetOrderBook
func NewGetOrderBookV4(r *http.Request) (*GetOrderBookV4, error) {
	b, err := newBase(r, baseOpts{
		supportedIncludes: includeTypeOrderBookV4All,
	})
	if err != nil {
		return nil, err
	}

	baseAsset := b.getString("base")
	if baseAsset == "" {
		return nil, validation.Errors{
			"id": errors.New("id format is: 'base:quote:order_book_id', but base is empty"),
		}
	}
	quoteAsset := b.getString("quote")
	if baseAsset == "" {
		return nil, validation.Errors{
			"id": errors.New("id format is: 'base:quote:order_book_id', but quote is empty"),
		}
	}
	orderBookID, err := b.getUint64("order_book_id")
	if err != nil {
		return nil, validation.Errors{
			"id":   errors.New("id format is: 'base:quote:orde_-book_id', but order_book_id is invalid"),
			"meta": err,
		}
	}

	request := GetOrderBookV4{
		base:        b,
		BaseAsset:   baseAsset,
		QuoteAsset:  quoteAsset,
		OrderBookID: orderBookID,
	}

	return &request, nil
}
