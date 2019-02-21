package requests

import (
	"net/http"
)

const (
	IncludeTypeCreateAtomicSwapBidRequestsBalance     = "request_details.base_balance"
	IncludeTypeCreateAtomicSwapBidRequestsQuoteAssets = "request_details.quote_assets"
)

var includeTypeCreateAtomicSwapBidRequests = map[string]struct{}{
	IncludeTypeCreateAtomicSwapBidRequestsBalance:     {},
	IncludeTypeCreateAtomicSwapBidRequestsQuoteAssets: {},
}

const (
	FilterTypeCreateAtomicSwapBidRequestsBalance = "request_details.base_balance"
)

var filterTypeCreateAtomicSwapBidRequests = map[string]struct{}{
	FilterTypeCreateAtomicSwapBidRequestsBalance: {},
}

type GetCreateAtomicSwapBidRequestsFilter struct {
	GetRequestListBaseFilters
	BaseBalance string `fig:"request_details.base_balance"`
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
