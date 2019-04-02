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

// GetOrderBook - processes request to get order book entries
func GetOrderBook(w http.ResponseWriter, r *http.Request) {
	coreRepo := ctx.CoreRepo(r)
	historyRepo := ctx.HistoryRepo(r)

	handler := getOrderBookHandler{
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
		ctx.Log(r).WithError(err).Error("failed to get order book entries", logan.F{
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
	Log         *logan.Entry
}

const secondaryMarketOrderBookID = 0

// GetOrderBook returns offer with related resources
func (h *getOrderBookHandler) GetOrderBook(request *requests.GetOrderBook) (*regources.OrderBookEntrysResponse, error) {
	if request.ID != secondaryMarketOrderBookID {
		sale, err := h.SalesQ.GetByID(request.ID)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to get sale by ID")
		}
		if sale == nil {
			return nil, nil
		}
	}

	q := h.OrderBooksQ.Page(*request.PageParams).FilterByOrderBookID(request.ID)

	if request.ShouldInclude(requests.IncludeTypeOrderBookBaseAssets) {
		q = q.WithBaseAsset()
	}

	if request.ShouldInclude(requests.IncludeTypeOrderBookQuoteAssets) {
		q = q.WithQuoteAsset()
	}

	if request.ShouldFilter(requests.FilterTypeOrderBookBaseAsset) {
		q = q.FilterByBaseAssetCode(request.Filters.BaseAsset)
	}

	if request.ShouldFilter(requests.FilterTypeOrderBookQuoteAsset) {
		q = q.FilterByQuoteAssetCode(request.Filters.QuoteAsset)
	}

	if request.ShouldFilter(requests.FilterTypeOrderBookIsBuy) {
		q = q.FilterByIsBuy(request.Filters.IsBuy)
	}

	coreOrderBookEntries, err := q.Select()

	if err != nil {
		return nil, errors.Wrap(err, "Failed to get offer list")
	}

	response := &regources.OrderBookEntrysResponse{
		Data:  make([]regources.OrderBookEntry, 0, len(coreOrderBookEntries)),
		Links: request.GetOffsetLinks(*request.PageParams),
	}

	for _, coreOrderBookEntry := range coreOrderBookEntries {
		response.Data = append(response.Data, resources.NewOrderBookEntry(coreOrderBookEntry))

		if request.ShouldInclude(requests.IncludeTypeOrderBookBaseAssets) {
			baseAsset := resources.NewAsset(*coreOrderBookEntry.BaseAsset)
			response.Included.Add(&baseAsset)
		}

		if request.ShouldInclude(requests.IncludeTypeOrderBookQuoteAssets) {
			quoteAsset := resources.NewAsset(*coreOrderBookEntry.QuoteAsset)
			response.Included.Add(&quoteAsset)
		}
	}

	return response, nil
}
