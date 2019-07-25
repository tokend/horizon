package requests

import (
	"net/http"
)

const (
	IncludeTypeCreateAtomicSwapBidRequestsQuoteAsset = "request_details.quote_asset"
	IncludeTypeCreateAtomicSwapBidRequestsAsk = "request_details.ask"
)

var includeTypeCreateAtomicSwapBidRequests = map[string]struct{}{
	IncludeTypeCreateAtomicSwapBidRequestsQuoteAsset: {},
	IncludeTypeCreateAtomicSwapBidRequestsAsk: {},
}

const (
	FilterTypeCreateAtomicSwapBidRequestsQuoteAsset = "request_details.quote_asset"
	FilterTypeCreateAtomicSwapBidRequestsAskID = "request_details.ask_id"
)

var filterTypeCreateAtomicSwapBidRequests = map[string]struct{}{
	FilterTypeCreateAtomicSwapBidRequestsQuoteAsset: {},
	FilterTypeCreateAtomicSwapBidRequestsAskID: {},
}

type GetCreateAtomicSwapBidRequestsFilter struct {
	GetRequestListBaseFilters
	QuoteAsset string `fig:"request_details.quote_asset"`
	AskID uint64 `fig:"request_details.ask_id"`
}

type GetCreateAtomicSwapBidRequests struct {
	*GetRequestsBase
	Filters GetCreateAtomicSwapBidRequestsFilter
}

func NewGetCreateAtomicSwapBidRequests(r *http.Request) (request GetCreateAtomicSwapBidRequests, err error) {
	request.GetRequestsBase, err = NewGetRequestsBase(
		r,
		&request.Filters,
		filterTypeCreateAtomicSwapBidRequests,
		includeTypeCreateAtomicSwapBidRequests,
	)
	if err != nil {
		return request, err
	}

	return request, nil
}
