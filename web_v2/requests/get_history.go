package requests

import (
	"net/http"

	"gitlab.com/tokend/horizon/db2"
)

const (
	// IncludeTypeHistoryOperation - defines if the operation should be included in the response
	IncludeTypeHistoryOperation = "operation"

	// FilterTypeHistoryAccount - defines if we need to filter the list by participant account address
	FilterTypeHistoryAccount = "account"
	// FilterTypeHistoryBalance - defines if we need to filter the list by participating balance
	FilterTypeHistoryBalance = "balance"
)

//GetHistory - represents params to be specified for Get History handler
type GetHistory struct {
	*base
	PageParams *db2.CursorPageParams
	Filters    struct {
		Account string `fig:"account"`
		Balance string `fig:"balance"`
	}
}

// NewGetHistory returns the new instance of GetHistory request
func NewGetHistory(r *http.Request) (*GetHistory, error) {
	b, err := newBase(r, baseOpts{
		supportedIncludes: map[string]struct{}{
			IncludeTypeHistoryOperation: {},
		},
		supportedFilters: map[string]struct{}{
			FilterTypeHistoryAccount: {},
			FilterTypeHistoryBalance: {},
		},
	})
	if err != nil {
		return nil, err
	}

	pagingParams, err := b.getCursorBasedPageParams()
	if err != nil {
		return nil, err
	}

	request := GetHistory{
		base:       b,
		PageParams: pagingParams,
	}

	err = b.populateFilters(&request.Filters)
	if err != nil {
		return nil, err
	}

	return &request, nil
}
