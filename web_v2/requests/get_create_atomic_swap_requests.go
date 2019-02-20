package requests

import (
	"net/http"
)

const (
	IncludeTypeCreateAtomicSwapRequestsQuoteAsset = "request_details.quote_asset"
)

var includeTypeCreateAtomicSwapRequests = map[string]struct{}{
	IncludeTypeCreateAtomicSwapRequestsQuoteAsset: {},
}

const (
	FilterTypeCreateAtomicSwapRequestsQuoteAsset = "request_details.quote_asset"
)

var filterTypeCreateAtomicSwapRequests = map[string]struct{}{
	FilterTypeCreateAtomicSwapRequestsQuoteAsset: {},
}

type GetCreateAtomicSwapRequestsFilter struct {
	GetRequestListBaseFilters
	QuoteAsset string `fig:"request_details.quote_asset"`
}

type GetCreateAtomicSwapRequests struct {
	*GetRequestsBase
	Filters GetCreateAtomicSwapRequestsFilter
}

func NewGetCreateAtomicSwapRequests(r *http.Request) (request GetCreateAtomicSwapRequests, err error) {
	request.GetRequestsBase, err = NewGetRequestsBase(
		r,
		&request.Filters,
		filterTypeCreateAtomicSwapRequests,
		includeTypeCreateAtomicSwapRequests,
	)
	if err != nil {
		return request, err
	}

	return request, nil
}
