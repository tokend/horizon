package history

import (
	"fmt"

	sq "github.com/lann/squirrel"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/swarmfund/horizon/db2"
	"gitlab.com/swarmfund/horizon/db2/sqx"
	"time"
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
	// ForTypes - filters requests by request type
	ForTypes(requestTypes []xdr.ReviewableRequestType) ReviewableRequestQI
	// Page specifies the paging constraints for the query being built by `q`.
	Page(page db2.PageQuery) ReviewableRequestQI
	// UpdatedAfter - selects requests updated after given timestamp
	UpdatedAfter(timestamp int64) ReviewableRequestQI
	// Select loads the results of the query specified by `q`
	Select() ([]ReviewableRequest, error)

	// Request Type specific filters. Filter for request type must be applied separately

	// Asset Management
	// AssetManagementByAsset - filters asset management requests by asset
	AssetManagementByAsset(assetCode string) ReviewableRequestQI

	// PreIssuance
	// PreIssuanceByAsset - filters pre issuance requests by asset
	PreIssuanceByAsset(assetCode string) ReviewableRequestQI

	// Issuance
	// IssuanceByAsset - filters issuance requests by asset
	IssuanceByAsset(assetCode string) ReviewableRequestQI

	// Withdrawal
	// WithdrawalByDestAsset - filters withdrawal requests by dest asset
	WithdrawalByDestAsset(assetCode string) ReviewableRequestQI

	// Sales
	// SalesByBaseAsset - filters sale requests by base asset
	SalesByBaseAsset(assetCode string) ReviewableRequestQI

	// Limits
	// LimitsByDocHash - filters limits request by document hash
	LimitsByDocHash(hash string) ReviewableRequestQI

	// Invoices
	// InvoiceByPayer - filters invoice requests by sender
	InvoiceByPayer(senderAccountID string) ReviewableRequestQI

	// KYC
	// KYCByAccountToUpdateKYC - filters update KYC requests by accountID of the owner of KYC
	KYCByAccountToUpdateKYC(accountID string) ReviewableRequestQI
	// KYCByMaskSet - filters update KYC requests by mask which must be set. If maskSetPartialEq is false, request will be returned
	// even if only part of the mask is set
	KYCByMaskSet(mask int64, maskSetPartialEq bool) ReviewableRequestQI
	// KYCByMaskNotSet - filters update KYC requests by mask which must not be set
	KYCByMaskNotSet(mask int64) ReviewableRequestQI
	// KYCByAccountTypeToSet - filters update KYC requests by account type which must be set.
	KYCByAccountTypeToSet(accountTypeToSet xdr.AccountType) ReviewableRequestQI
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

	query := sq.Insert("reviewable_request").SetMap(map[string]interface{}{
		"id":            request.ID,
		"requestor":     request.Requestor,
		"reviewer":      request.Reviewer,
		"reference":     request.Reference,
		"reject_reason": request.RejectReason,
		"request_type":  request.RequestType,
		"request_state": request.RequestState,
		"hash":          request.Hash,
		"details":       request.Details,
		"created_at":    request.CreatedAt,
		"updated_at":    request.UpdatedAt,
	})

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
		"updated_at":    request.UpdatedAt,
	}).Where("id = ?", request.ID)

	_, err := q.parent.Exec(query)
	return err
}

// ByID - selects request by id. Returns nil, nil if not found
func (q *ReviewableRequestQ) ByID(requestID uint64) (*ReviewableRequest, error) {
	if q.Err != nil {
		return nil, q.Err
	}

	query := q.sql.Where("id = ?", requestID)

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

// ByDetailsEq - filters by specified key value from the details. Note: do not pass key passed by user
func (q *ReviewableRequestQ) ByDetailsEq(key, value string) ReviewableRequestQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where(fmt.Sprintf("details->>'%s' = ?", key), value)
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

// UpdatedAfter - selects requests updated after given timestamp
func (q *ReviewableRequestQ) UpdatedAfter(timestamp int64) ReviewableRequestQI {
	if q.Err != nil {
		return q
	}

	tm := time.Unix(timestamp, 0)
	tmf := tm.Format(time.RFC3339)

	q.sql = q.sql.Where(fmt.Sprintf("updated_at >= '%s'::timestamp", tmf))
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

func (q *ReviewableRequestQ) AssetManagementByAsset(assetCode string) ReviewableRequestQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("(details->'asset_create'->>'asset' = ? OR details->'asset_update'->>'asset' = ?)", assetCode, assetCode)
	return q
}


// PreIssuance
// PreIssuanceByAsset - filters pre issuance requests by asset
func (q *ReviewableRequestQ) PreIssuanceByAsset(assetCode string) ReviewableRequestQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("details->'pre_issuance_create'->>'asset' = ?", assetCode)
	return q
}

// Issuance
// IssuanceByAsset - filters issuance requests by asset
func (q *ReviewableRequestQ) IssuanceByAsset(assetCode string) ReviewableRequestQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("details->'issuance_create'->>'asset' = ?", assetCode)
	return q
}

// Withdrawal
// WithdrawalByDestAsset - filters withdrawal requests by dest asset
func (q *ReviewableRequestQ) WithdrawalByDestAsset(assetCode string) ReviewableRequestQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("(details->'withdraw'->>'dest_asset_code' = ? OR details->'two_step_withdrawal'->>'dest_asset_code' = ?)", assetCode, assetCode)
	return q
}

// Sales
// SalesByBaseAsset - filters sale requests by base asset
func (q *ReviewableRequestQ) SalesByBaseAsset(assetCode string) ReviewableRequestQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("details->'sale'->>'base_asset' = ?", assetCode)
	return q
}

// Limits
// LimitsByDocHash - filters limits request by document hash
func (q *ReviewableRequestQ) LimitsByDocHash(hash string) ReviewableRequestQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("details->'limits_update'->>'document_hash' = ?", hash)
	return q
}

func (q *ReviewableRequestQ) InvoiceByPayer(senderAccountID string) ReviewableRequestQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("details->'invoice'->>'sender_account_id' = ?", senderAccountID)
	return q
}

// KYC
// KYCByAccountToUpdateKYC - filters update KYC requests by accountID of the owner of KYC
func (q *ReviewableRequestQ) KYCByAccountToUpdateKYC(accountID string) ReviewableRequestQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("details->'update_kyc'->>'updated_account_id' = ?", accountID)
	return q
}
// KYCByMaskSet - filters update KYC requests by mask which must be set. If mustBeEq is false, request will be returned
// even if only part of the mask is set
func (q *ReviewableRequestQ) KYCByMaskSet(mask int64, maskSetPartialEq bool) ReviewableRequestQI {
	if q.Err != nil {
		return q
	}

	if maskSetPartialEq {
		q.sql = q.sql.Where("(details->'update_kyc'->>'pending_tasks')::integer & ? <> 0", mask)
	} else {
		q.sql = q.sql.Where("(details->'update_kyc'->>'pending_tasks')::integer & ? = ?", mask, mask)
	}
	return q
}

// KYCByMaskNotSet - filters update KYC requests by mask which must not be set
func (q *ReviewableRequestQ) KYCByMaskNotSet(mask int64) ReviewableRequestQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("~(details->'update_kyc'->>'pending_tasks')::integer & ? = ?", mask, mask)
	return q
}

// KYCByAccountTypeToSet - filters update KYC requests by account type which must be set.
func (q *ReviewableRequestQ) KYCByAccountTypeToSet(accountTypeToSet xdr.AccountType) ReviewableRequestQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("(details->'update_kyc'->'account_type_to_set'->>'int')::integer = ?", int32(accountTypeToSet))
	return q
}

var selectReviewableRequest = sq.Select("id", "requestor", "reviewer", "reference", "reject_reason", "request_type", "request_state", "hash",
	"details", "created_at", "updated_at").From("reviewable_request")
