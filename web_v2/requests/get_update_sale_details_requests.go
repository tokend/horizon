package requests

import (
	"net/http"
)

const (
	IncludeTypeUpdateSaleDetailsRequestsSale = "request_details.sale"
)

var includeTypeUpdateSaleDetailsRequests = map[string]struct{}{
	IncludeTypeUpdateSaleDetailsRequestsSale: {},
}

type GetUpdateSaleDetailsRequests struct {
	*GetRequestsBase
	Filters GetRequestListBaseFilters
}

func NewGetUpdateSaleDetailsRequests(r *http.Request) (request GetUpdateSaleDetailsRequests, err error) {
	request.GetRequestsBase, err = NewGetRequestsBase(
		r,
		&request.Filters,
		map[string]struct{}{},
		includeTypeUpdateSaleDetailsRequests,
	)
	if err != nil {
		return request, err
	}

	return request, nil
}
