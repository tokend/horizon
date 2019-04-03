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
func (h *getOrderBookV4Handler) GetOrderBookV4(request *requests.GetOrderBookV4) (*resources.OrderBookResponse, error) {
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

	if request.ShouldInclude(requests.IncludeTypeOrderBookV4BaseAssets) {
		q = q.WithBaseAsset()
	}

	if request.ShouldInclude(requests.IncludeTypeOrderBookV4QuoteAssets) {
		q = q.WithQuoteAsset()
	}

	isBuyQ := q.FilterByIsBuy(true)
	isSellQ := q.FilterByIsBuy(false)

	coreBuyEntries, err := isBuyQ.Select()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get buy order book entries")
	}

	coreSellEntries, err := isSellQ.Select()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get sell order book entries")
	}

	response := &resources.OrderBookResponse{
		Data: resources.OrderBook{
			Relationships: resources.OrderBookRelations{
				BuyEntries: &regources.RelationCollection{
					Data: make([]regources.Key, 0, len(coreBuyEntries)),
				},
				SellEntries: &regources.RelationCollection{
					Data: make([]regources.Key, 0, len(coreSellEntries)),
				},
			},
		},
	}

	for _, coreEntry := range coreBuyEntries {
		response.Data.Relationships.BuyEntries.Data = append(
			response.Data.Relationships.BuyEntries.Data,
			resources.NewOrderBookEntryKey(coreEntry.ID),
		)

		if request.ShouldInclude(requests.IncludeTypeOrderBookV4OrderBookEntries) {
			entry := resources.NewOrderBookEntry(coreEntry)
			response.Included.Add(&entry)
		}

		if request.ShouldInclude(requests.IncludeTypeOrderBookV4BaseAssets) {
			baseAsset := resources.NewAsset(*coreEntry.BaseAsset)
			response.Included.Add(&baseAsset)
		}

		if request.ShouldInclude(requests.IncludeTypeOrderBookV4QuoteAssets) {
			quoteAsset := resources.NewAsset(*coreEntry.QuoteAsset)
			response.Included.Add(&quoteAsset)
		}
	}

	// TODO: reduce duplication
	for _, coreEntry := range coreSellEntries {
		response.Data.Relationships.SellEntries.Data = append(
			response.Data.Relationships.BuyEntries.Data,
			resources.NewOrderBookEntryKey(coreEntry.ID),
		)

		if request.ShouldInclude(requests.IncludeTypeOrderBookV4OrderBookEntries) {
			entry := resources.NewOrderBookEntry(coreEntry)
			response.Included.Add(&entry)
		}

		if request.ShouldInclude(requests.IncludeTypeOrderBookV4BaseAssets) {
			baseAsset := resources.NewAsset(*coreEntry.BaseAsset)
			response.Included.Add(&baseAsset)
		}

		if request.ShouldInclude(requests.IncludeTypeOrderBookV4QuoteAssets) {
			quoteAsset := resources.NewAsset(*coreEntry.QuoteAsset)
			response.Included.Add(&quoteAsset)
		}
	}

	return response, nil
}
