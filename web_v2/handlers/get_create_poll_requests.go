package handlers

import (
	"net/http"

	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/core2"
	regources "gitlab.com/tokend/regources/generated"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/horizon/web_v2/ctx"
	"gitlab.com/tokend/horizon/web_v2/requests"
)

func GetCreatePollRequests(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetCreatePollRequests(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	historyRepo := ctx.HistoryRepo(r)
	coreRepo := ctx.CoreRepo(r)
	handler := getCreatePollRequestsHandler{
		R:         request,
		RequestsQ: history2.NewReviewableRequestsQ(historyRepo),
		AccountsQ: core2.NewAccountsQ(coreRepo),
		Log:       ctx.Log(r),
	}

	constraints := []string{
		request.GetRequestsBase.Filters.Requestor[0],
		request.GetRequestsBase.Filters.Reviewer[0],
	}

	if !isAllowed(r, w, constraints...) {
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

type getCreatePollRequestsHandler struct {
	R         requests.GetCreatePollRequests
	Base      getRequestListBaseHandler
	RequestsQ history2.ReviewableRequestsQ
	AccountsQ core2.AccountsQ
	Log       *logan.Entry
}

func (h *getCreatePollRequestsHandler) MakeAll(w http.ResponseWriter, request requests.GetCreatePollRequests) error {
	q := h.RequestsQ.FilterByRequestType(uint64(xdr.ReviewableRequestTypeCreatePoll))

	if request.ShouldFilter(requests.FilterTypeCreatePollRequestsPermissionType) {
		q = q.FilterByCreatePollPermissionType(request.Filters.PermissionType[0])
	}

	if request.ShouldFilter(requests.FilterTypeCreatePollRequestsVoteConfirmationRequired) {
		q = q.FilterByCreatePollVoteConfirmationRequired(request.Filters.VoteConfirmationRequired[0])
	}
	if request.ShouldFilter(requests.FilterTypeCreatePollRequestsResultProvider) {
		q = q.FilterByCreatePollResultProvider(request.Filters.ResultProvider[0])
	}

	return h.Base.SelectAndRender(w, *request.GetRequestsBase, q, h.RenderRecord)
}

func (h *getCreatePollRequestsHandler) RenderRecord(included *regources.Included, record history2.ReviewableRequest) (regources.ReviewableRequest, error) {
	resource := h.Base.PopulateResource(*h.R.GetRequestsBase, included, record)

	return resource, nil
}
