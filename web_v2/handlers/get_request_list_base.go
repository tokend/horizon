package handlers

import (
	"math"
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
		response.Links = request.GetCursorLinks(request.PageParams, response.Data[len(response.Data)-1].ID)
	} else {
		response.Links = request.GetCursorLinks(request.PageParams, "")
	}
}

func (h *getRequestListBaseHandler) SelectAndRender(
	w http.ResponseWriter,
	request requests.GetRequestsBase,
	requestsQ history2.ReviewableRequestsQ,
	renderer func(*regources.Included, history2.ReviewableRequest) (regources.ReviewableRequest, error),
) error {

	q := h.ApplyFilters(request, requestsQ)

	recordsAll, err := q.Select()
	if err != nil {
		return errors.Wrap(err, "failed to get reviewable request list")
	}

	if request.UseOffsetPageParams {
		q = q.PageOffset(request.OffsetPageParams)
	} else {
		q = q.Page(request.PageParams)
	}

	records, err := q.Select()
	if err != nil {
		return errors.Wrap(err, "Failed to get reviewable request list")
	}

	if request.Filters.ID != nil && *request.Filters.ID != 0 {
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

		if request.UseOffsetPageParams {
			response.Links = request.GetOffsetLinks(request.OffsetPageParams)

			err = response.PutMeta(requests.MetaPageParams{
				CurrentPage: request.OffsetPageParams.PageNumber,
				TotalPages:  uint64(math.Ceil(float64(len(recordsAll)) / float64(request.OffsetPageParams.Limit))),
			})
		} else {
			h.PopulateLinks(response, request)

			err = response.PutMeta(requests.MetaCursorParams{
				CurrentCursor: request.PageParams.Cursor,
				TotalPages:    uint64(math.Ceil(float64(len(recordsAll)) / float64(request.PageParams.Limit))),
			})
		}

		if err != nil {
			return errors.Wrap(err, "failed to put meta to response")
		}

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
	if request.Filters.Requestor != nil {
		q = q.FilterByRequestorAddress(*request.Filters.Requestor)
	}

	if request.Filters.Reviewer != nil {
		q = q.FilterByReviewerAddress(*request.Filters.Reviewer)
	}

	if request.Filters.State != nil {
		q = q.FilterByState(*request.Filters.State)
	}

	if request.Filters.Type != nil {
		q = q.FilterByRequestType(*request.Filters.Type)
	}

	if request.Filters.PendingTasks != nil {
		q = q.FilterByPendingTasks(*request.Filters.PendingTasks)
	}

	if request.Filters.PendingTasksNotSet != nil {
		q = q.FilterPendingTasksNotSet(*request.Filters.PendingTasksNotSet)
	}

	if request.Filters.PendingTasksAnyOf != nil {
		q = q.FilterByPendingTasksAnyOf(*request.Filters.PendingTasksAnyOf)
	}

	if request.Filters.CreatedAfter != nil {
		q = q.FilterByCreatedAtAfter(*request.Filters.CreatedAfter)
	}

	if request.Filters.CreatedBefore != nil {
		q = q.FilterByCreatedAtBefore(*request.Filters.CreatedBefore)
	}

	if request.Filters.ID != nil && *request.Filters.ID != 0 {
		q = q.FilterByID(*request.Filters.ID)
	}

	if request.Filters.AllTasks != nil {
		q = q.FilterByAllTasks(*request.Filters.AllTasks)
	}

	if request.Filters.AllTasksNotSet != nil {
		q = q.FilterByAllTasksNotSet(*request.Filters.AllTasksNotSet)
	}

	if request.Filters.AllTasksAnyOf != nil {
		q = q.FilterByAllTasksAnyOf(*request.Filters.AllTasksAnyOf)
	}

	if request.Filters.UpdatedAfter != nil {
		q = q.FilterByUpdatedAtAfter(*request.Filters.UpdatedAfter)
	}

	if request.Filters.UpdatedBefore != nil {
		q = q.FilterByUpdatedAtBefore(*request.Filters.UpdatedBefore)
	}

	return q
}
