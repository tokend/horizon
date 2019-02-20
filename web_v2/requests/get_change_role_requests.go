package requests

import (
	"net/http"
)

const (
	IncludeTypeChangeRoleRequestsAccount = "request_details.account"
)

var includeTypeChangeRoleRequests = map[string]struct{}{
	IncludeTypeChangeRoleRequestsAccount: {},
}

const (
	FilterTypeChangeRoleRequestsAccount          = "request_details.destination_account"
	FilterTypeChangeRoleRequestsAccountRoleToSet = "request_details.account_role_to_set"
)

var filterTypeChangeRoleRequests = map[string]struct{}{
	FilterTypeChangeRoleRequestsAccount:          {},
	FilterTypeChangeRoleRequestsAccountRoleToSet: {},
}

type GetChangeRoleRequestsFilter struct {
	GetRequestListBaseFilters
	Account     string `fig:"request_details.destination_account"`
	AccountRole int32  `fig:"request_details.account_role_to_set"`
}

type GetChangeRoleRequests struct {
	*GetRequestsBase
	Filters GetChangeRoleRequestsFilter
}

func NewGetChangeRoleRequests(r *http.Request) (request GetChangeRoleRequests, err error) {
	request.GetRequestsBase, err = NewGetRequestsBase(
		r,
		&request.Filters,
		filterTypeChangeRoleRequests,
		includeTypeChangeRoleRequests,
	)
	if err != nil {
		return request, err
	}

	return request, nil
}
