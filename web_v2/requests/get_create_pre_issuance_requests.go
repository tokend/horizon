package requests

import (
	"net/http"
)

const (
	FilterTypeCreatePreIssuanceRequestsAsset = "request_details.asset"
)

var filterTypeCreatePreIssuanceRequests = map[string]struct{}{
	FilterTypeCreatePreIssuanceRequestsAsset: {},
}

const (
	IncludeTypeCreatePreIssuanceRequestsAsset = "request_details.asset"
)

var includeTypeCreatePreIssuanceRequests = map[string]struct{}{
	IncludeTypeCreatePreIssuanceRequestsAsset: {},
}

type GetCreatePreIssuanceRequestsFilter struct {
	GetRequestListBaseFilters
	Asset []string `filter:"request_details.asset"`
}

type GetCreatePreIssuanceRequests struct {
	*GetRequestsBase
	Filters GetCreatePreIssuanceRequestsFilter
}

func NewGetCreatePreIssuanceRequests(r *http.Request) (request GetCreatePreIssuanceRequests, err error) {
	request.GetRequestsBase, err = NewGetRequestsBase(
		r,
		&request.Filters,
		filterTypeCreatePreIssuanceRequests,
		includeTypeCreatePreIssuanceRequests,
	)
	if err != nil {
		return request, err
	}

	return request, nil
}
