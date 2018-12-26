package history2

/*
import (
	"time"

	sq "github.com/lann/squirrel"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2"
	"gitlab.com/tokend/horizon/db2/sqx"
)

type ReviewableRequestQ struct {
	repo *db2.Repo
}

// NewAccountsQ
func NewReviewableRequestQ(repo *db2.Repo) *ReviewableRequestQ {
	return &ReviewableRequestQ{
		repo: repo,
	}
}

// ByID - selects request by id. Returns nil, nil if not found
func (q *ReviewableRequestQ) ByID(requestID uint64) (*ReviewableRequest, error) {
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

func (q *ReviewableRequestQ) ByIDs(requestIDs []int64) ReviewableRequestQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where(sq.Eq{"id": requestIDs})
	return q
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

func (q *ReviewableRequestQ) ForCounterparty(counterparty string) ReviewableRequestQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("((reviewer = ?) or (requestor = ?))", counterparty, counterparty)
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

	q.sql = q.sql.Where("updated_at >= ?::timestamp", tmf)
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

func (q *ReviewableRequestQ) CountQuery() ReviewableRequestQI {
	if q.Err != nil {
		return q
	}
	q.sql = sq.Select("COUNT(*)").From("reviewable_request")
	return q
}

func (q *ReviewableRequestQ) Count() (int64, error) {
	if q.Err != nil {
		return 0, q.Err
	}

	var result int64
	q.Err = q.parent.Get(&result, q.sql)
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

// Withdraw
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

func (q *ReviewableRequestQ) ContractsByContractNumber(contractNumber string) ReviewableRequestQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("details->'contract'->'details'->>'contract_number' = ?", contractNumber)
	return q
}

func (q *ReviewableRequestQ) ASwapByBidID(bidID int64) ReviewableRequestQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("details->'atomic_swap'->>'bid_id' = ?", bidID)
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

var selectReviewableRequest = sq.Select("id", "requestor", "reviewer", "reference", "reject_reason",
	"request_type", "request_state", "hash", "details", "created_at", "updated_at", "all_tasks", "pending_tasks",
	"external_details").From("reviewable_request")
*/