package requests

import (
	"net/http"

	"gitlab.com/distributed_lab/urlval"
)

const (
	IncludeTypeUpdateSaleDetailsRequestsSale = "request_details.sale"
)

var includeTypeUpdateSaleDetailsRequests = map[string]struct{}{
	IncludeTypeUpdateSaleDetailsRequestsSale: {},
}

type GetUpdateSaleDetailsRequests struct {
	GetRequestsBase
}

func NewGetUpdateSaleDetailsRequests(r *http.Request) (request GetUpdateSaleDetailsRequests, err error) {
	request.GetRequestsBase, err = NewGetRequestsBase(
		r,
		map[string]struct{}{},
		includeTypeUpdateSaleDetailsRequests,
	)
	if err != nil {
		return request, err
	}

	err = urlval.Decode(r.URL.Query(), &request)
	if err != nil {
		return request, err
	}

	err = PopulateRequest(&request.GetRequestsBase)
	if err != nil {
		return request, err
	}

	return request, nil
}
