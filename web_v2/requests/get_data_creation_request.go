package requests

import (
	"net/http"

	"gitlab.com/distributed_lab/urlval"
)

type GetDataCreationRequests struct {
	GetRequestsBase
}

func NewGetDataCreationRequests(r *http.Request) (request GetDataCreationRequests, err error) {
	request.GetRequestsBase, err = NewGetRequestsBase(
		r,
		map[string]struct{}{},
		map[string]struct{}{},
	)
	if err != nil {
		return request, err
	}

	err = urlval.DecodeSilently(r.URL.Query(), &request)
	if err != nil {
		return request, err
	}

	err = PopulateRequest(&request.GetRequestsBase)
	if err != nil {
		return request, err
	}

	return request, nil
}
