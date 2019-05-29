package requests

import (
	"net/http"
)

const (
	IncludeTypeCreateWithdrawRequestsBalance = "request_details.balance"
	IncludeTypeCreateWithdrawRequestsAsset   = "request_details.asset"
)

var includeTypeCreateWithdrawRequests = map[string]struct{}{
	IncludeTypeCreateWithdrawRequestsBalance: {},
	IncludeTypeCreateWithdrawRequestsAsset:   {},
}

const (
	FilterTypeCreateWithdrawRequestsBalance = "request_details.balance"
	FilterTypeCreateWithdrawRequestsAsset   = "request_details.asset"
)

var filterTypeCreateWithdrawRequests = map[string]struct{}{
	FilterTypeCreateWithdrawRequestsBalance: {},
	FilterTypeCreateWithdrawRequestsAsset:   {},
}

type GetCreateWithdrawRequestsFilter struct {
	GetRequestListBaseFilters
	Balance string `fig:"request_details.balance"`
	Asset   string `fig:"request_details.asset"`
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
