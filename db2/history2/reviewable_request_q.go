package history2

import (
	sq "github.com/lann/squirrel"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2"
)

// ReviewableRequestsQ - helper struct to get reviewable requests from db
type ReviewableRequestsQ struct {
	repo     *db2.Repo
	selector sq.SelectBuilder
}

// NewReviewableRequestsQ - creates new instance of ReviewableRequestsQ
func NewReviewableRequestsQ(repo *db2.Repo) ReviewableRequestsQ {
	return ReviewableRequestsQ{
		repo: repo,
		selector: sq.Select(
			"reviewable_requests.id",
			"reviewable_requests.requestor",
			"reviewable_requests.reviewer",
			"reviewable_requests.reference",
			"reviewable_requests.reject_reason",
			"reviewable_requests.request_type",
			"reviewable_requests.request_state",
			"reviewable_requests.hash",
			"reviewable_requests.details",
			"reviewable_requests.created_at",
			"reviewable_requests.updated_at",
			"reviewable_requests.all_tasks",
			"reviewable_requests.pending_tasks",
			"reviewable_requests.external_details",
		).From("reviewable_request reviewable_requests"),
	}
}

// FilterByState - returns q with filter by request state
func (q ReviewableRequestsQ) FilterByState(state uint64) ReviewableRequestsQ { // TODO: move state to regources?
	q.selector = q.selector.Where("reviewable_requests.request_state = ?", state)
	return q
}

// FilterByRequestorAddress - returns q with filter by requestor address
func (q ReviewableRequestsQ) FilterByRequestorAddress(address string) ReviewableRequestsQ {
	q.selector = q.selector.Where("reviewable_requests.requestor = ?", address)
	return q
}

// FilterByReviewerAddress - returns q with filter by reviewer
func (q ReviewableRequestsQ) FilterByReviewerAddress(address string) ReviewableRequestsQ {
	q.selector = q.selector.Where("reviewable_requests.reviewer", address)
	return q
}

// FilterByPendingTasks - returns q with filter by pending tasks
func (q ReviewableRequestsQ) FilterByPendingTasks(mask uint64) ReviewableRequestsQ {
	q.selector = q.selector.Where("(details->'update_kyc'->>'pending_tasks')::integer & ? = ?", mask, mask)
	return q
}

// FilterByPendingTasksAnyOf - returns q with filter by pending tasks
func (q ReviewableRequestsQ) FilterByPendingTasksAnyOf(mask uint64) ReviewableRequestsQ {
	q.selector = q.selector.Where("(details->'update_kyc'->>'pending_tasks')::integer & ? <> 0", mask)
	return q
}

// FilterPendingTasksNotSet - returns q with filter by pending tasks that aren't set
func (q ReviewableRequestsQ) FilterPendingTasksNotSet(mask uint64) ReviewableRequestsQ {
	q.selector = q.selector.Where("~(details->'update_kyc'->>'pending_tasks')::integer & ? = ?", mask, mask)
	return q
}

// FilterByRequestType - returns q with filter by request type
func (q ReviewableRequestsQ) FilterByRequestType(requestType uint64) ReviewableRequestsQ {
	q.selector = q.selector.Where("request_type = ?", requestType)
	return q
}

// FilterByID - returns q with filter by ID
func (q ReviewableRequestsQ) FilterByID(id uint64) ReviewableRequestsQ {
	q.selector = q.selector.Where("reviewable_requests.id = ?", id)
	return q
}

// GetByID loads a row from `reviewable_requests`, by ID
// returns nil, nil - if request does not exists
func (q ReviewableRequestsQ) GetByID(id uint64) (*ReviewableRequest, error) {
	return q.FilterByID(id).Get()
}

// Page - apply paging params to the query
func (q ReviewableRequestsQ) Page(pageParams db2.CursorPageParams) ReviewableRequestsQ {
	q.selector = pageParams.ApplyTo(q.selector, "reviewable_requests.id")
	return q
}

// Get - loads a row from `reviewable_requests`
// returns nil, nil - if request does not exists
// returns error if more than one ReviewableRequest found
func (q ReviewableRequestsQ) Get() (*ReviewableRequest, error) {
	var result ReviewableRequest
	err := q.repo.Get(&result, q.selector)
	if err != nil {
		if q.repo.NoRows(err) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load account")
	}

	return &result, nil
}

// Select - selects ReviewableRequest from db using specified filters. Returns nil, nil - if one does not exists
func (q ReviewableRequestsQ) Select() ([]ReviewableRequest, error) {
	var result []ReviewableRequest
	err := q.repo.Select(&result, q.selector)
	if err != nil {
		if q.repo.NoRows(err) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to select reviewable requests")
	}

	return result, nil
}
