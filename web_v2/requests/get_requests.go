package requests

import (
	"net/http"

	"gitlab.com/distributed_lab/urlval"
)

const (
	// IncludeTypeReviewableRequestListDetails - defines if we need to include request details to response
	IncludeTypeReviewableRequestListDetails = "request_details"

	// FilterTypeRequestListParticipant defines if we need to filter the list by participant (AccountID)
	FilterTypeRequestListParticipant = "participant"
	// FilterTypeRequestListRequestor - defines if we need to filter the list by requestor
	FilterTypeRequestListRequestor = "requestor"
	// FilterTypeRequestListReviewer - defines if we need to filter the list by reviewer
	FilterTypeRequestListReviewer = "reviewer"
	// FilterTypeRequestListState - defines if we need to filter the list by request state
	FilterTypeRequestListState = "state"
	// FilterTypeRequestListType - defines if we need to filter the list by request type
	FilterTypeRequestListType = "type"
	// FilterTypeRequestListPendingTasks - defines if we need to filter the list by pending tasks
	FilterTypeRequestListPendingTasks = "pending_tasks"
	// FilterTypeRequestListPendingTasksNotSet - defines if we need to filter the list by pending tasks that aren't set
	FilterTypeRequestListPendingTasksNotSet = "pending_tasks_not_set"
	// FilterTypeRequestListPendingTasksAnyOf - defines if we need to filter the list by any of pending tasks
	FilterTypeRequestListPendingTasksAnyOf = "pending_tasks_any_of"
	// FilterTypeRequestListCreatedBefore - defines if we need to filter the list by creation time before specific timestamp
	FilterTypeRequestListCreatedBefore = "created_before"
	// FilterTypeRequestListCreatedAfter - defines if we need to filter the list by creation time after specific timestamp
	FilterTypeRequestListCreatedAfter = "created_after"
	// FilterTypeRequestListAllTasks - defines if we need to filter the list by all tasks
	FilterTypeRequestListAllTasks = "all_tasks"
	// FilterTypeRequestListAllTasksAnyOf - defines if we need to filter the list by any of all tasks
	FilterTypeRequestListAllTasksAnyOf = "all_tasks_any_of"
	// FilterTypeRequestListAllTasksNotSet - defines if we need to filter the list by all tasks that aren't set
	FilterTypeRequestListAllTasksNotSet = "all_tasks_not_set"
	// FilterTypeRequestListUpdatedBefore - defines if we need to filter the list by update time before specific timestamp
	FilterTypeRequestListUpdatedBefore = "updated_before"
	// FilterTypeRequestListUpdatedAfter - defines if we need to filter the list by update time after specific timestamp
	FilterTypeRequestListUpdatedAfter = "updated_after"
)

var includeTypeReviewableRequestListAll = map[string]struct{}{
	IncludeTypeReviewableRequestListDetails: {},
}

var filterTypeRequestListAll = map[string]struct{}{
	FilterTypeRequestListParticipant:        {},
	FilterTypeRequestListRequestor:          {},
	FilterTypeRequestListReviewer:           {},
	FilterTypeRequestListState:              {},
	FilterTypeRequestListType:               {},
	FilterTypeRequestListPendingTasks:       {},
	FilterTypeRequestListPendingTasksNotSet: {},
	FilterTypeRequestListPendingTasksAnyOf:  {},
	FilterTypeRequestListCreatedBefore:      {},
	FilterTypeRequestListCreatedAfter:       {},
	FilterTypeRequestListAllTasks:           {},
	FilterTypeRequestListAllTasksNotSet:     {},
	FilterTypeRequestListAllTasksAnyOf:      {},
	FilterTypeRequestListUpdatedBefore:      {},
	FilterTypeRequestListUpdatedAfter:       {},
}

type GetRequests struct {
	GetRequestsBase
}

func NewGetRequests(r *http.Request) (request GetRequests, err error) {
	request.GetRequestsBase, err = NewGetRequestsBase(
		r,
		map[string]struct{}{},
		map[string]struct{}{},
	)
	if err != nil {
		return request, err
	}

	err = urlval.DecodeSilently(r.URL.Query(), &request)
	if err != nil {
		return request, err
	}

	err = PopulateRequest(&request.GetRequestsBase)
	if err != nil {
		return request, err
	}

	return request, nil
}
