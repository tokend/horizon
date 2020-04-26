package requests

import (
	"gitlab.com/tokend/horizon/bridge"
	"net/http"
)

const (
	IncludeTypeAskListBaseBalances = "base_balance"
	IncludeTypeAskListOwners       = "owner"
	IncludeTypeAskListBaseAssets   = "base_asset"
	IncludeTypeAskListQuoteAssets  = "quote_assets"

	FilterTypeAskListOwner       = "owner"
	FilterTypeAskListBaseAsset   = "base_asset"
	FilterTypeAskListQuoteAssets = "quote_assets"
)

var includeTypeAskListAll = map[string]struct{}{
	IncludeTypeAskListBaseBalances: {},
	IncludeTypeAskListOwners:       {},
	IncludeTypeAskListBaseAssets:   {},
	IncludeTypeAskListQuoteAssets:  {},
}

var filterTypeAskListAll = map[string]struct{}{
	FilterTypeAskListOwner:       {},
	FilterTypeAskListBaseAsset:   {},
	FilterTypeAskListQuoteAssets: {},
}

// GetAtomicSwapAskList - represents params to be specified by user for Get AtomicSwapAskList handler
type GetAtomicSwapAskList struct {
	*base
	Filters struct {
		Owner       string   `fig:"owner"`
		BaseAsset   string   `fig:"base_asset"`
		QuoteAssets []string `fig:"quote_assets"`
	}
	PageParams *bridge.OffsetPageParams
}

// NewGetAtomicSwapAskList returns new instance of GetAtomicSwapAskList request
func NewGetAtomicSwapAskList(r *http.Request) (*GetAtomicSwapAskList, error) {
	b, err := newBase(r, baseOpts{
		supportedIncludes: includeTypeAskListAll,
		supportedFilters:  filterTypeAskListAll,
	})
	if err != nil {
		return nil, err
	}

	// bid relations has not asset relation, we use balance relation
	if _, ok := b.include[IncludeTypeAskBaseAsset]; ok {
		b.include[IncludeTypeAskBaseBalance] = struct{}{}
	}

	pageParams, err := b.getOffsetBasedPageParams()
	if err != nil {
		return nil, err
	}

	request := GetAtomicSwapAskList{
		base:       b,
		PageParams: pageParams,
	}

	err = b.populateFilters(&request.Filters)
	if err != nil {
		return nil, err
	}

	return &request, nil
}
