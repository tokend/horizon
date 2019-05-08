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
		salesBaseHandler: salesBaseHandler{
			SalesQ:                history2.NewSalesQ(historyRepo),
			AccountSpecificRulesQ: history2.NewAccountSpecificRulesQ(historyRepo),
			AssetsQ:               core2.NewAssetsQ(coreRepo),
			saleCapConverter:      converter,
			Log:                   ctx.Log(r),
		},
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
	salesBaseHandler
}

// GetSaleList returns the list of assets with related resources
func (h *getSaleListHandler) GetSaleList(request *requests.GetSaleList) (*regources.SalesResponse, error) {
	q := applySaleFilters(request.SalesBase, h.SalesQ)

	historySales, err := q.Page(*request.PageParams).Select()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get asset list")
	}

	response := &regources.SalesResponse{
		Data:  make([]regources.Sale, 0, len(historySales)),
		Links: request.GetOffsetLinks(*request.PageParams),
	}

	err = h.populateResponse(historySales, request.SalesBase, response)
	if err != nil {
		return nil, errors.Wrap(err, "failed to populate response")
	}

	return response, nil
}
