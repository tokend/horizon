package requests

import (
	"gitlab.com/tokend/horizon/db2"
	"net/http"
)

type GetRequestsBase struct {
	*base
	Filters    GetRequestListBaseFilters
	PageParams *db2.CursorPageParams
}

func NewGetRequestsBase(
	r *http.Request, filterDst interface{}, filters map[string]struct{}, includes map[string]struct{},
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

	err = b.populateFilters(filterDst)
	if err != nil {
		return nil, err
	}

	var baseFilters GetRequestListBaseFilters
	err = b.populateFilters(&baseFilters)
	if err != nil {
		return nil, err
	}
	baseFilters.ID, err = b.getUint64ID()
	if err != nil {
		return nil, err
	}

	return &GetRequestsBase{
		base:       b,
		Filters:    baseFilters,
		PageParams: pageParams,
	}, nil
}
