package handlers

import (
	"net/http"

	"gitlab.com/tokend/regources/v2"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/horizon/web_v2/ctx"
	"gitlab.com/tokend/horizon/web_v2/requests"
)

func GetRequests(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetRequests(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	historyRepo := ctx.HistoryRepo(r)
	handler := getRequestListHandler{
		R:         request,
		RequestsQ: history2.NewReviewableRequestsQ(historyRepo),
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

type getRequestListHandler struct {
	R         requests.GetRequests
	Base      getRequestListBaseHandler
	RequestsQ history2.ReviewableRequestsQ
	Log       *logan.Entry
}

func (h *getRequestListHandler) MakeAll(w http.ResponseWriter, request requests.GetRequests) error {
	q := h.RequestsQ

	// apply custom filters here

	return h.Base.SelectAndRender(w, *request.GetRequestsBase, q, h.RenderRecord)
}

func (h *getRequestListHandler) RenderRecord(included *regources.Included, record history2.ReviewableRequest) (regources.ReviewableRequest, error) {
	resource := h.Base.PopulateResource(*h.R.GetRequestsBase, included, record)

	// apply custom includes here

	return resource, nil
}