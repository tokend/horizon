package requests

import (
	"net/http"

	"gitlab.com/tokend/horizon/db2"
)

const (
	// IncludeTypeSaleParticipationBaseAsset - defines if the base asset should be included in the response
	IncludeTypeSaleParticipationBaseAsset = "base_asset"
	// IncludeTypeSaleParticipationQuoteAsset - defines if the quote asset should be included in the response
	IncludeTypeSaleParticipationQuoteAsset = "quote_asset"

	// FilterTypeSaleListParticipant - defines if we need to filter response by participant
	FilterTypeSaleParticipationParticipant = "participant"
	// FilterTypeSaleListQuoteAsset - defines if we need to filter response by quote_asset
	FilterTypeSaleParticipationQuoteAsset = "quote_asset"
)

var includeTypeSaleParticipationAll = map[string]struct{}{
	IncludeTypeSaleParticipationQuoteAsset: {},
	IncludeTypeSaleParticipationBaseAsset:  {},
}

var filterTypeSaleParticipationAll = map[string]struct{}{
	FilterTypeSaleParticipationParticipant: {},
	FilterTypeSaleParticipationQuoteAsset:  {},
}

type GetSaleParticipation struct {
	*base
	SaleID  uint64
	Filters struct {
		QuoteAsset  string `json:"quote_asset"`
		Participant string `json:"participant"`
	}
	PageParams *db2.CursorPageParams
}

func NewGetSaleParticipation(r *http.Request) (*GetSaleParticipation, error) {
	b, err := newBase(r, baseOpts{
		supportedFilters:  filterTypeSaleParticipationAll,
		supportedIncludes: includeTypeSaleParticipationAll,
	})
	if err != nil {
		return nil, err
	}

	id, err := b.getUint64ID()
	if err != nil {
		return nil, err
	}

	pageParams, err := b.getCursorBasedPageParams()
	if err != nil {
		return nil, err
	}

	request := &GetSaleParticipation{
		base:       b,
		SaleID:     id,
		PageParams: pageParams,
	}

	err = b.populateFilters(&request.Filters)
	if err != nil {
		return nil, err
	}

	return request, nil
}
