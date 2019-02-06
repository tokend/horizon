package handlers

import (
	"net/http"

	"gitlab.com/tokend/horizon/db2/core2"
	"gitlab.com/tokend/horizon/db2/history2"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/web_v2/ctx"
	"gitlab.com/tokend/horizon/web_v2/requests"
	"gitlab.com/tokend/horizon/web_v2/resources"
	"gitlab.com/tokend/regources/v2"
)

// GetSale - processes request to get sale and it's details by sale ID
func GetSale(w http.ResponseWriter, r *http.Request) {
	historyRepo := ctx.HistoryRepo(r)
	coreRepo := ctx.CoreRepo(r)

	converter := newSaleCapConverterForHandler(w, r)
	if converter == nil {
		return
	}

	handler := getSaleHandler{
		SalesQ:           history2.NewSalesQ(historyRepo),
		AssetsQ:          core2.NewAssetsQ(coreRepo),
		saleCapConverter: converter,
		Log:              ctx.Log(r),
	}

	request, err := requests.NewGetSale(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	result, err := handler.GetSale(request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get sale", logan.F{
			"request": request,
		})
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if result == nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}

	ape.Render(w, result)
}

type getSaleHandler struct {
	SalesQ           history2.SalesQ
	AssetsQ          core2.AssetsQ
	saleCapConverter *saleCapConverter
	Log              *logan.Entry
}

// GetSale returns sale with related resources
func (h *getSaleHandler) GetSale(request *requests.GetSale) (*regources.SaleResponse, error) {
	historySale, err := h.SalesQ.GetByID(request.ID)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get sale by ID")
	}
	if historySale == nil {
		return nil, nil
	}

	err = h.saleCapConverter.PopulateSaleCap(historySale)
	if err != nil {
		return nil, errors.Wrap(err, "failed to populate sale cap")
	}

	response := &regources.SaleResponse{
		Data: resources.NewSale(*historySale),
	}

	quoteAssets := &regources.RelationCollection{
		Data: make([]regources.Key, 0, len(historySale.QuoteAssets.QuoteAssets)),
	}

	for _, historyQuoteAsset := range historySale.QuoteAssets.QuoteAssets {
		quoteAsset := resources.NewSaleQuoteAsset(historyQuoteAsset)
		quoteAssets.Data = append(quoteAssets.Data, quoteAsset.Key)

		if request.ShouldInclude(requests.IncludeTypeSaleQuoteAssets) {
			response.Included.Add(&quoteAsset)
		}
	}
	response.Data.Relationships.QuoteAssets = quoteAssets

	defaultQuoteAsset := resources.NewSaleDefaultQuoteAsset(*historySale)
	response.Data.Relationships.DefaultQuoteAsset = defaultQuoteAsset.AsRelation()

	if request.ShouldInclude(requests.IncludeTypeSaleDefaultQuoteAsset) {
		response.Included.Add(&defaultQuoteAsset)
	}

	if request.ShouldInclude(requests.IncludeTypeSaleBaseAsset) {
		// FIXME: ingest assets to history and join
		coreAsset, err := h.AssetsQ.GetByCode(historySale.BaseAsset)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to get base asset by code")
		}

		asset := resources.NewAsset(*coreAsset)
		response.Included.Add(&asset)
	}

	return response, nil
}
