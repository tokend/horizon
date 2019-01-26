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
	Type                     xdr.ReviewableRequestType `json:"type"`
	AssetCreation            *AssetCreationRequest     `json:"asset_create,omitempty"`
	AssetUpdate              *AssetUpdateRequest       `json:"asset_update,omitempty"`
	PreIssuanceCreate        *PreIssuanceRequest       `json:"pre_issuance_create,omitempty"`
	IssuanceCreate           *IssuanceRequest          `json:"issuance_create,omitempty"`
	Withdraw                 *WithdrawalRequest        `json:"withdraw,omitempty"`
	TwoStepWithdraw          *WithdrawalRequest        `json:"two_step_withdrawal"`
	Sale                     *SaleRequest              `json:"sale,omitempty"`
	LimitsUpdate             *LimitsUpdateRequest      `json:"limits_update"`
	AmlAlert                 *AmlAlertRequest          `json:"aml_alert"`
	ChangeRole               *ChangeRoleRequest        `json:"change_role,omitempty"`
	UpdateSaleDetails        *UpdateSaleDetailsRequest `json:"update_sale_details"`
	UpdateSaleEndTimeRequest *UpdateSaleEndTimeRequest `json:"update_sale_end_time_request"`
	AtomicSwapBidCreation    *AtomicSwapBidCreation    `json:"atomic_swap_bid_creation"`
	AtomicSwap               *AtomicSwap               `json:"atomic_swap"`
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

//AssetCreationRequest - asset creation request details
type AssetCreationRequest struct {
	Asset                  string            `json:"asset"`
	Policies               int32             `json:"policies"`
	PreIssuedAssetSigner   string            `json:"pre_issued_asset_signer"`
	MaxIssuanceAmount      string            `json:"max_issuance_amount"`
	InitialPreissuedAmount string            `json:"initial_preissued_amount"`
	Details                regources.Details `json:"details"`
}

//AssetUpdateRequest - asset update request details
type AssetUpdateRequest struct {
	Asset    string            `json:"asset"`
	Policies int32             `json:"policies"`
	Details  regources.Details `json:"details"`
}

//PreIssuanceRequest - request details
type PreIssuanceRequest struct {
	Asset     string `json:"asset"`
	Amount    string `json:"amount"`
	Signature string `json:"signature"`
	Reference string `json:"reference"`
}

//IssuanceRequest - request details
type IssuanceRequest struct {
	Asset    string            `json:"asset"`
	Amount   string            `json:"amount"`
	Receiver string            `json:"receiver"`
	Details  regources.Details `json:"external_details"`
}

//WithdrawalRequest - request details
type WithdrawalRequest struct {
	BalanceID       string            `json:"balance_id"`
	Amount          string            `json:"amount"`
	FixedFee        string            `json:"fixed_fee"`
	PercentFee      string            `json:"percent_fee"`
	Details         regources.Details `json:"external_details"`
	ReviewerDetails regources.Details `json:"reviewer_details"`
}

//SaleRequest - request details
type SaleRequest struct {
	BaseAsset           string                     `json:"base_asset"`
	DefaultQuoteAsset   string                     `json:"quote_asset"`
	StartTime           time.Time                  `json:"start_time"`
	EndTime             time.Time                  `json:"end_time"`
	SoftCap             string                     `json:"soft_cap"`
	HardCap             string                     `json:"hard_cap"`
	Details             regources.Details          `json:"details"`
	QuoteAssets         []regources.SaleQuoteAsset `json:"quote_assets"`
	SaleType            xdr.SaleType               `json:"sale_type"`
	BaseAssetForHardCap string                     `json:"base_asset_for_hard_cap"`
}

//LimitsUpdateRequest - request details
type LimitsUpdateRequest struct {
	DocumentHash string            `json:"document_hash"`
	Details      regources.Details `json:"details"`
}

//AmlAlertRequest - request details
type AmlAlertRequest struct {
	BalanceID string `json:"balance_id"`
	Amount    string `json:"amount"`
	Reason    string `json:"reason"`
}

//ChangeRoleRequest - request details
type ChangeRoleRequest struct {
	DestinationAccount string            `json:"destination_account"`
	AccountRoleToSet   uint64            `json:"account_role_to_set"`
	KYCData            regources.Details `json:"kyc_data"`
	SequenceNumber     uint32            `json:"sequence_number"`
}

//UpdateSaleDetailsRequest - request details
type UpdateSaleDetailsRequest struct {
	SaleID     uint64            `json:"sale_id"`
	NewDetails regources.Details `json:"new_details"`
}

//UpdateSaleEndTimeRequest - request details
type UpdateSaleEndTimeRequest struct {
	SaleID     uint64    `json:"sale_id"`
	NewEndTime time.Time `json:"new_end_time"`
}

//AtomicSwapBidCreation - request details
type AtomicSwapBidCreation struct {
	BaseBalance string                 `json:"base_balance"`
	BaseAmount  uint64                 `json:"base_amount"`
	Details     regources.Details      `json:"details"`
	QuoteAssets []regources.AssetPrice `json:"quote_assets"`
}

//AtomicSwap - request details
type AtomicSwap struct {
	BidID      uint64 `json:"bid_id"`
	BaseAmount uint64 `json:"base_amount"`
	QuoteAsset string `json:"quote_asset"`
}
