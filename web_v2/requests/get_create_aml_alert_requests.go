package requests

import (
	"net/http"
)

const (
	IncludeTypeCreateAmlAlertRequestsBalance = "request_details.balance"
)

var includeTypeCreateAmlAlertRequests = map[string]struct{}{
	IncludeTypeCreateAmlAlertRequestsBalance: {},
}

const (
	FilterTypeCreateAmlAlertRequestsBalance = "request_details.balance"
)

var filterTypeCreateAmlAlertRequests = map[string]struct{}{
	FilterTypeCreateAmlAlertRequestsBalance: {},
}

type GetCreateAmlAlertRequestsFilter struct {
	GetRequestListBaseFilters
	Balance *string `filter:"request_details.balance"`
}

type GetCreateAmlAlertRequests struct {
	*GetRequestsBase
	Filters GetCreateAmlAlertRequestsFilter
}

func NewGetCreateAmlAlertRequests(r *http.Request) (request GetCreateAmlAlertRequests, err error) {
	request.GetRequestsBase, err = NewGetRequestsBase(
		r,
		&request.Filters,
		filterTypeCreateAmlAlertRequests,
		includeTypeCreateAmlAlertRequests,
	)
	if err != nil {
		return request, err
	}

	return request, nil
}
