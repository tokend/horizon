package handlers

import (
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/horizon/web_v2/ctx"
	"net/http"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/tokend/horizon/web_v2/requests"
)

func GetRequestForAccount(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetRequestForAccount(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	historyRepo := ctx.HistoryRepo(r)
	handler := getRequestListHandler{
		R:         *request,
		Log:       ctx.Log(r),
		RequestsQ: history2.NewReviewableRequestsQ(historyRepo),
	}

	if !isAllowed(r, w, request.GetRequestsBase.Filters.Requestor, request.GetRequestsBase.Filters.Reviewer) {
		return
	}

	if err = handler.MakeAll(w, *request); err != nil {
		ctx.Log(r).WithError(err).Error("failed to get request list", logan.F{
			"request": request,
		})
		ape.RenderErr(w, problems.InternalError())
		return
	}
}
