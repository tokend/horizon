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

// GetOrderBook - processes request to get order book
func GetOrderBook(w http.ResponseWriter, r *http.Request) {
	coreRepo := ctx.CoreRepo(r)
	historyRepo := ctx.HistoryRepo(r)

	handler := getOrderBookHandler{
		AssetsQ:     core2.NewAssetsQ(coreRepo),
		OrderBooksQ: core2.NewOrderBooksQ(coreRepo),
		SalesQ:      history2.NewSalesQ(historyRepo),
		Log:         ctx.Log(r),
	}

	request, err := requests.NewGetOrderBook(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	result, err := handler.GetOrderBook(request)
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

type getOrderBookHandler struct {
	OrderBooksQ core2.OrderBooksQ
	SalesQ      history2.SalesQ
	AssetsQ     core2.AssetsQ
	Log         *logan.Entry
}

// GetOrderBook returns order book with related resources
func (h *getOrderBookHandler) GetOrderBook(request *requests.GetOrderBook) (*regources.OrderBookResponse, error) {
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

	if request.ShouldInclude(requests.IncludeTypeOrderBookBaseAsset) {
		baseAsset := resources.NewAsset(*coreBaseAsset)
		response.Included.Add(&baseAsset)
	}

	if request.ShouldInclude(requests.IncludeTypeOrderBookQuoteAsset) {
		quoteAsset := resources.NewAsset(*coreQuoteAsset)
		response.Included.Add(&quoteAsset)
	}

	q := h.
		OrderBooksQ.
		FilterByOrderBookID(request.OrderBookID).
		FilterByBaseAssetCode(request.BaseAsset).
		FilterByQuoteAssetCode(request.QuoteAsset).
		Limit(request.MaxEntries)

	coreBuyEntries, err := q.FilterByIsBuy(true).OrderByPrice("desc").Select()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get buy entries")
	}

	coreSellEntries, err := q.FilterByIsBuy(false).OrderByPrice("asc").Select()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get sell entries")
	}

	for _, coreBuyEntry := range coreBuyEntries {
		entry := resources.NewOrderBookEntry(coreBuyEntry)
		response.Data.Relationships.BuyEntries.Data = append(response.Data.Relationships.BuyEntries.Data, entry.Key)

		if request.ShouldInclude(requests.IncludeTypeOrderBookBuyEntries) {
			response.Included.Add(&entry)
		}
	}

	for _, coreSellEntry := range coreSellEntries {
		entry := resources.NewOrderBookEntry(coreSellEntry)
		response.Data.Relationships.SellEntries.Data = append(response.Data.Relationships.SellEntries.Data, entry.Key)

		if request.ShouldInclude(requests.IncludeTypeOrderBookSellEntries) {
			response.Included.Add(&entry)
		}
	}

	return response, nil
}
