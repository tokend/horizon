package requests

import (
	"gitlab.com/tokend/horizon/bridge"
	"net/http"
)

const (
	// IncludeTypeHistoryOperation - defines if the operation should be included in the response
	IncludeTypeHistoryOperation = "operation"
	// IncludeTypeHistoryEffect - defines if particular effect should be included
	IncludeTypeHistoryEffect = "effect"
	//IncludeTypeHistoryOperationDetails - defines if the operation details should be included
	IncludeTypeHistoryOperationDetails = "operation.details"
	//IncludeTypeHistoryAsset - defines if the asset should be included
	IncludeTypeHistoryAsset = "asset"

	// FilterTypeHistoryAccount - defines if we need to filter the list by participant account address
	FilterTypeHistoryAccount = "account"
	// FilterTypeHistoryBalance - defines if we need to filter the list by participating balance
	FilterTypeHistoryBalance = "balance"
	// FilterTypeHistoryAsset - defines if we need to filter the list by asset
	FilterTypeHistoryAsset = "asset"
	// FilterTypeHistoryIDs
	FilterTypeHistoryIDs = "id"
)

//GetHistory - represents params to be specified for Get History handler
type GetHistory struct {
	*base
	PageParams *bridge.CursorPageParams
	Filters    struct {
		Account string `fig:"account"`
		Balance string `fig:"balance"`
		Asset   string `fig:"asset"`
	}
}

// NewGetHistory returns the new instance of GetHistory request
func NewGetHistory(r *http.Request) (*GetHistory, error) {
	b, err := newBase(r, baseOpts{
		supportedIncludes: map[string]struct{}{
			IncludeTypeHistoryOperation:        {},
			IncludeTypeHistoryEffect:           {},
			IncludeTypeHistoryOperationDetails: {},
			IncludeTypeHistoryAsset:            {},
		},
		supportedFilters: map[string]struct{}{
			FilterTypeHistoryAccount: {},
			FilterTypeHistoryBalance: {},
			FilterTypeHistoryAsset:   {},
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
