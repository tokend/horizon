package requests

import (
	"gitlab.com/tokend/horizon/db2"
	"gitlab.com/tokend/regources/rgenerated"
	"net/http"
	"time"
)

const (
	// IncludeTypeSaleListBaseAssets - defines if the base assets should be included in the response
	IncludeTypeSaleListBaseAssets = "base_asset"
	// IncludeTypeSaleListQuoteAssets - defines if the quote assets should be included in the response
	IncludeTypeSaleListQuoteAssets = "quote_assets"
	// IncludeTypeSaleListDefaultQuoteAsset - defines if the default quote asset should be included in the response
	IncludeTypeSaleListDefaultQuoteAsset = "default_quote_asset"

	// FilterTypeSaleListOwner - defines if we need to filter resopnse by owner
	FilterTypeSaleListOwner = "owner"
	// FilterTypeSaleListBaseAsset - defines if we need to filter resopnse by base_asset
	FilterTypeSaleListBaseAsset = "base_asset"
	// FilterTypeSaleListMaxEndTime - defines if we need to filter response by max_end_time
	FilterTypeSaleListMaxEndTime = "max_end_time"
	// FilterTypeSaleListMaxStartTime - defines if we need to filter response by max_start_time
	FilterTypeSaleListMaxStartTime = "max_start_time"
	// FilterTypeSaleListMinStartTime - defines if we need to filter response by min_start_time
	FilterTypeSaleListMinStartTime = "min_start_time"
	// FilterTypeSaleListMinEndTime - defines if we need to filter response by min_end_time
	FilterTypeSaleListMinEndTime = "min_end_time"
	// FilterTypeSaleListState - defines if we need to filter response by state
	FilterTypeSaleListState = "state"
	// FilterTypeSaleListSaleType - defines if we need to filter response by sale_type
	FilterTypeSaleListSaleType = "sale_type"
	// FilterTypeSaleListMinHardCap - defines if we need to filter response by min_hard_cap
	FilterTypeSaleListMinHardCap = "min_hard_cap"
	// FilterTypeSaleListMinSoftCap - defines if we need to filter response by min_soft_cap
	FilterTypeSaleListMinSoftCap = "min_soft_cap"
	// FilterTypeSaleListMaxHardCap - defines if we need to filter response by max_hard_cap
	FilterTypeSaleListMaxHardCap = "max_hard_cap"
	// FilterTypeSaleListMaxSoftCap - defines if we need to filter response by max_soft_cap
	FilterTypeSaleListMaxSoftCap = "max_soft_cap"
)

var includeTypeSaleListAll = map[string]struct{}{
	IncludeTypeSaleListBaseAssets:        {},
	IncludeTypeSaleListQuoteAssets:       {},
	IncludeTypeSaleListDefaultQuoteAsset: {},
}

var filterTypeSaleListAll = map[string]struct{}{
	FilterTypeSaleListOwner:        {},
	FilterTypeSaleListBaseAsset:    {},
	FilterTypeSaleListMaxEndTime:   {},
	FilterTypeSaleListMaxStartTime: {},
	FilterTypeSaleListMinStartTime: {},
	FilterTypeSaleListMinEndTime:   {},
	FilterTypeSaleListState:        {},
	FilterTypeSaleListSaleType:     {},
	FilterTypeSaleListMinHardCap:   {},
	FilterTypeSaleListMinSoftCap:   {},
	FilterTypeSaleListMaxHardCap:   {},
	FilterTypeSaleListMaxSoftCap:   {},
}

// GetSaleList - represents params to be specified by user for getSaleList handler
type GetSaleList struct {
	*base
	Filters struct {
		Owner        string            `json:"owner"`
		BaseAsset    string            `json:"base_asset"`
		MaxEndTime   *time.Time        `json:"max_end_time"`
		MaxStartTime *time.Time        `json:"max_start_time"`
		MinStartTime *time.Time        `json:"min_start_time"`
		MinEndTime   *time.Time        `json:"min_end_time"`
		State        uint64            `json:"state"`
		SaleType     uint64            `json:"sale_type"`
		MinHardCap   rgenerated.Amount `json:"min_hard_cap"`
		MinSoftCap   rgenerated.Amount `json:"min_soft_cap"`
		MaxHardCap   rgenerated.Amount `json:"max_hard_cap"`
		MaxSoftCap   rgenerated.Amount `json:"max_soft_cap"`
	}
	PageParams *db2.OffsetPageParams
}

// NewGetSaleList returns new instance of GetSaleList request
func NewGetSaleList(r *http.Request) (*GetSaleList, error) {
	b, err := newBase(r, baseOpts{
		supportedIncludes: includeTypeSaleListAll,
		supportedFilters:  filterTypeSaleListAll,
	})
	if err != nil {
		return nil, err
	}

	pageParams, err := b.getOffsetBasedPageParams()
	if err != nil {
		return nil, err
	}

	request := GetSaleList{
		base:       b,
		PageParams: pageParams,
	}

	err = b.populateFilters(&request.Filters)
	if err != nil {
		return nil, err
	}

	return &request, nil
}
