package requests

import (
	"github.com/spf13/cast"
	"net/http"
)

type GetRequestForAccount struct {
	*base
	GetRequestListBaseFilters
}

func NewGetRequestForAccount(r *http.Request) (request *GetRequests, err error) {
	b, err := newBase(r, baseOpts{
		supportedIncludes: includeTypeReviewableRequestListAll,
		supportedFilters:  filterTypeRequestListAll,
	})

	if err != nil {
		return nil, err
	}

	accountID, err := newAccountAddress(b, "id")
	if err != nil {
		return nil, err
	}

	requestID, err := b.getUint64("request_id")
	if err != nil {
		return nil, err
	}

	return newGetRequestForAccount(b, accountID, requestID)
}

func newGetRequestForAccount(b *base, accountID string, requestID uint64) (*GetRequests, error) {
	filters := GetRequestListBaseFilters{
		ID:        requestID,
		Requestor: accountID,
	}

	populateFilters(b, filters)

	page, err := b.getCursorBasedPageParams()
	if err != nil {
		return nil, err
	}

	return &GetRequests{
		Filters: filters,
		GetRequestsBase: &GetRequestsBase{
			base:       b,
			Filters:    filters,
			PageParams: page,
		},
	}, nil
}

func populateFilters(b *base, filters GetRequestListBaseFilters) {
	b.filter = map[string]string{
		"id":        cast.ToString(filters.ID),
		"requestor": filters.Requestor,
	}
}
