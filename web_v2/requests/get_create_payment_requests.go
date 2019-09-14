package requests

import (
	"net/http"
)

type GetCreatePaymentRequestsFilter struct {
	GetRequestListBaseFilters
}

type GetCreatePaymentRequests struct {
	*GetRequestsBase
	Filters GetCreatePaymentRequestsFilter
}

func NewGetCreatePaymentRequests(r *http.Request) (request GetCreatePaymentRequests, err error) {
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
