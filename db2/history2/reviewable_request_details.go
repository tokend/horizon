package history2

import (
	"time"

	"database/sql/driver"

	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2"
	"gitlab.com/tokend/regources/v2"
)

//ReviewableRequestDetails - stores in union switch details of the reviewable requests
type ReviewableRequestDetails struct {
	Type                xdr.ReviewableRequestType   `json:"type"`
	CreateAsset         *CreateAssetRequest         `json:"create_asset,omitempty"`
	UpdateAsset         *UpdateAssetRequest         `json:"update_asset,omitempty"`
	CreatePreIssuance   *CreatePreIssuanceRequest   `json:"create_pre_issuance,omitempty"`
	CreateIssuance      *CreateIssuanceRequest      `json:"create_issuance,omitempty"`
	CreateWithdraw      *CreateWithdrawalRequest    `json:"create_withdraw,omitempty"`
	CreateSale          *CreateSaleRequest          `json:"create_sale,omitempty"`
	UpdateLimits        *UpdateLimitsRequest        `json:"update_limits"`
	CreateAmlAlert      *CreateAmlAlertRequest      `json:"create_aml_alert"`
	ChangeRole          *ChangeRoleRequest          `json:"change_role,omitempty"`
	UpdateSaleDetails   *UpdateSaleDetailsRequest   `json:"update_sale_details"`
	CreateAtomicSwapBid *CreateAtomicSwapBidRequest `json:"create_atomic_swap_bid"`
	CreateAtomicSwap    *CreateAtomicSwapRequest    `json:"create_atomic_swap"`
}

//Value - implements db driver method for auto marshal
func (r ReviewableRequestDetails) Value() (driver.Value, error) {
	result, err := db2.DriverValue(r)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal details")
	}

	return result, nil
}

//Scan - implements db driver method for auto unmarshal
func (r *ReviewableRequestDetails) Scan(src interface{}) error {
	err := db2.DriveScan(src, r)
	if err != nil {
		return errors.Wrap(err, "failed to scan details")
	}

	return nil
}

//CreateAssetRequest - asset creation request details
type CreateAssetRequest struct {
	Asset                  string            `json:"asset"`
	Policies               int32             `json:"policies"`
	PreIssuedAssetSigner   string            `json:"pre_issued_asset_signer"`
	MaxIssuanceAmount      string            `json:"max_issuance_amount"`
	InitialPreissuedAmount string            `json:"initial_preissued_amount"`
	CreatorDetails         regources.Details `json:"details"`
}

//UpdateAssetRequest - asset update request details
type UpdateAssetRequest struct {
	Asset          string            `json:"asset"`
	Policies       int32             `json:"policies"`
	CreatorDetails regources.Details `json:"details"`
}

//CreatePreIssuanceRequest - request details
type CreatePreIssuanceRequest struct {
	Asset          string            `json:"asset"`
	Amount         string            `json:"amount"`
	Signature      string            `json:"signature"`
	Reference      string            `json:"reference"`
	CreatorDetails regources.Details `json:"creator_details"`
}

//CreateIssuanceRequest - request details
type CreateIssuanceRequest struct {
	Asset          string            `json:"asset"`
	Amount         string            `json:"amount"`
	Receiver       string            `json:"receiver"`
	CreatorDetails regources.Details `json:"external_details"`
}

//CreateWithdrawalRequest - request details
type CreateWithdrawalRequest struct {
	BalanceID       string            `json:"balance_id"`
	Amount          string            `json:"amount"`
	FixedFee        string            `json:"fixed_fee"`
	PercentFee      string            `json:"percent_fee"`
	CreatorDetails  regources.Details `json:"creator_details"`
	ReviewerDetails regources.Details `json:"reviewer_details"`
}

//CreateSaleRequest - request details
type CreateSaleRequest struct {
	BaseAsset           string                 `json:"base_asset"`
	DefaultQuoteAsset   string                 `json:"quote_asset"`
	StartTime           time.Time              `json:"start_time"`
	EndTime             time.Time              `json:"end_time"`
	SoftCap             string                 `json:"soft_cap"`
	HardCap             string                 `json:"hard_cap"`
	CreatorDetails      regources.Details      `json:"creator_details"`
	QuoteAssets         []regources.AssetPrice `json:"quote_assets"`
	SaleType            xdr.SaleType           `json:"sale_type"`
	BaseAssetForHardCap string                 `json:"base_asset_for_hard_cap"`
}

//UpdateLimitsRequest - request details
type UpdateLimitsRequest struct {
	CreatorDetails regources.Details `json:"creator_details"`
}

//CreateAmlAlertRequest - request details
type CreateAmlAlertRequest struct {
	BalanceID      string `json:"balance_id"`
	Amount         string `json:"amount"`
	CreatorDetails string `json:"creator_details"`
}

//ChangeRoleRequest - request details
type ChangeRoleRequest struct {
	DestinationAccount string            `json:"destination_account"`
	AccountRoleToSet   uint64            `json:"account_role_to_set"`
	KYCData            regources.Details `json:"kyc_data"`
	SequenceNumber     uint32            `json:"sequence_number"`
	CreatorDetails     regources.Details `json:"creator_details"`
}

//UpdateSaleDetailsRequest - request details
type UpdateSaleDetailsRequest struct {
	SaleID         uint64            `json:"sale_id"`
	CreatorDetails regources.Details `json:"creator_details"`
}

//CreateAtomicSwapBidRequest - request details
type CreateAtomicSwapBidRequest struct {
	BaseBalance    string                 `json:"base_balance"`
	BaseAmount     uint64                 `json:"base_amount"`
	CreatorDetails regources.Details      `json:"creator_details"`
	QuoteAssets    []regources.AssetPrice `json:"quote_assets"`
}

//CreateAtomicSwapRequest - request details
type CreateAtomicSwapRequest struct {
	BidID          uint64            `json:"bid_id"`
	BaseAmount     uint64            `json:"base_amount"`
	QuoteAsset     string            `json:"quote_asset"`
	CreatorDetails regources.Details `json:"creator_details"`
}
