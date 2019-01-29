package requests

import (
	"gitlab.com/tokend/horizon/db2"
	"net/http"
)

const (
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
type GetReviewableRequestList struct {
	*base
	Filters struct {
		Requestor           string `fig:"requestor"`
		Reviewer            string `fig:"reviewer"`
		State               uint64 `fig:"state"`
		Type                uint64 `fig:"type"`
		PendingTasks        uint64 `fig:"pending_tasks"`
		PendingTasksAnyOf   uint64 `fig:"pending_tasks_any_of"`
		PendingTasksNotSet  uint64 `fig:"pending_tasks_not_set"`
		MissingPendingTasks uint64 `fig:"missing_pending_tasks"`
	}

	PageParams *db2.CursorPageParams
}

// NewGetReviewableRequestList - returns new instance of GetReviewableRequestList
func NewGetReviewableRequestList(r *http.Request) (*GetReviewableRequestList, error) {
	b, err := newBase(r, baseOpts{
		supportedFilters: filterTypeRequestListAll,
	})
	if err != nil {
		return nil, err
	}

	pageParams, err := b.getCursorBasedPageParams()
	if err != nil {
		return nil, err
	}

	request := GetReviewableRequestList{
		base:       b,
		PageParams: pageParams,
	}

	err = b.populateFilters(&request.Filters)
	if err != nil {
		return nil, err
	}

	return &request, nil
}
