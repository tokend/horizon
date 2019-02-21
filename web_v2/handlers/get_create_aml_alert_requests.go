package handlers

import (
	"net/http"

	"gitlab.com/tokend/horizon/db2/core2"
	"gitlab.com/tokend/horizon/web_v2/resources"

	"gitlab.com/tokend/go/xdr"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/horizon/web_v2/ctx"
	"gitlab.com/tokend/horizon/web_v2/requests"
	"gitlab.com/tokend/regources/v2"
)

func GetCreateAmlAlertRequests(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetCreateAmlAlertRequests(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	historyRepo := ctx.HistoryRepo(r)
	coreRepo := ctx.CoreRepo(r)
	handler := getCreateAmlAlertRequestsHandler{
		R:         request,
		RequestsQ: history2.NewReviewableRequestsQ(historyRepo),
		BalancesQ: core2.NewBalancesQ(coreRepo),
		Log:       ctx.Log(r),
	}

	if !isAllowed(r, w, request.Filters.Requestor, request.Filters.Reviewer) {
		return
	}

	err = handler.MakeAll(w, request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get request list", logan.F{
			"request": request,
		})
		ape.RenderErr(w, problems.InternalError())
		return
	}
}

type getCreateAmlAlertRequestsHandler struct {
	R         requests.GetCreateAmlAlertRequests
	Base      getRequestListBaseHandler
	RequestsQ history2.ReviewableRequestsQ
	BalancesQ core2.BalancesQ
	Log       *logan.Entry
}

func (h *getCreateAmlAlertRequestsHandler) MakeAll(w http.ResponseWriter, request requests.GetCreateAmlAlertRequests) error {
	q := h.RequestsQ.FilterByRequestType(uint64(xdr.ReviewableRequestTypeCreateAmlAlert))

	if request.ShouldFilter(requests.FilterTypeCreateAmlAlertRequestsBalance) {
		q = q.FilterByAmlAlertBalance(request.Filters.Balance)
	}

	return h.Base.SelectAndRender(w, *request.GetRequestsBase, q, h.RenderRecord)
}

func (h *getCreateAmlAlertRequestsHandler) RenderRecord(included *regources.Included, record history2.ReviewableRequest) (regources.ReviewableRequest, error) {
	resource := h.Base.PopulateResource(*h.R.GetRequestsBase, included, record)

	if h.R.ShouldInclude(requests.IncludeTypeCreateAmlAlertRequestsBalance) {
		balance, err := h.BalancesQ.GetByAddress(record.Details.CreateAmlAlert.BalanceID)
		if err != nil {
			return regources.ReviewableRequest{}, errors.Wrap(err, "failed to get balance")
		}

		if balance == nil {
			return regources.ReviewableRequest{}, errors.New("balance not found")
		}
		resource := resources.NewBalance(balance)
		included.Add(resource)
	}

	return resource, nil
}