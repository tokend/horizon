package requests

import (
	"net/http"
)

const (
	FilterTypeUpdateAssetRequestsAsset = "request_details.asset"
)

var filterTypeUpdateAssetRequests = map[string]struct{}{
	FilterTypeUpdateAssetRequestsAsset: {},
}

const (
	IncludeTypeUpdateAssetRequestsAsset = "request_details.asset"
)

var includeTypeUpdateAssetRequests = map[string]struct{}{
	IncludeTypeUpdateAssetRequestsAsset: {},
}

type GetUpdateAssetRequestsFilter struct {
	GetRequestListBaseFilters
	Asset *string `filter:"request_details.asset"`
}

type GetUpdateAssetRequests struct {
	*GetRequestsBase
	Filters GetUpdateAssetRequestsFilter
}

func NewGetUpdateAssetRequests(r *http.Request) (request GetUpdateAssetRequests, err error) {
	request.GetRequestsBase, err = NewGetRequestsBase(
		r,
		&request.Filters,
		filterTypeUpdateAssetRequests,
		includeTypeUpdateAssetRequests,
	)
	if err != nil {
		return request, err
	}

	return request, nil
}
