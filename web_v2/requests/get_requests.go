package requests

import (
	"net/http"
)

const (
	// IncludeTypeReviewableRequestListDetails - defines if we need to include request details to response
	IncludeTypeReviewableRequestListDetails = "request_details"

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
)

var includeTypeReviewableRequestListAll = map[string]struct{}{
	IncludeTypeReviewableRequestListDetails: {},
}

var filterTypeRequestListAll = map[string]struct{}{
	FilterTypeRequestListRequestor:          {},
	FilterTypeRequestListReviewer:           {},
	FilterTypeRequestListState:              {},
	FilterTypeRequestListType:               {},
	FilterTypeRequestListPendingTasks:       {},
	FilterTypeRequestListPendingTasksNotSet: {},
	FilterTypeRequestListPendingTasksAnyOf:  {},
}

// GetReviewableRequestList represents params to be specified by user for getReviewableRequestList handler
//type GetReviewableBaseRequestList struct {
//	*base
//	BaseFilters GetReviewableRequestListFilters
//	PageParams  *pgdb.CursorPageParams
//}

type GetRequestListBaseFilters struct {
	ID                  []uint64 `filter:"id"`
	Requestor           []string `filter:"requestor"`
	Reviewer            []string `filter:"reviewer"`
	State               []uint64 `filter:"state"`
	Type                []uint64 `filter:"type"`
	PendingTasks        []uint64 `filter:"pending_tasks"`
	PendingTasksAnyOf   []uint64 `filter:"pending_tasks_any_of"`
	PendingTasksNotSet  []uint64 `filter:"pending_tasks_not_set"`
	MissingPendingTasks []uint64 `filter:"missing_pending_tasks"`
}
type GetRequests struct {
	*GetRequestsBase
	Filters GetRequestListBaseFilters
}

func NewGetRequests(r *http.Request) (request GetRequests, err error) {
	request.GetRequestsBase, err = NewGetRequestsBase(
		r,
		&request.Filters,
		map[string]struct{}{},
		map[string]struct{}{},
	)
	if err != nil {
		return request, err
	}

	return request, nil
}
