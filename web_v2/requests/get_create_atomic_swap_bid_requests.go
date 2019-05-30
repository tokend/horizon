package requests

import (
	"net/http"
)

const (
	IncludeTypeCreateAtomicSwapBidRequestsQuoteAsset = "request_details.quote_asset"
)

var includeTypeCreateAtomicSwapBidRequests = map[string]struct{}{
	IncludeTypeCreateAtomicSwapBidRequestsQuoteAsset: {},
}

const (
	FilterTypeCreateAtomicSwapBidRequestsQuoteAsset = "request_details.quote_asset"
)

var filterTypeCreateAtomicSwapBidRequests = map[string]struct{}{
	FilterTypeCreateAtomicSwapBidRequestsQuoteAsset: {},
}

type GetCreateAtomicSwapBidRequestsFilter struct {
	GetRequestListBaseFilters
	QuoteAsset string `fig:"request_details.quote_asset"`
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
