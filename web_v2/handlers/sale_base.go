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
	response *regources.SalesResponse) error {

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
			quoteAsset := resources.NewSaleQuoteAsset(historyQuoteAsset)
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
