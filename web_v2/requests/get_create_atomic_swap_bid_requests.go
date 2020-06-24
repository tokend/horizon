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
	FilterTypeCreateAtomicSwapBidRequestsAskID      = "request_details.ask_id"
	FilterTypeCreateAtomicSwapBidRequestsAskOwner   = "request_details.ask_owner"
)

var filterTypeCreateAtomicSwapBidRequests = map[string]struct{}{
	FilterTypeCreateAtomicSwapBidRequestsQuoteAsset: {},
	FilterTypeCreateAtomicSwapBidRequestsAskID:      {},
	FilterTypeCreateAtomicSwapBidRequestsAskOwner:   {},
}

type GetCreateAtomicSwapBidRequestsFilter struct {
	GetRequestListBaseFilters
	QuoteAsset *string `filter:"request_details.quote_asset"`
	AskID      *uint64 `filter:"request_details.ask_id"`
	AskOwner   *string `filter:"request_details.ask_owner"`
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
