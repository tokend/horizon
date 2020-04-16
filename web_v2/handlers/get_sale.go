package handlers

import (
	"net/http"

	"gitlab.com/tokend/horizon/db2/history2"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/tokend/horizon/web_v2/ctx"
	"gitlab.com/tokend/horizon/web_v2/requests"
	regources "gitlab.com/tokend/regources/generated"
)

// GetSale - processes request to get sale and it's details by sale ID
func GetSale(w http.ResponseWriter, r *http.Request) {
	historyRepo := ctx.HistoryRepo(r)

	converter := newSaleCapConverterForHandler(w, r)
	if converter == nil {
		return
	}

	handler := getSaleHandler{
		getSaleBase{
			SalesQ:           history2.NewSalesQ(historyRepo),
			AssetsQ:          history2.NewAssetQ(historyRepo),
			saleCapConverter: converter,
			Log:              ctx.Log(r),
		},
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
	getSaleBase
}

// GetSale returns sale with related resources
func (h *getSaleHandler) GetSale(request *requests.GetSale) (*regources.SaleResponse, error) {
	q := h.SalesQ.FilterByID(request.ID)

	if request.ShouldInclude(requests.IncludeTypeSaleListBaseAssets) {
		q = q.WithAsset()
	}

	return h.getAndPopulateResponse(q, request)
}
