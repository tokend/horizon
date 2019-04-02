package handlers

import (
	"net/http"

	"gitlab.com/tokend/horizon/db2/history2"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2/core2"
	"gitlab.com/tokend/horizon/web_v2/ctx"
	"gitlab.com/tokend/horizon/web_v2/requests"
	"gitlab.com/tokend/horizon/web_v2/resources"
	regources "gitlab.com/tokend/regources/generated"
)

// GetSaleList - processes request to get the list of sales
func GetSaleList(w http.ResponseWriter, r *http.Request) {
	historyRepo := ctx.HistoryRepo(r)
	coreRepo := ctx.CoreRepo(r)

	converter := newSaleCapConverterForHandler(w, r)
	if converter == nil {
		return
	}

	handler := getSaleListHandler{
		SalesQ:           history2.NewSalesQ(historyRepo),
		AssetsQ:          core2.NewAssetsQ(coreRepo),
		saleCapConverter: converter,
		Log:              ctx.Log(r),
	}

	request, err := requests.NewGetSaleList(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	result, err := handler.GetSaleList(request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get sale list", logan.F{
			"request": request,
		})
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, result)
}

type getSaleListHandler struct {
	SalesQ           history2.SalesQ
	AssetsQ          core2.AssetsQ
	saleCapConverter *saleCapConverter
	Log              *logan.Entry
}

// GetSaleList returns the list of assets with related resources
func (h *getSaleListHandler) GetSaleList(request *requests.GetSaleList) (*regources.SalesResponse, error) {
	q := h.SalesQ

	if request.ShouldFilter(requests.FilterTypeSaleListOwner) {
		q = q.FilterByOwner(request.Filters.Owner)
	}

	if request.ShouldFilter(requests.FilterTypeSaleListBaseAsset) {
		q = q.FilterByBaseAsset(request.Filters.BaseAsset)
	}

	if request.ShouldFilter(requests.FilterTypeSaleListMaxEndTime) {
		q = q.FilterByMaxEndTime(*request.Filters.MaxEndTime)
	}

	if request.ShouldFilter(requests.FilterTypeSaleListMaxStartTime) {
		q = q.FilterByMaxStartTime(*request.Filters.MaxStartTime)
	}

	if request.ShouldFilter(requests.FilterTypeSaleListMinStartTime) {
		q = q.FilterByMinStartTime(*request.Filters.MinStartTime)
	}

	if request.ShouldFilter(requests.FilterTypeSaleListMinEndTime) {
		q = q.FilterByMinEndTime(*request.Filters.MinEndTime)
	}

	if request.ShouldFilter(requests.FilterTypeSaleListState) {
		q = q.FilterByState(request.Filters.State)
	}

	if request.ShouldFilter(requests.FilterTypeSaleListSaleType) {
		q = q.FilterBySaleType(request.Filters.SaleType)
	}

	if request.ShouldFilter(requests.FilterTypeSaleListMinHardCap) {
		q = q.FilterByMinHardCap(uint64(request.Filters.MinHardCap))
	}

	if request.ShouldFilter(requests.FilterTypeSaleListMinSoftCap) {
		q = q.FilterByMinSoftCap(uint64(request.Filters.MinSoftCap))
	}

	if request.ShouldFilter(requests.FilterTypeSaleListMaxHardCap) {
		q = q.FilterByMaxHardCap(uint64(request.Filters.MaxHardCap))
	}

	if request.ShouldFilter(requests.FilterTypeSaleListMaxSoftCap) {
		q = q.FilterByMaxSoftCap(uint64(request.Filters.MaxSoftCap))
	}

	historySales, err := q.Page(*request.PageParams).Select()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get asset list")
	}

	response := &regources.SalesResponse{
		Data:  make([]regources.Sale, 0, len(historySales)),
		Links: request.GetOffsetLinks(*request.PageParams),
	}

	for _, historySale := range historySales {
		sale := resources.NewSale(historySale)
		err = h.saleCapConverter.PopulateSaleCap(&historySale)
		if err != nil {
			return nil, errors.Wrap(err, "failed to populate sale cap")
		}
		sale.Relationships.QuoteAssets = &regources.RelationCollection{
			Data: make([]regources.Key, 0, len(historySale.QuoteAssets.QuoteAssets)),
		}

		for _, historyQuoteAsset := range historySale.QuoteAssets.QuoteAssets {
			quoteAsset := resources.NewSaleQuoteAsset(historyQuoteAsset)
			sale.Relationships.QuoteAssets.Data = append(sale.Relationships.QuoteAssets.Data, quoteAsset.Key)

			if request.ShouldInclude(requests.IncludeTypeSaleQuoteAssets) {
				response.Included.Add(&quoteAsset)
			}
		}

		defaultQuoteAsset := resources.NewSaleDefaultQuoteAsset(historySale)
		sale.Relationships.DefaultQuoteAsset = defaultQuoteAsset.AsRelation()

		if request.ShouldInclude(requests.IncludeTypeSaleDefaultQuoteAsset) {
			response.Included.Add(&defaultQuoteAsset)
		}

		if request.ShouldInclude(requests.IncludeTypeSaleListBaseAssets) {
			// FIXME: ingest assets to history and join
			coreAsset, err := h.AssetsQ.GetByCode(historySale.BaseAsset)
			if err != nil {
				return nil, errors.Wrap(err, "Failed to get base asset by code")
			}

			asset := resources.NewAsset(*coreAsset)
			response.Included.Add(&asset)
		}

		response.Data = append(response.Data, sale)
	}

	return response, nil
}
