package requests

import (
	"net/http"
)

const (
	IncludeTypeCreateAtomicSwapAskRequestsBalance     = "request_details.base_balance"
	IncludeTypeCreateAtomicSwapAskRequestsQuoteAssets = "request_details.quote_assets"
)

var includeTypeCreateAtomicSwapAskRequests = map[string]struct{}{
	IncludeTypeCreateAtomicSwapAskRequestsBalance:     {},
	IncludeTypeCreateAtomicSwapAskRequestsQuoteAssets: {},
}

const (
	FilterTypeCreateAtomicSwapAskRequestsBalance = "request_details.base_balance"
)

var filterTypeCreateAtomicSwapAskRequests = map[string]struct{}{
	FilterTypeCreateAtomicSwapAskRequestsBalance: {},
}

type GetCreateAtomicSwapAskRequestsFilter struct {
	GetRequestListBaseFilters
	BaseBalance string `fig:"request_details.base_balance"`
}

type GetCreateAtomicSwapAskRequests struct {
	*GetRequestsBase
	Filters GetCreateAtomicSwapAskRequestsFilter
}

func NewGetCreateAtomicSwapAskRequests(r *http.Request) (request GetCreateAtomicSwapAskRequests, err error) {
	request.GetRequestsBase, err = NewGetRequestsBase(
		r,
		&request.Filters,
		filterTypeCreateAtomicSwapAskRequests,
		includeTypeCreateAtomicSwapAskRequests,
	)
	if err != nil {
		return request, err
	}

	return request, nil
}
