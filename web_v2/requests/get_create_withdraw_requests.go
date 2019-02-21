package requests

import (
	"net/http"
)

const (
	IncludeTypeCreateWithdrawRequestsBalance = "request_details.balance"
)

var includeTypeCreateWithdrawRequests = map[string]struct{}{
	IncludeTypeCreateWithdrawRequestsBalance: {},
}

const (
	FilterTypeCreateWithdrawRequestsBalance = "request_details.balance"
)

var filterTypeCreateWithdrawRequests = map[string]struct{}{
	FilterTypeCreateWithdrawRequestsBalance: {},
}

type GetCreateWithdrawRequestsFilter struct {
	GetRequestListBaseFilters
	Balance string `fig:"request_details.balance"`
}

type GetCreateWithdrawRequests struct {
	*GetRequestsBase
	Filters GetCreateWithdrawRequestsFilter
}

func NewGetCreateWithdrawRequests(r *http.Request) (request GetCreateWithdrawRequests, err error) {
	request.GetRequestsBase, err = NewGetRequestsBase(
		r,
		&request.Filters,
		filterTypeCreateWithdrawRequests,
		includeTypeCreateWithdrawRequests,
	)
	if err != nil {
		return request, err
	}

	return request, nil
}
