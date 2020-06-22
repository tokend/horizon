package requests

import (
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/urlval"
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
		Asset      []string `filter:"asset"`
		AssetOwner []string `filter:"asset_owner"`
		Owner      []string `filter:"owner"`
	}
	PageParams *pgdb.OffsetPageParams
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

	var pageParams pgdb.OffsetPageParams
	err=urlval.Decode(r.URL.Query(), &pageParams)
	request := GetBalanceList{
		base:       b,
		PageParams: &pageParams,
	}

  	request.Filters = struct {
		Asset      []string `filter:"asset"`
		AssetOwner []string `filter:"asset_owner"`
		Owner      []string `filter:"owner"`
	}{[]string{""},[]string{""},[]string{""}}

	err=urlval.Decode(r.URL.Query(), &request.Filters)

	return &request, nil
}
