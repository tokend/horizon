package handlers

import (
	"net/http"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3/errors"

	"gitlab.com/tokend/horizon/db2/history2"
	"gitlab.com/tokend/horizon/web_v2/requests"
	"gitlab.com/tokend/horizon/web_v2/resources"
	regources "gitlab.com/tokend/regources/generated"
)

type getRequestListBaseHandler struct {
}

func (h *getRequestListBaseHandler) PopulateLinks(
	response *regources.ReviewableRequestListResponse, request requests.GetRequestsBase,
) {
	if len(response.Data) > 0 {
		response.Links = request.GetCursorLinks(*request.PageParams, response.Data[len(response.Data)-1].ID)
	} else {
		response.Links = request.GetCursorLinks(*request.PageParams, "")
	}
}

func (h *getRequestListBaseHandler) SelectAndRender(
	w http.ResponseWriter,
	request requests.GetRequestsBase,
	requestsQ history2.ReviewableRequestsQ,
	renderer func(*regources.Included, history2.ReviewableRequest) (regources.ReviewableRequest, error),
) error {

	q := h.ApplyFilters(request, requestsQ)

	records, err := q.Select()
	if err != nil {
		return errors.Wrap(err, "Failed to get reviewable request list")
	}

	if request.Filters.ID[0] != 0 {
		if len(records) == 0 {
			ape.RenderErr(w, problems.NotFound())
			return nil
		}

		var response regources.ReviewableRequestResponse
		response.Data, err = renderer(&response.Included, records[0])
		if err != nil {
			return errors.Wrap(err, "failed to render record")
		}

		ape.Render(w, response)
		return nil
	} else {
		response := &regources.ReviewableRequestListResponse{
			Data: make([]regources.ReviewableRequest, 0, len(records)),
		}

		for _, record := range records {
			resource, err := renderer(&response.Included, record)
			if err != nil {
				return errors.Wrap(err, "failed to render record")
			}
			response.Data = append(response.Data, resource)
		}

		h.PopulateLinks(response, request)

		ape.Render(w, response)
		return nil
	}
}

func (h *getRequestListBaseHandler) PopulateResource(
	request requests.GetRequestsBase, included *regources.Included, record history2.ReviewableRequest,
) regources.ReviewableRequest {
	reviewableRequest := resources.NewRequest(record)
	reviewableRequestDetails := resources.NewRequestDetails(record)
	reviewableRequest.Relationships.RequestDetails = reviewableRequestDetails.GetKey().AsRelation()

	if request.ShouldInclude(requests.IncludeTypeReviewableRequestListDetails) {
		included.Add(reviewableRequestDetails)
	}
	return reviewableRequest
}

func (h *getRequestListBaseHandler) ApplyFilters(
	request requests.GetRequestsBase, q history2.ReviewableRequestsQ,
) history2.ReviewableRequestsQ {
	q = q.Page(*request.PageParams)
	if request.ShouldFilter(requests.FilterTypeRequestListRequestor) {
		q = q.FilterByRequestorAddress(request.Filters.Requestor[0])
	}

	if request.ShouldFilter(requests.FilterTypeRequestListReviewer) {
		q = q.FilterByReviewerAddress(request.Filters.Reviewer[0])
	}

	if request.ShouldFilter(requests.FilterTypeRequestListState) {
		q = q.FilterByState(request.Filters.State[0])
	}

	if request.ShouldFilter(requests.FilterTypeRequestListType) {
		q = q.FilterByRequestType(request.Filters.Type[0])
	}

	if request.ShouldFilter(requests.FilterTypeRequestListPendingTasks) {
		q = q.FilterByPendingTasks(request.Filters.PendingTasks[0])
	}

	if request.ShouldFilter(requests.FilterTypeRequestListPendingTasksNotSet) {
		q = q.FilterPendingTasksNotSet(request.Filters.PendingTasksNotSet[0])
	}

	if request.ShouldFilter(requests.FilterTypeRequestListPendingTasksAnyOf) {
		q = q.FilterByPendingTasksAnyOf(request.Filters.PendingTasksAnyOf[0])
	}

	if request.Filters.ID[0] != 0 {
		q = q.FilterByID(request.Filters.ID[0])
	}

	return q
}
