package requests

import (
	"net/http"
)

type GetUpdateLimitsRequests struct {
	*GetRequestsBase
	Filters GetRequestListBaseFilters
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
