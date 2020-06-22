package requests

import (
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/urlval"
	"net/http"
)

type GetRequestsBase struct {
	*base
	Filters    GetRequestListBaseFilters
	PageParams *pgdb.CursorPageParams
}

func NewGetRequestsBase(
	r *http.Request,
	filterDst interface{},
	filters map[string]struct{},
	includes map[string]struct{},
) (*GetRequestsBase, error) {

	// merge filters
	mergedFilters := map[string]struct{}{}
	for k := range filters {
		mergedFilters[k] = struct{}{}
	}
	for k := range filterTypeRequestListAll {
		mergedFilters[k] = struct{}{}
	}
	// merge includes
	mergedIncludes := map[string]struct{}{}
	for k := range includes {
		mergedIncludes[k] = struct{}{}
	}
	for k := range includeTypeReviewableRequestListAll {
		mergedIncludes[k] = struct{}{}
	}

	b, err := newBase(r, baseOpts{
		supportedFilters:  mergedFilters,
		supportedIncludes: mergedIncludes,
	})
	if err != nil {
		return nil, err
	}

	pageParams, err := b.getCursorBasedPageParams()
	if err != nil {
		return nil, err
	}


	err=urlval.Decode(r.URL.Query(), filterDst)

	var baseFilters=
	 	GetRequestListBaseFilters{
		ID: []uint64{0},
		Requestor: []string{""},
	 	Reviewer: []string{""},
	 	State: []uint64{0},
	 	Type: []uint64{0},
	 	PendingTasks: []uint64{0},
	 	PendingTasksAnyOf: []uint64{0},
	 	PendingTasksNotSet: []uint64{0},
	 	MissingPendingTasks: []uint64{0},
	 	}
	err=urlval.Decode(r.URL.Query(),&baseFilters)

	ID, err := b.getUint64ID()
	baseFilters.ID=[]uint64{ID}
	if err != nil {
		return nil, err
	}

	return &GetRequestsBase{
		base:       b,
		Filters:    baseFilters,
		PageParams: pageParams,
	}, nil
}
