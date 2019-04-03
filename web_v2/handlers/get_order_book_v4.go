package handlers

import (
	"fmt"
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
	regources "gitlab.com/tokend/regources/generated"
)

// GetOrderBookV4 - processes request to get order book
func GetOrderBookV4(w http.ResponseWriter, r *http.Request) {
	coreRepo := ctx.CoreRepo(r)
	historyRepo := ctx.HistoryRepo(r)

	handler := getOrderBookV4Handler{
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
	Log         *logan.Entry
}

// GetOrderBookV4 returns order book with related resources
func (h *getOrderBookV4Handler) GetOrderBookV4(request *requests.GetOrderBookV4) (*regources.OrderBookResponse, error) {
	if request.OrderBookID != secondaryMarketOrderBookID {
		sale, err := h.SalesQ.GetByID(request.OrderBookID)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to get sale by ID")
		}
		if sale == nil {
			return nil, nil
		}
	}

	q := h.OrderBooksQ.
		FilterByOrderBookID(request.OrderBookID).
		FilterByBaseAssetCode(request.BaseAsset).
		FilterByQuoteAssetCode(request.QuoteAsset)

	buyEntriesQ := q.FilterByIsBuy(true)
	sellEntriesQ := q.FilterByIsBuy(false)

	if request.ShouldInclude(requests.IncludeTypeOrderBookV4BuyEntriesBaseAssets) {
		fmt.Println("inside if")
		buyEntriesQ = buyEntriesQ.WithBaseAsset()
	}

	if request.ShouldInclude(requests.IncludeTypeOrderBookV4BuyEntriesQuoteAssets) {
		buyEntriesQ = buyEntriesQ.WithQuoteAsset()
	}

	if request.ShouldInclude(requests.IncludeTypeOrderBookV4SellEntriesBaseAssets) {
		sellEntriesQ = sellEntriesQ.WithBaseAsset()
	}

	if request.ShouldInclude(requests.IncludeTypeOrderBookV4SellEntriesQuoteAssets) {
		sellEntriesQ = sellEntriesQ.WithQuoteAsset()
	}

	coreBuyEntries, err := buyEntriesQ.Select()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get buy order book entries")
	}

	coreSellEntries, err := sellEntriesQ.Select()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get sell order book entries")
	}

	response := &regources.OrderBookResponse{
		Data: regources.OrderBook{
			Relationships: regources.OrderBookRelationships{
				BuyEntries: &regources.RelationCollection{
					Data: make([]regources.Key, 0, len(coreBuyEntries)),
				},
				SellEntries: &regources.RelationCollection{
					Data: make([]regources.Key, 0, len(coreSellEntries)),
				},
			},
		},
	}

	for _, coreBuyEntry := range coreBuyEntries {
		response.Data.Relationships.BuyEntries.Data = append(
			response.Data.Relationships.BuyEntries.Data,
			resources.NewOrderBookEntryKey(coreBuyEntry.ID),
		)

		if request.ShouldInclude(requests.IncludeTypeOrderBookV4BuyEntries) {
			entry := resources.NewOrderBookEntry(coreBuyEntry)
			response.Included.Add(&entry)
		}

		if request.ShouldInclude(requests.IncludeTypeOrderBookV4BuyEntriesBaseAssets) {
			baseAsset := resources.NewAsset(*coreBuyEntry.BaseAsset)
			response.Included.Add(&baseAsset)
		}

		if request.ShouldInclude(requests.IncludeTypeOrderBookV4BuyEntriesQuoteAssets) {
			quoteAsset := resources.NewAsset(*coreBuyEntry.QuoteAsset)
			response.Included.Add(&quoteAsset)
		}
	}

	for _, coreSellEntry := range coreSellEntries {
		response.Data.Relationships.SellEntries.Data = append(
			response.Data.Relationships.BuyEntries.Data,
			resources.NewOrderBookEntryKey(coreSellEntry.ID),
		)

		if request.ShouldInclude(requests.IncludeTypeOrderBookV4SellEntries) {
			entry := resources.NewOrderBookEntry(coreSellEntry)
			response.Included.Add(&entry)
		}

		if request.ShouldInclude(requests.IncludeTypeOrderBookV4SellEntriesBaseAssets) {
			baseAsset := resources.NewAsset(*coreSellEntry.BaseAsset)
			response.Included.Add(&baseAsset)
		}

		if request.ShouldInclude(requests.IncludeTypeOrderBookV4SellEntriesQuoteAssets) {
			quoteAsset := resources.NewAsset(*coreSellEntry.QuoteAsset)
			response.Included.Add(&quoteAsset)
		}
	}

	return response, nil
}
