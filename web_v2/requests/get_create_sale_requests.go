package requests

import (
	"net/http"
)

const (
	FilterTypeCreateSaleRequestsBaseAsset         = "request_details.base_asset"
	FilterTypeCreateSaleRequestsDefaultQuoteAsset = "request_details.default_quote_asset"
)

const (
	IncludeTypeCreateSaleRequestsBaseAsset         = "request_details.base_asset"
	IncludeTypeCreateSaleRequestsDefaultQuoteAsset = "request_details.default_quote_asset"
	IncludeTypeCreateSaleRequestsQuoteAssets       = "request_details.quote_assets"
)

var filterTypeCreateSaleRequests = map[string]struct{}{
	FilterTypeCreateSaleRequestsBaseAsset:         {},
	FilterTypeCreateSaleRequestsDefaultQuoteAsset: {},
}

var includeTypeCreateSaleRequests = map[string]struct{}{
	IncludeTypeCreateSaleRequestsBaseAsset:         {},
	IncludeTypeCreateSaleRequestsDefaultQuoteAsset: {},
	IncludeTypeCreateSaleRequestsQuoteAssets:       {},
}

type GetCreateSaleRequestsFilter struct {
	GetRequestListBaseFilters
	BaseAsset         string `fig:"request_details.base_asset"`
	DefaultQuoteAsset string `fig:"request_details.default_quote_asset"`
}

type GetCreateSaleRequests struct {
	*GetRequestsBase
	Filters GetCreateSaleRequestsFilter
}

func NewGetCreateSaleRequests(r *http.Request) (request GetCreateSaleRequests, err error) {
	request.GetRequestsBase, err = NewGetRequestsBase(
		r,
		&request.Filters,
		filterTypeCreateSaleRequests,
		includeTypeCreateSaleRequests,
	)
	if err != nil {
		return request, err
	}

	return request, nil
}
