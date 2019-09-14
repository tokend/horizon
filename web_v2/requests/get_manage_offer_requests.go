package requests

import (
	"net/http"
)

type GetManageOfferRequestsFilter struct {
	GetRequestListBaseFilters
}

type GetManageOfferRequests struct {
	*GetRequestsBase
	Filters GetManageOfferRequestsFilter
}

func NewGetManageOfferRequests(r *http.Request) (request GetManageOfferRequests, err error) {
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
