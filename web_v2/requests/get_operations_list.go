package requests

import (
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/urlval"
	"net/http"
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

	request.Filters.Types=[]int{0}
	err=urlval.Decode(r.URL.Query(),&request.Filters)


	return &request, nil
}

//func (r *GetOperations) getIntSlice(name string) ([]int, error) {
//	valuesStr := strings.Split(r.getString(name), ",")
//
//	if len(valuesStr) > 0 {
//		valuesInt := make([]int, 0, len(valuesStr))
//		for _, v := range valuesStr {
//			if v != "" {
//				valueInt, err := strconv.Atoi(v)
//				if err != nil {
//					return nil, validation.Errors{
//						v: err,
//					}
//				}
//
//				valuesInt = append(valuesInt, valueInt)
//			}
//		}
//
//		return valuesInt, nil
//	}
//
//	return []int{}, nil
//}
//
//func (r *GetOperations) populateFilters() (err error) {
//	r.Filters.Types, err = r.getIntSlice(
//		fmt.Sprintf("filter[%s]", FilterTypeOperationsListTypes),
//	)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
