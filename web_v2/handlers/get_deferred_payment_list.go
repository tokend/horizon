package handlers

import (
	"net/http"

	"gitlab.com/tokend/horizon/web_v2/resources"

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

func GetDeferredPaymentList(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetDeferredPaymentList(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	hRepo := ctx.HistoryRepo(r)
	handler := getDeferredPaymentListHandler{
		DeferredPaymentQ: history2.NewDeferredPaymentQ(hRepo),
		Log:              ctx.Log(r),
	}

	deferredPaymentOwners := []*string{}

	if request.ShouldFilter(requests.FilterTypeDeferredPaymentListSourceBalance) {
		source, err := handler.BalanceQ.GetByAddress(request.Filters.SourceBalance)
		if err != nil {
			ape.RenderErr(w, problems.NotFound())
			return
		}

		deferredPaymentOwners = append(deferredPaymentOwners, &source.AccountAddress)
	}

	if request.ShouldFilter(requests.FilterTypeDeferredPaymentListDestination) {
		deferredPaymentOwners = append(deferredPaymentOwners, &request.Filters.Destination)
	}

	if !isAllowed(r, w, deferredPaymentOwners...) {
		return
	}

	response, err := handler.GetDeferredPaymentList(request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get deferredPayment")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, response)
}

type getDeferredPaymentListHandler struct {
	DeferredPaymentQ history2.DeferredPaymentQ
	AccountQ         core2.AccountsQ
	BalanceQ         core2.BalancesQ
	Log              *logan.Entry
}

func (h *getDeferredPaymentListHandler) GetDeferredPaymentList(request *requests.GetDeferredPaymentList) (*regources.DeferredPaymentListResponse, error) {
	q := h.DeferredPaymentQ

	if request.ShouldFilter(requests.FilterTypeDeferredPaymentListSourceBalance) {
		q = q.FilterBySourceBalance(request.Filters.SourceBalance)
	}

	if request.ShouldFilter(requests.FilterTypeDeferredPaymentListDestination) {
		q = q.FilterByDestinationAccount(request.Filters.Destination)
	}

	q = q.Page(request.PageParams)

	deferredPaymentSet, err := q.Select()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get deferredPayment list")
	}

	response := regources.DeferredPaymentListResponse{
		Data: make([]regources.DeferredPayment, 0, len(deferredPaymentSet)),
	}

	for _, deferredPaymentEntry := range deferredPaymentSet {
		response.Data = append(response.Data, resources.NewDeferredPayment(deferredPaymentEntry))
	}

	if len(response.Data) > 0 {
		response.Links = request.GetCursorLinks(*request.PageParams, response.Data[len(response.Data)-1].ID)
	} else {
		response.Links = request.GetCursorLinks(*request.PageParams, "")
	}

	return &response, nil
}
