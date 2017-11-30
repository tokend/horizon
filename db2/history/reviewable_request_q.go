package history

import (
	sq "github.com/lann/squirrel"
)

// ReviewableRequestQI - provides methods to operate reviewable request
type ReviewableRequestQI interface {
	// Insert - inserts new request
	Insert(request ReviewableRequest) error
	// Update - update request using it's ID
	Update(request ReviewableRequest) error
	// ByID - selects request by id. Returns nil, nil if not found
	ByID(requestID uint64) (*ReviewableRequest, error)
	// Cancel - sets request state to `ReviewableRequestStateCanceled`
	Cancel(requestID uint64) error
	// Approve - sets request state to ReviewableRequestStateApproved and cleans reject reason
	Approve(requestID uint64) error
	// PermanentReject - sets request state to ReviewableRequestStatePermanentlyRejected and sets reject reason
	PermanentReject(requestID uint64, rejectReason string) error
}

type ReviewableRequestQ struct {
	Err    error
	parent *Q
	sql    sq.SelectBuilder
}

// Insert - inserts new request
func (q *ReviewableRequestQ) Insert(request ReviewableRequest) error {
	if q.Err != nil {
		return q.Err
	}

	query := sq.Insert("reviewable_request").Columns("id", "requestor", "reviewer", "reference", "reject_reason",
		"request_type", "request_state", "hash", "details").Values(
		request.ID,
		request.Requestor,
		request.Reviewer,
		request.Reference,
		request.RejectReason,
		request.RequestType,
		request.RequestState,
		request.Hash,
		request.Details,
	)

	_, err := q.parent.Exec(query)
	return err
}

// Update - update request using it's ID
func (q *ReviewableRequestQ) Update(request ReviewableRequest) error {
	if q.Err != nil {
		return q.Err
	}

	query := sq.Update("reviewable_request").SetMap(map[string]interface{}{
		"requestor":     request.Requestor,
		"reviewer":      request.Reviewer,
		"reject_reason": request.RejectReason,
		"request_type":  request.RequestType,
		"request_state": request.RequestState,
		"hash":          request.Hash,
		"details":       request.Details,
	})

	_, err := q.parent.Exec(query)
	return err
}

// ByID - selects request by id. Returns nil, nil if not found
func (q *ReviewableRequestQ) ByID(requestID uint64) (*ReviewableRequest, error) {
	if q.Err != nil {
		return nil, q.Err
	}

	query := sq.Select("id", "requestor", "reviewer", "reference", "reject_reason", "request_type", "request_state", "hash",
		"details").From("reviewable_request").Where("id = ?", requestID)

	var result ReviewableRequest
	err := q.parent.Get(&result, query)
	if q.parent.NoRows(err) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &result, nil
}

// Cancel - sets request state to `ReviewableRequestStateCanceled`
func (q *ReviewableRequestQ) Cancel(requestID uint64) error {
	if q.Err != nil {
		return q.Err
	}

	query := sq.Update("reviewable_request").
		Set("request_state", ReviewableRequestStateCanceled).Where("id = ?", requestID)

	_, err := q.parent.Exec(query)
	return err
}

// Approve - sets request state to ReviewableRequestStateApproved and cleans reject reason
func (q *ReviewableRequestQ) Approve(requestID uint64) error {
	return q.setStateRejectReason(requestID, ReviewableRequestStateApproved, "")
}

// PermanentReject - sets request state to ReviewableRequestStatePermanentlyRejected and sets reject reason
func (q *ReviewableRequestQ) PermanentReject(requestID uint64, rejectReason string) error {
	return q.setStateRejectReason(requestID, ReviewableRequestStatePermanentlyRejected, rejectReason)
}

func (q *ReviewableRequestQ) setStateRejectReason(requestID uint64, requestState ReviewableRequestState, rejectReason string) error {
	if q.Err != nil {
		return q.Err
	}

	query := sq.Update("reviewable_request").
		Set("request_state", requestState).
		Set("reject_reason", rejectReason).Where("id = ?", requestID)

	_, err := q.parent.Exec(query)
	return err
}
