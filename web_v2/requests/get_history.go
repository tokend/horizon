package requests

import (
	"net/http"

	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/urlval"
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
	// FilterTypeEffect - defines if we need to filter the list by effect type
	FilterTypeHistoryEffectType = "effect"
	// FilterTypeHistoryIDs
	FilterTypeHistoryIDs = "id"
)

//GetHistory - represents params to be specified for Get History handler
type GetHistory struct {
	*base
	PageParams pgdb.CursorPageParams
	Filters    struct {
		Account    *string `filter:"account"`
		Balance    *string `filter:"balance"`
		Asset      *string `filter:"asset"`
		EffectType []int64 `filter:"effect"`
	}
	Includes struct {
		Operation        bool `include:"operation"`
		Effect           bool `include:"effect"`
		OperationDetails bool `include:"operation.details"`
		Asset            bool `include:"asset"`
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
			FilterTypeHistoryAccount:    {},
			FilterTypeHistoryBalance:    {},
			FilterTypeHistoryAsset:      {},
			FilterTypeHistoryEffectType: {},
		},
	})
	if err != nil {
		return nil, err
	}

	request := GetHistory{
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

	return &request, nil
}
