package requests

import (
	"net/http"
)

const (
	FilterTypeCreatePollRequestsPermissionType           = "request_details.permission_type"
	FilterTypeCreatePollRequestsVoteConfirmationRequired = "request_details.vote_confirmation_required"
	FilterTypeCreatePollRequestsResultProvider           = "request_details.result_provider"
)

var filterTypeCreatePollRequests = map[string]struct{}{
	FilterTypeCreatePollRequestsPermissionType:           {},
	FilterTypeCreatePollRequestsVoteConfirmationRequired: {},
	FilterTypeCreatePollRequestsResultProvider:           {},
}

type GetCreatePollRequestsFilter struct {
	GetRequestListBaseFilters
	PermissionType           *uint32 `filter:"request_details.permission_type"`
	VoteConfirmationRequired *bool   `filter:"request_details.vote_confirmation_required"`
	ResultProvider           *string `filter:"request_details.result_provider"`
}

type GetCreatePollRequests struct {
	*GetRequestsBase
	Filters GetCreatePollRequestsFilter
}

func NewGetCreatePollRequests(r *http.Request) (request GetCreatePollRequests, err error) {
	request.GetRequestsBase, err = NewGetRequestsBase(
		r,
		&request.Filters,
		filterTypeCreatePollRequests,
		map[string]struct{}{},
	)
	if err != nil {
		return request, err
	}

	return request, nil
}
