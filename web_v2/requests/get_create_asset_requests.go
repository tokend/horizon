package requests

import (
	"net/http"
)

const (
	FilterTypeCreateAssetRequestsAsset = "request_details.asset"
)

var filterTypeCreateAssetRequests = map[string]struct{}{
	FilterTypeCreateAssetRequestsAsset: {},
}

const (
	IncludeTypeCreateAssetRequestsAsset = "request_details.asset"
)

var includeTypeCreateAssetRequests = map[string]struct{}{
	IncludeTypeCreateAssetRequestsAsset: {},
}

type GetCreateAssetRequestsFilter struct {
	GetRequestListBaseFilters
	Asset string `fig:"request_details.asset"`
}

type GetCreateAssetRequests struct {
	*GetRequestsBase
	Filters GetCreateAssetRequestsFilter
}

func NewGetCreateAssetRequests(r *http.Request) (request GetCreateAssetRequests, err error) {
	request.GetRequestsBase, err = NewGetRequestsBase(
		r,
		&request.Filters,
		filterTypeCreateAssetRequests,
		includeTypeCreateAssetRequests,
	)
	if err != nil {
		return request, err
	}

	return request, nil
}
