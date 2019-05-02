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
	regources "gitlab.com/tokend/regources/generated"
)

// GetSaleListForAccount - processes request to get the list of sales
func GetSaleListForAccount(w http.ResponseWriter, r *http.Request) {
	historyRepo := ctx.HistoryRepo(r)
	coreRepo := ctx.CoreRepo(r)

	converter := newSaleCapConverterForHandler(w, r)
	if converter == nil {
		return
	}

	handler := getSaleListForAccountHandler{
		AccountSpecificRulesQ: history2.NewAccountSpecificRulesQ(historyRepo),
		salesBaseHandler: salesBaseHandler{
			SalesQ:           history2.NewSalesQ(historyRepo),
			AssetsQ:          core2.NewAssetsQ(coreRepo),
			saleCapConverter: converter,
			Log:              ctx.Log(r),
		},
	}

	request, err := requests.NewGetSaleListForAccount(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	if !isAllowed(r, w, request.Address) {
		return
	}

	result, err := handler.GetSaleListForAccount(request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get sale list", logan.F{
			"request": request,
		})
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, result)
}

type getSaleListForAccountHandler struct {
	salesBaseHandler
	AccountSpecificRulesQ history2.AccountSpecificRulesQ
}

// GetSaleListForAccount returns the list of assets with related resources
func (h *getSaleListForAccountHandler) GetSaleListForAccount(request *requests.GetSaleListForAccount) (*regources.SalesResponse, error) {
	q := request.ApplyFilters(h.SalesQ).Whitelisted(request.Address)

	historySales, err := q.CursorPage(*request.PageParams).Select()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get sales list")
	}
	if historySales == nil {
		return nil, nil
	}

	response := &regources.SalesResponse{
		Data: make([]regources.Sale, 0, len(historySales)),
	}

	err = h.populateResponse(historySales, request.SalesBase, response)
	if err != nil {
		return nil, errors.Wrap(err, "failed to populate response")
	}

	h.populateLinks(response, request)

	return response, nil
}

func (h *getSaleListForAccountHandler) populateLinks(
	response *regources.SalesResponse, request *requests.GetSaleListForAccount,
) {
	if len(response.Data) > 0 {
		response.Links = request.GetCursorLinks(*request.PageParams, response.Data[len(response.Data)-1].ID)
	} else {
		response.Links = request.GetCursorLinks(*request.PageParams, "")
	}
}
