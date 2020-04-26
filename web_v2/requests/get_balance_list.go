package requests

import (
	"gitlab.com/tokend/horizon/bridge"
	"net/http"
)

const (
	// IncludeTypeBalanceListState - defines if the state of the balance should be included in the response
	IncludeTypeBalanceListState = "state"

	// IncludeTypeBalanceListOwner - defines if the owner of the balance should be included in the response
	IncludeTypeBalanceListOwner = "owner"

	// FilterTypeBalanceListAsset - defines if we need to filter the list by asset
	FilterTypeBalanceListAsset = "asset"

	// FilterTypeBalanceListOwner - defines if we need to filter the list by owner
	FilterTypeBalanceListOwner = "owner"

	// FilterTypeBalanceListAssetOwner - defines if we need to filter the list by asset owner
	FilterTypeBalanceListAssetOwner = "asset_owner"
)

var includeTypeBalanceListAll = map[string]struct{}{
	IncludeTypeBalanceListState: {},
	IncludeTypeBalanceListOwner: {},
}

var filterTypeBalanceListAll = map[string]struct{}{
	FilterTypeBalanceListAsset:      {},
	FilterTypeBalanceListAssetOwner: {},
	FilterTypeBalanceListOwner:      {},
}

// GetBalanceList - represents params to be specified by user for getBalanceList handler
type GetBalanceList struct {
	*base
	Filters struct {
		Asset      string `fig:"asset"`
		AssetOwner string `json:"asset_owner"`
		Owner      string `json:"owner"`
	}
	PageParams *bridge.OffsetPageParams
}

// NewGetBalanceList - returns new instance of GetBalanceList
func NewGetBalanceList(r *http.Request) (*GetBalanceList, error) {
	b, err := newBase(r, baseOpts{
		supportedIncludes: includeTypeBalanceListAll,
		supportedFilters:  filterTypeBalanceListAll,
	})
	if err != nil {
		return nil, err
	}

	pageParams, err := b.getOffsetBasedPageParams()
	if err != nil {
		return nil, err
	}

	request := GetBalanceList{
		base:       b,
		PageParams: pageParams,
	}

	err = b.populateFilters(&request.Filters)
	if err != nil {
		return nil, err
	}

	return &request, nil
}
