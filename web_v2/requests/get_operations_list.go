package requests

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"net/http"

	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/urlval"
)

const (
	IncludeTypeOperationsListOperationDetails = "operation.details"

	FilterTypeOperationsListTypes = "types"
	FilterTypeOperationsSource    = "source"
)

type GetOperations struct {
	*base
	PageParams pgdb.CursorPageParams
	PageNumber *uint64 `page:"number"`
	Filters    struct {
		Types  []int   `filter:"types"`
		Source *string `filter:"source"`
	}
	Includes struct {
		OperationDetails bool `include:"operation.details"`
	}
}

func NewGetOperations(r *http.Request) (*GetOperations, error) {
	b, err := newBase(r, baseOpts{
		supportedIncludes: map[string]struct{}{
			IncludeTypeOperationsListOperationDetails: {},
		},
		supportedFilters: map[string]struct{}{
			FilterTypeOperationsListTypes: {},
			FilterTypeOperationsSource:    {},
		},
	})
	if err != nil {
		return nil, err
	}

	request := GetOperations{
		base: b,
	}

	err = urlval.DecodeSilently(r.URL.Query(), &request)
	if err != nil {
		return nil, err
	}

	err = SetDefaultCursorPageParams(&request.PageParams)
	if err != nil {
		return nil, err
	}

	// use part of cursor params struct to prevent decode same token twice
	if request.PageNumber != nil {
		params := pgdb.OffsetPageParams{
			Limit:      request.PageParams.Limit,
			Order:      request.PageParams.Order,
			PageNumber: *request.PageNumber,
		}

		err = request.SetDefaultOffsetPageParams(&params)
		if err != nil {
			return nil, errors.Wrap(err, "failed to set default page params")
		}

		request.PageParams.Limit = params.Limit
		request.PageParams.Order = params.Order
		request.PageNumber = &params.PageNumber
	}

	return &request, nil
}
