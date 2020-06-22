package requests

import (
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/urlval"
	"net/http"
)

const (
	// IncludeTypeAssetPairListBaseAssets - defines if the base asset should be included in the response
	IncludeTypeAssetPairListBaseAssets = "base_asset"
	// IncludeTypeAssetPairListQuoteAssets - defines if the quote asset should be included in the response
	IncludeTypeAssetPairListQuoteAssets = "quote_asset"

	// FilterTypeAssetPairListBaseAsset - defines if we need to filter the list by base asset
	FilterTypeAssetPairListBaseAsset = "base_asset"
	// FilterTypeAssetPairListQuoteAsset - defines if we need to filter the list by quote asset
	FilterTypeAssetPairListQuoteAsset = "quote_asset"
	// FilterTypeAssetPairListAsset - defines if we need to filter the list by asset (no matter base or quote it is)
	FilterTypeAssetPairListAsset = "asset"
	// FilterTypeAssetPairListPolicy - defines if we need to filter the list by policy
	FilterTypeAssetPairListPolicy = "policy"
)

var includeTypeAssetPairListAll = map[string]struct{}{
	IncludeTypeAssetPairListBaseAssets:  {},
	IncludeTypeAssetPairListQuoteAssets: {},
}

var filterTypeAssetPairListAll = map[string]struct{}{
	FilterTypeAssetPairListBaseAsset:  {},
	FilterTypeAssetPairListQuoteAsset: {},
	FilterTypeAssetPairListAsset:      {},
	FilterTypeAssetPairListPolicy:     {},
}

// GetAssetPairList - represents params to be specified by user for getAssetPairList handler
type GetAssetPairList struct {
	*base
	Filters GetAssetPairListFilters
	PageParams *pgdb.OffsetPageParams
}
type GetAssetPairListFilters struct {
	Policy     []uint64 `filter:"policy"`
	Asset      []string `filter:"asset"`
	BaseAsset  []string `filter:"base_asset"`
	QuoteAsset []string `filter:"quote_asset"`

}
// NewGetAssetPairList returns new instance of GetAssetPairList request
func NewGetAssetPairList(r *http.Request) (*GetAssetPairList, error) {
	b, err := newBase(r, baseOpts{
		supportedIncludes: includeTypeAssetPairListAll,
		supportedFilters:  filterTypeAssetPairListAll,
	})
	if err != nil {
		return nil, err
	}

	pageParams, err := b.getOffsetBasedPageParams()
	if err != nil {
		return nil, err
	}

	request := GetAssetPairList{
		base:       b,
		PageParams: pageParams,
	}
  	request.Filters=
  		GetAssetPairListFilters{
  		Policy: []uint64{0},
		}
	err=urlval.Decode(r.URL.Query(),&request.Filters)


	return &request, nil
}
