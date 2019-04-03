package requests

import (
	"github.com/go-ozzo/ozzo-validation"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"net/http"
)

const (
	// IncludeTypeOrderBookV4OrderBookEntries
	IncludeTypeOrderBookV4OrderBookEntries = "order_book_entries"
	// IncludeTypeOrderBookV4BaseAssets - defines if base assets should be included in the response
	IncludeTypeOrderBookV4BaseAssets = "order_book_entries.base_asset"
	// IncludeTypeOrderBookV4QuoteAssets = defines if quote assets should be included in the response
	IncludeTypeOrderBookV4QuoteAssets = "order_book_entries.quote_asset"
)

var includeTypeOrderBookV4All = map[string]struct{}{
	IncludeTypeOrderBookV4OrderBookEntries: {},
	IncludeTypeOrderBookV4BaseAssets:       {},
	IncludeTypeOrderBookV4QuoteAssets:      {},
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
