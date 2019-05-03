package handlers

import (
	"net/http"

	"github.com/google/jsonapi"
	"gitlab.com/tokend/horizon/web_v2/resources"

	regources "gitlab.com/tokend/regources/generated"

	"gitlab.com/tokend/horizon/db2/history2"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/web_v2/ctx"
	"gitlab.com/tokend/horizon/web_v2/requests"
)

func GetSaleWhitelist(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetSaleWhitelist(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	handler := getSaleWhiteListHandler{
		SalesQ:                history2.NewSalesQ(ctx.HistoryRepo(r)),
		AccountSpecificRulesQ: history2.NewAccountSpecificRulesQ(ctx.HistoryRepo(r)),
		Log: ctx.Log(r),
	}

	sale, err := handler.SalesQ.GetByID(request.SaleID)
	if err != nil {
		ctx.Log(r).WithError(err).WithFields(logan.F{
			"sale_id": request.SaleID,
		}).Error("failed to get sale by ID")
	}
	if sale == nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}

	if !isAllowed(r, w, sale.OwnerAddress) {
		return
	}

	result, err := handler.getSaleWhiteList(request)
	if errObj, ok := err.(*jsonapi.ErrorObject); ok {
		ape.RenderErr(w, errObj)
		return
	}
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get whitelist", logan.F{
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

type getSaleWhiteListHandler struct {
	SalesQ                history2.SalesQ
	AccountSpecificRulesQ history2.AccountSpecificRulesQ
	Log                   *logan.Entry
}

// GetSale returns sale with related resources
func (h *getSaleWhiteListHandler) getSaleWhiteList(request *requests.GetSaleWhitelist) (*regources.SaleWhitelistsResponse, error) {
	rules, err := h.AccountSpecificRulesQ.
		ForSale(request.SaleID).
		Permission(false).
		Page(*request.PageParams).
		Select()

	if err != nil {
		return nil, errors.Wrap(err, "failed to get account specific rules for sale")
	}

	if rules == nil {
		return nil, nil
	}

	response := &regources.SaleWhitelistsResponse{}

	for _, rule := range rules {
		if rule.Address == nil {
			continue
		}

		response.Data = append(response.Data, resources.NewSaleWhitelist(request.SaleID, *rule.Address))
	}

	h.populateLinks(response, request)

	return response, nil
}

func (h *getSaleWhiteListHandler) populateLinks(
	response *regources.SaleWhitelistsResponse, request *requests.GetSaleWhitelist,
) {
	if len(response.Data) > 0 {
		response.Links = request.GetCursorLinks(*request.PageParams, response.Data[len(response.Data)-1].ID)
	} else {
		response.Links = request.GetCursorLinks(*request.PageParams, "")
	}
}
