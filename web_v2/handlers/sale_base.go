package handlers

import (
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2/core2"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/horizon/web_v2/requests"
	"gitlab.com/tokend/horizon/web_v2/resources"
	regources "gitlab.com/tokend/regources/generated"
)

type salesBaseHandler struct {
	SalesQ  history2.SalesQ
	AssetsQ core2.AssetsQ

	saleCapConverter *saleCapConverter
	Log              *logan.Entry
}

func (h *salesBaseHandler) populateResponse(historySales []history2.Sale,
	request requests.SalesBase,
	response *regources.SaleListResponse) error {

	for _, historySale := range historySales {
		sale := resources.NewSale(historySale)
		err := h.saleCapConverter.PopulateSaleCap(&historySale)
		if err != nil {
			return errors.Wrap(err, "failed to populate sale cap")
		}
		sale.Relationships.QuoteAssets = &regources.RelationCollection{
			Data: make([]regources.Key, 0, len(historySale.QuoteAssets.QuoteAssets)),
		}

		defaultQuoteAsset := resources.NewSaleDefaultQuoteAsset(historySale)
		sale.Relationships.DefaultQuoteAsset = defaultQuoteAsset.AsRelation()

		if request.ShouldInclude(requests.IncludeTypeSaleDefaultQuoteAsset) {
			response.Included.Add(&defaultQuoteAsset)
		}

		for _, historyQuoteAsset := range historySale.QuoteAssets.QuoteAssets {
			quoteAsset := resources.NewSaleQuoteAsset(historyQuoteAsset, historySale.ID)
			sale.Relationships.QuoteAssets.Data = append(sale.Relationships.QuoteAssets.Data, quoteAsset.Key)

			if request.ShouldInclude(requests.IncludeTypeSaleQuoteAssets) {
				response.Included.Add(&quoteAsset)
			}
		}

		if request.ShouldInclude(requests.IncludeTypeSaleListBaseAssets) {
			// FIXME: ingest assets to history and join
			coreAsset, err := h.AssetsQ.GetByCode(historySale.BaseAsset)
			if err != nil {
				return errors.Wrap(err, "Failed to get base asset by code")
			}

			asset := resources.NewAsset(*coreAsset)
			response.Included.Add(&asset)
		}

		response.Data = append(response.Data, sale)
	}
	return nil
}

func applySaleFilters(s requests.SalesBase, q history2.SalesQ) history2.SalesQ {
	if s.ShouldFilter(requests.FilterTypeSaleListOwner) {
		q = q.FilterByOwner(s.Filters.Owner)
	}

	if s.ShouldFilter(requests.FilterTypeSaleListBaseAsset) {
		q = q.FilterByBaseAsset(s.Filters.BaseAsset)
	}

	if s.ShouldFilter(requests.FilterTypeSaleListMaxEndTime) {
		q = q.FilterByMaxEndTime(*s.Filters.MaxEndTime)
	}

	if s.ShouldFilter(requests.FilterTypeSaleListMaxStartTime) {
		q = q.FilterByMaxStartTime(*s.Filters.MaxStartTime)
	}

	if s.ShouldFilter(requests.FilterTypeSaleListMinStartTime) {
		q = q.FilterByMinStartTime(*s.Filters.MinStartTime)
	}

	if s.ShouldFilter(requests.FilterTypeSaleListMinEndTime) {
		q = q.FilterByMinEndTime(*s.Filters.MinEndTime)
	}

	if s.ShouldFilter(requests.FilterTypeSaleListState) {
		q = q.FilterByState(s.Filters.State)
	}

	if s.ShouldFilter(requests.FilterTypeSaleListSaleType) {
		q = q.FilterBySaleType(s.Filters.SaleType)
	}

	if s.ShouldFilter(requests.FilterTypeSaleListMinHardCap) {
		q = q.FilterByMinHardCap(uint64(s.Filters.MinHardCap))
	}

	if s.ShouldFilter(requests.FilterTypeSaleListMinSoftCap) {
		q = q.FilterByMinSoftCap(uint64(s.Filters.MinSoftCap))
	}

	if s.ShouldFilter(requests.FilterTypeSaleListMaxHardCap) {
		q = q.FilterByMaxHardCap(uint64(s.Filters.MaxHardCap))
	}

	if s.ShouldFilter(requests.FilterTypeSaleListMaxSoftCap) {
		q = q.FilterByMaxSoftCap(uint64(s.Filters.MaxSoftCap))
	}

	return q
}
