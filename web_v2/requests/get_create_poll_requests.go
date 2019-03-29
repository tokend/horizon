package requests

import (
	"net/http"
)

const (
	FilterTypeCreatePollRequestsPermissionType   = "request_details.permission_type"
	FilterTypeCreatePollRequestsVoteConfirmation = "request_details.vote_confirmation"
	FilterTypeCreatePollRequestsResultProvider   = "request_details.result_provider"
)

var filterTypeCreatePollRequests = map[string]struct{}{
	FilterTypeCreatePollRequestsPermissionType:   {},
	FilterTypeCreatePollRequestsVoteConfirmation: {},
	FilterTypeCreatePollRequestsResultProvider:   {},
}

type GetCreatePollRequestsFilter struct {
	GetRequestListBaseFilters
	PermissionType           uint32 `fig:"request_details.permission_type"`
	VoteConfirmationRequired bool   `fig:"request_details.vote_confirmation"`
	ResultProvider           string `fig:"request_details.result_provider"`
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
