package requests

import (
	"net/http"
)

const (
	IncludeTypeKYCRecoveryRequestsAccount = "request_details.target_account"
)

var includeTypeKYCRecoveryRequests = map[string]struct{}{
	IncludeTypeKYCRecoveryRequestsAccount: {},
}

const (
	FilterTypeKYCRecoveryRequestsAccount = "request_details.target_account"
)

var filterTypeKYCRecoveryRequests = map[string]struct{}{
	FilterTypeKYCRecoveryRequestsAccount: {},
}

type GetKYCRecoveryRequestsFilter struct {
	GetRequestListBaseFilters
	Account []string `filter:"request_details.target_account"`
}

type GetKYCRecoveryRequests struct {
	*GetRequestsBase
	Filters GetKYCRecoveryRequestsFilter
}

func NewGetKYCRecoveryRequests(r *http.Request) (request GetKYCRecoveryRequests, err error) {
	request.Filters=
		GetKYCRecoveryRequestsFilter{
		Account: []string{},
		}
	request.GetRequestsBase, err = NewGetRequestsBase(
		r,
		&request.Filters,
		filterTypeKYCRecoveryRequests,
		includeTypeKYCRecoveryRequests,
	)
	if err != nil {
		return request, err
	}

	return request, nil
}
