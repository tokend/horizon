package requests

import (
	"net/http"
)

type GetUpdateLimitsRequestsFilter struct {
	GetRequestListBaseFilters
}

type GetUpdateLimitsRequests struct {
	*GetRequestsBase
	Filters GetUpdateLimitsRequestsFilter
}

func NewGetUpdateLimitsRequests(r *http.Request) (request GetUpdateLimitsRequests, err error) {
	request.GetRequestsBase, err = NewGetRequestsBase(
		r,
		&request.Filters,
		map[string]struct{}{},
		map[string]struct{}{},
	)
	if err != nil {
		return request, err
	}

	return request, nil
}
