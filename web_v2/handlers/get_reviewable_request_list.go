package handlers

import (
	"net/http"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/horizon/web_v2/ctx"
	"gitlab.com/tokend/horizon/web_v2/requests"
	"gitlab.com/tokend/horizon/web_v2/resources"
	"gitlab.com/tokend/regources/v2"
)

// GetReviewableRequestList - processes request to get the list of reviewable requests and their details
func GetReviewableRequestList(w http.ResponseWriter, r *http.Request) {
	historyRepo := ctx.HistoryRepo(r)
	handler := getReviewableRequestListHandler{
		RequestsQ: history2.NewReviewableRequestsQ(historyRepo),
		Log:       ctx.Log(r),
	}

	request, err := requests.NewGetReviewableRequestList(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	// TODO check allowed?

	result, err := handler.GetReviewableRequestList(request)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get request list", logan.F{
			"request": request,
		})
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, result)
}

type getReviewableRequestListHandler struct {
	RequestsQ history2.ReviewableRequestsQ
	Log       *logan.Entry
}

// GetReviewableRequestList returns reviewable request list with related resources
func (h *getReviewableRequestListHandler) GetReviewableRequestList(
	request *requests.GetReviewableRequestList) (*regources.ReviewableRequestsResponse, error) {

	q := h.RequestsQ.Page(*request.PageParams)

	if request.ShouldFilter(requests.FilterTypeRequestListRequestor) {
		q = q.FilterByRequestorAddress(request.Filters.Requestor)
	}

	if request.ShouldFilter(requests.FilterTypeRequestListReviewer) {
		q = q.FilterByReviewerAddress(request.Filters.Reviewer)
	}

	if request.ShouldFilter(requests.FilterTypeRequestListState) {
		q = q.FilterByState(request.Filters.State)
	}

	if request.ShouldFilter(requests.FilterTypeRequestListType) {
		q = q.FilterByRequestType(request.Filters.Type)
	}

	if request.ShouldFilter(requests.FilterTypeRequestListPendingTasks) {
		q = q.FilterByPendingTasks(request.Filters.PendingTasks)
	}

	if request.ShouldFilter(requests.FilterTypeRequestListPendingTasksNotSet) {
		q = q.FilterPendingTasksNotSet(request.Filters.PendingTasksNotSet)
	}

	if request.ShouldFilter(requests.FilterTypeRequestListPendingTasksAnyOf) {
		q = q.FilterByPendingTasksAnyOf(request.Filters.PendingTasksAnyOf)
	}

	historyRecords, err := q.Select()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get reviewable request list")
	}

	response := &regources.ReviewableRequestsResponse{
		Data: make([]regources.ReviewableRequest, 0, len(historyRecords)),
	}

	for _, historyRecord := range historyRecords {
		reviewableRequest := resources.NewRequest(historyRecord)
		reviewableRequestDetails := resources.NewRequestDetails(historyRecord)
		reviewableRequest.Relationships.RequestDetails = reviewableRequestDetails.GetKey().AsRelation()

		response.Data = append(response.Data, reviewableRequest)

		if request.ShouldInclude(requests.IncludeTypeReviewableRequestListDetails) {
			response.Included.Add(reviewableRequestDetails)
		}
	}

	if len(response.Data) > 0 {
		response.Links = request.GetCursorLinks(*request.PageParams, response.Data[len(response.Data)-1].ID)
	} else {
		response.Links = request.GetCursorLinks(*request.PageParams, "")
	}

	return response, nil
}
