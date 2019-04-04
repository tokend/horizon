package handlers

import (
	"net/http"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2/core2"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/horizon/web_v2/ctx"
	"gitlab.com/tokend/horizon/web_v2/requests"
	"gitlab.com/tokend/horizon/web_v2/resources"
	"gitlab.com/tokend/regources/generated"
)

// GetOrderBookV4 - processes request to get order book
func GetOrderBookV4(w http.ResponseWriter, r *http.Request) {
	coreRepo := ctx.CoreRepo(r)
	historyRepo := ctx.HistoryRepo(r)

	handler := getOrderBookV4Handler{
		AssetsQ:     core2.NewAssetsQ(coreRepo),
		OrderBooksQ: core2.NewOrderBooksQ(coreRepo),
		SalesQ:      history2.NewSalesQ(historyRepo),
		Log:         ctx.Log(r),
	}

	request, err := requests.NewGetOrderBookV4(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	result, err := handler.GetOrderBookV4(request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get order book", logan.F{
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

type getOrderBookV4Handler struct {
	OrderBooksQ core2.OrderBooksQ
	SalesQ      history2.SalesQ
	AssetsQ     core2.AssetsQ
	Log         *logan.Entry
}

// GetOrderBookV4 returns order book with related resources
func (h *getOrderBookV4Handler) GetOrderBookV4(request *requests.GetOrderBookV4) (*regources.OrderBookResponse, error) {
	if request.OrderBookID != secondaryMarketOrderBookID {
		coreSale, err := h.SalesQ.GetByID(request.OrderBookID)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to get sale by ID")
		}
		if coreSale == nil {
			return nil, nil
		}
	}

	coreBaseAsset, err := h.AssetsQ.GetByCode(request.BaseAsset)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get base asset by code")
	}
	if coreBaseAsset == nil {
		return nil, nil
	}

	coreQuoteAsset, err := h.AssetsQ.GetByCode(request.QuoteAsset)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get quote asset by code")
	}
	if coreQuoteAsset == nil {
		return nil, nil
	}

	response := &regources.OrderBookResponse{
		Data: resources.NewOrderBook(request.BaseAsset, request.QuoteAsset, request.OrderBookID),
	}
	response.Data.Relationships.BuyEntries = &regources.RelationCollection{
		Data: []regources.Key{},
	}
	response.Data.Relationships.SellEntries = &regources.RelationCollection{
		Data: []regources.Key{},
	}

	if request.ShouldInclude(requests.IncludeTypeOrderBookV4BaseAsset) {
		baseAsset := resources.NewAsset(*coreBaseAsset)
		response.Included.Add(&baseAsset)
	}

	if request.ShouldInclude(requests.IncludeTypeOrderBookV4QuoteAsset) {
		quoteAsset := resources.NewAsset(*coreQuoteAsset)
		response.Included.Add(&quoteAsset)
	}

	coreOrderBookEntries, err := h.OrderBooksQ.
		FilterByOrderBookID(request.OrderBookID).
		FilterByBaseAssetCode(request.BaseAsset).
		FilterByQuoteAssetCode(request.QuoteAsset).
		Select()

	if err != nil {
		return nil, errors.Wrap(err, "Failed to get order book entries")
	}

	for _, coreOrderBookEntry := range coreOrderBookEntries {
		orderBookEntry := resources.NewOrderBookEntry(coreOrderBookEntry)
		if orderBookEntry.Attributes.IsBuy {
			response.Data.Relationships.BuyEntries.Data = append(
				response.Data.Relationships.BuyEntries.Data,
				orderBookEntry.Key,
			)
			if request.ShouldInclude(requests.IncludeTypeOrderBookV4BuyEntries) {
				response.Included.Add(&orderBookEntry)
			}
		} else {
			response.Data.Relationships.SellEntries.Data = append(
				response.Data.Relationships.SellEntries.Data,
				orderBookEntry.Key,
			)
			if request.ShouldInclude(requests.IncludeTypeOrderBookV4SellEntries) {
				response.Included.Add(&orderBookEntry)
			}
		}
	}

	return response, nil
}
