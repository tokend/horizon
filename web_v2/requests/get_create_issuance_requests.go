package requests

import (
	"net/http"
)

const (
	FilterTypeCreateIssuanceRequestsAsset    = "request_details.asset"
	FilterTypeCreateIssuanceRequestsReceiver = "request_details.receiver"
)

var filterTypeCreateIssuanceRequests = map[string]struct{}{
	FilterTypeCreateIssuanceRequestsAsset:    {},
	FilterTypeCreateIssuanceRequestsReceiver: {},
}

const (
	IncludeTypeCreateIssuanceRequestsAsset = "request_details.asset"
)

var includeTypeCreateIssuanceRequests = map[string]struct{}{
	IncludeTypeCreateIssuanceRequestsAsset: {},
}

type GetCreateIssuanceRequestsFilter struct {
	GetRequestListBaseFilters
	Asset    string `fig:"request_details.asset"`
	Receiver string `fig:"request_details.receiver"`
}

type GetCreateIssuanceRequests struct {
	*GetRequestsBase
	Filters GetCreateIssuanceRequestsFilter
}

func NewGetCreateIssuanceRequests(r *http.Request) (request GetCreateIssuanceRequests, err error) {
	request.GetRequestsBase, err = NewGetRequestsBase(
		r,
		&request.Filters,
		filterTypeCreateIssuanceRequests,
		includeTypeCreateIssuanceRequests,
	)
	if err != nil {
		return request, err
	}

	return request, nil
}
