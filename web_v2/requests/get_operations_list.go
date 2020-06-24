package requests

import (
	"net/http"

	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/urlval"
)

const (
	IncludeTypeOperationsListOperationDetails = "operation.details"

	FilterTypeOperationsListTypes = "types"
)

type GetOperations struct {
	*base
	PageParams *pgdb.CursorPageParams
	Filters    struct {
		Types []int `filter:"types"`
	}
}

func NewGetOperations(r *http.Request) (*GetOperations, error) {
	b, err := newBase(r, baseOpts{
		supportedIncludes: map[string]struct{}{
			IncludeTypeOperationsListOperationDetails: {},
		},
		supportedFilters: map[string]struct{}{
			FilterTypeOperationsListTypes: {},
		},
	})
	if err != nil {
		return nil, err
	}

	pagingParams, err := b.getCursorBasedPageParams()
	if err != nil {
		return nil, err
	}

	request := GetOperations{
		base:       b,
		PageParams: pagingParams,
	}

	err = urlval.Decode(r.URL.Query(), &request.Filters)

	return &request, nil
}
