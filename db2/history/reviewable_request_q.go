package history

import (
	sq "github.com/lann/squirrel"
	"gitlab.com/swarmfund/horizon/db2"
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/horizon/db2/sqx"
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
	// ForRequestor - filters requests by requestor
	ForRequestor(requestor string) ReviewableRequestQI
	// ForReviewer - filters requests by reviewer
	ForReviewer(reviewer string) ReviewableRequestQI
	// ForState - filters requests by state
	ForState(state int64) ReviewableRequestQI
	// ForType - filters requests by type
	ForType(requestType int64) ReviewableRequestQI
	// ForAsset - filters requests by asset
	ForAsset(asset string) ReviewableRequestQI
	// ForTypes - filters requests by request type
	ForTypes(requestTypes []xdr.ReviewableRequestType) ReviewableRequestQI
	// Page specifies the paging constraints for the query being built by `q`.
	Page(page db2.PageQuery) ReviewableRequestQI
	// Select loads the results of the query specified by `q`
	Select() ([]ReviewableRequest, error)
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
	}).Where("id = ?", request.ID)

	_, err := q.parent.Exec(query)
	return err
}

// ByID - selects request by id. Returns nil, nil if not found
func (q *ReviewableRequestQ) ByID(requestID uint64) (*ReviewableRequest, error) {
	if q.Err != nil {
		return nil, q.Err
	}

	query := selectReviewableRequest.Where("id = ?", requestID)

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

// ForRequestor - filters requests by requestor
func (q *ReviewableRequestQ) ForRequestor(requestor string) ReviewableRequestQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("requestor = ?", requestor)
	return q
}

// ForReviewer - filters requests by reviewer
func (q *ReviewableRequestQ) ForReviewer(reviewer string) ReviewableRequestQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("reviewer = ?", reviewer)
	return q
}

// ForState - filters requests by state
func (q *ReviewableRequestQ) ForState(state int64) ReviewableRequestQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("request_state = ?", state)
	return q
}

// ForType - filters requests by type
func (q *ReviewableRequestQ) ForType(requestType int64) ReviewableRequestQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("request_type = ?", requestType)
	return q
}

// ForAsset - filters requests by asset
func (q *ReviewableRequestQ) ForAsset(asset string) ReviewableRequestQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("details->>'asset' = ?", asset)
	return q
}

// ForTypes - filters requests by request type
func (q *ReviewableRequestQ) ForTypes(requestTypes []xdr.ReviewableRequestType) ReviewableRequestQI {
	if q.Err != nil {
		return q
	}

	if len(requestTypes) == 0 {
		return q
	}

	query, values := sqx.InForReviewableRequestTypes("request_type", requestTypes...)

	q.sql = q.sql.Where(query, values...)
	return q
}

// Page specifies the paging constraints for the query being built by `q`.
func (q *ReviewableRequestQ) Page(page db2.PageQuery) ReviewableRequestQI {
	if q.Err != nil {
		return q
	}

	q.sql, q.Err = page.ApplyTo(q.sql, "id")
	return q
}

// Select loads the results of the query specified by `q`
func (q *ReviewableRequestQ) Select() ([]ReviewableRequest, error) {
	if q.Err != nil {
		return nil, q.Err
	}

	var result []ReviewableRequest
	q.Err = q.parent.Select(&result, q.sql)
	return result, q.Err
}

var selectReviewableRequest = sq.Select("id", "requestor", "reviewer", "reference", "reject_reason", "request_type", "request_state", "hash",
	"details").From("reviewable_request")
