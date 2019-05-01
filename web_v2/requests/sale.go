package requests

import (
	"time"

	history "gitlab.com/tokend/horizon/db2/history2"

	regources "gitlab.com/tokend/regources/generated"
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

type SalesBase struct {
	*base
	Filters struct {
		Owner        string           `json:"owner"`
		BaseAsset    string           `json:"base_asset"`
		MaxEndTime   *time.Time       `json:"max_end_time"`
		MaxStartTime *time.Time       `json:"max_start_time"`
		MinStartTime *time.Time       `json:"min_start_time"`
		MinEndTime   *time.Time       `json:"min_end_time"`
		State        uint64           `json:"state"`
		SaleType     uint64           `json:"sale_type"`
		MinHardCap   regources.Amount `json:"min_hard_cap"`
		MinSoftCap   regources.Amount `json:"min_soft_cap"`
		MaxHardCap   regources.Amount `json:"max_hard_cap"`
		MaxSoftCap   regources.Amount `json:"max_soft_cap"`
	}
}

func (s *SalesBase) ApplyFilters(q history.SalesQ) history.SalesQ {
	if s.ShouldFilter(FilterTypeSaleListOwner) {
		q = q.FilterByOwner(s.Filters.Owner)
	}

	if s.ShouldFilter(FilterTypeSaleListBaseAsset) {
		q = q.FilterByBaseAsset(s.Filters.BaseAsset)
	}

	if s.ShouldFilter(FilterTypeSaleListMaxEndTime) {
		q = q.FilterByMaxEndTime(*s.Filters.MaxEndTime)
	}

	if s.ShouldFilter(FilterTypeSaleListMaxStartTime) {
		q = q.FilterByMaxStartTime(*s.Filters.MaxStartTime)
	}

	if s.ShouldFilter(FilterTypeSaleListMinStartTime) {
		q = q.FilterByMinStartTime(*s.Filters.MinStartTime)
	}

	if s.ShouldFilter(FilterTypeSaleListMinEndTime) {
		q = q.FilterByMinEndTime(*s.Filters.MinEndTime)
	}

	if s.ShouldFilter(FilterTypeSaleListState) {
		q = q.FilterByState(s.Filters.State)
	}

	if s.ShouldFilter(FilterTypeSaleListSaleType) {
		q = q.FilterBySaleType(s.Filters.SaleType)
	}

	if s.ShouldFilter(FilterTypeSaleListMinHardCap) {
		q = q.FilterByMinHardCap(uint64(s.Filters.MinHardCap))
	}

	if s.ShouldFilter(FilterTypeSaleListMinSoftCap) {
		q = q.FilterByMinSoftCap(uint64(s.Filters.MinSoftCap))
	}

	if s.ShouldFilter(FilterTypeSaleListMaxHardCap) {
		q = q.FilterByMaxHardCap(uint64(s.Filters.MaxHardCap))
	}

	if s.ShouldFilter(FilterTypeSaleListMaxSoftCap) {
		q = q.FilterByMaxSoftCap(uint64(s.Filters.MaxSoftCap))
	}

	return q
}
