package requests

import (
	"net/http"

	"gitlab.com/tokend/horizon/db2"
)

const (
	IncludeTypeBidListBaseBalances = "base_balance"
	IncludeTypeBidListOwners       = "owner"
	IncludeTypeBidListBaseAssets   = "base_asset"
	IncludeTypeBidListQuoteAssets  = "quote_assets"

	FilterTypeBidListOwner       = "owner"
	FilterTypeBidListBaseAsset   = "base_asset"
	FilterTypeBidListQuoteAssets = "quote_assets"
)

var includeTypeBidListAll = map[string]struct{}{
	IncludeTypeBidListBaseBalances: {},
	IncludeTypeBidListOwners:       {},
	IncludeTypeBidListBaseAssets:   {},
	IncludeTypeBidListQuoteAssets:  {},
}

var filterTypeBidListAll = map[string]struct{}{
	FilterTypeBidListOwner:       {},
	FilterTypeBidListBaseAsset:   {},
	FilterTypeBidListQuoteAssets: {},
}

// GetAtomicSwapBidList - represents params to be specified by user for Get AtomicSwapBidList handler
type GetAtomicSwapBidList struct {
	*base
	Filters struct {
		Owner       string   `fig:"owner"`
		BaseAsset   string   `fig:"base_asset"`
		QuoteAssets []string `fig:"quote_assets"`
	}
	PageParams *db2.OffsetPageParams
}

// NewGetAtomicSwapBidList returns new instance of GetAtomicSwapBidList request
func NewGetAtomicSwapBidList(r *http.Request) (*GetAtomicSwapBidList, error) {
	b, err := newBase(r, baseOpts{
		supportedIncludes: includeTypeBidListAll,
		supportedFilters:  filterTypeBidListAll,
	})
	if err != nil {
		return nil, err
	}

	// bid relations has not asset relation, we use balance relation
	if _, ok := b.include[IncludeTypeBidBaseAsset]; ok {
		b.include[IncludeTypeBidBaseBalance] = struct{}{}
	}

	pageParams, err := b.getOffsetBasedPageParams()
	if err != nil {
		return nil, err
	}

	request := GetAtomicSwapBidList{
		base:       b,
		PageParams: pageParams,
	}

	err = b.populateFilters(&request.Filters)
	if err != nil {
		return nil, err
	}

	return &request, nil
}
