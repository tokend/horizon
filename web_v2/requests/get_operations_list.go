package requests

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
	"gitlab.com/tokend/horizon/db2"
)

const (
	IncludeTypeOperationsListOperationDetails = "operation.details"

	FilterTypeOperationsListTypes = "types"
)

type GetOperations struct {
	*base
	PageParams *db2.CursorPageParams
	Filters    struct {
		Types []int
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

	err = request.populateFilters()
	if err != nil {
		return nil, err
	}

	return &request, nil
}

func (r *GetOperations) getIntSlice(name string) ([]int, error) {
	valuesStr := strings.Split(r.getString(name), ",")

	if len(valuesStr) > 0 {
		valuesInt := make([]int, 0, len(valuesStr))
		for _, v := range valuesStr {
			if v != "" {
				valueInt, err := strconv.Atoi(v)
				if err != nil {
					return nil, validation.Errors{
						v: err,
					}
				}

				valuesInt = append(valuesInt, valueInt)
			}
		}

		return valuesInt, nil
	}

	return []int{}, nil
}

func (r *GetOperations) populateFilters() (err error) {
	r.Filters.Types, err = r.getIntSlice(
		fmt.Sprintf("filter[%s]", FilterTypeOperationsListTypes),
	)
	if err != nil {
		return err
	}

	return nil
}
