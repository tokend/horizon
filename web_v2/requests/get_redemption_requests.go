package requests

import (
	"net/http"
)

const (
	FilterTypeRedemptionRequestsSourceBalance      = "request_details.source_balance"
	FilterTypeRedemptionRequestsDestinationAccount = "request_details.dest_account"
)

var filterTypeRedemptionRequests = map[string]struct{}{
	FilterTypeRedemptionRequestsSourceBalance:      {},
	FilterTypeRedemptionRequestsDestinationAccount: {},
}

const (
	IncludeTypeRedemptionRequestsSourceBalance      = "request_details.source_balance"
	IncludeTypeRedemptionRequestsDestinationAccount = "request_details.dest_account"
)

var includeTypeRedemptionRequests = map[string]struct{}{
	IncludeTypeRedemptionRequestsSourceBalance:      {},
	IncludeTypeRedemptionRequestsDestinationAccount: {},
}

type GetRedemptionRequestsFilter struct {
	GetRequestListBaseFilters
	SourceBalance      *string `filter:"request_details.source_balance"`
	DestinationAccount *string `filter:"request_details.dest_account"`
}

type GetRedemptionRequests struct {
	*GetRequestsBase
	Filters GetRedemptionRequestsFilter
}

func NewGetRedemptionRequests(r *http.Request) (request GetRedemptionRequests, err error) {
	request.GetRequestsBase, err = NewGetRequestsBase(
		r,
		&request.Filters,
		filterTypeRedemptionRequests,
		includeTypeRedemptionRequests,
	)
	if err != nil {
		return request, err
	}

	return request, nil
}
