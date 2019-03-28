package history

import (
	"gitlab.com/tokend/regources"
	"time"

	"database/sql/driver"

	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2"
)

type ReviewableRequestDetails struct {
	AssetCreation         *AssetCreationRequest     `json:"create_asset,omitempty"`
	AssetUpdate           *AssetUpdateRequest       `json:"update_asset,omitempty"`
	PreIssuanceCreate     *PreIssuanceRequest       `json:"create_pre_issuance,omitempty"`
	IssuanceCreate        *IssuanceRequest          `json:"create_issuance,omitempty"`
	Withdraw              *WithdrawalRequest        `json:"create_withdraw,omitempty"`
	TwoStepWithdraw       *WithdrawalRequest        `json:"two_step_withdrawal"`
	Sale                  *SaleRequest              `json:"create_sale,omitempty"`
	LimitsUpdate          *LimitsUpdateRequest      `json:"update_limits,omitempty"`
	AmlAlert              *AmlAlertRequest          `json:"create_aml_alert,omitempty"`
	ChangeRole            *ChangeRoleRequest        `json:"change_role,omitempty"`
	UpdateSaleDetails     *UpdateSaleDetailsRequest `json:"update_sale_details,omitempty"`
	PromotionUpdate       *PromotionUpdateRequest   `json:"promotion_update"`
	Invoice               *InvoiceRequest           `json:"invoice"`
	Contract              *ContractRequest          `json:"contract"`
	AtomicSwapBidCreation *AtomicSwapBidCreation    `json:"create_atomic_swap_bid,omitempty"`
	AtomicSwap            *AtomicSwap               `json:"create_atomic_swap,omitempty"`
}

func (r ReviewableRequestDetails) Value() (driver.Value, error) {
	result, err := db2.DriverValue(r)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal details")
	}

	return result, nil
}

func (r *ReviewableRequestDetails) Scan(src interface{}) error {
	err := db2.DriveScan(src, r)
	if err != nil {
		return errors.Wrap(err, "failed to scan details")
	}

	return nil
}

type AssetCreationRequest struct {
	Asset                  string                 `json:"asset"`
	Type                   uint64                 `json:"type"`
	Policies               int32                  `json:"policies"`
	PreIssuedAssetSigner   string                 `json:"pre_issued_asset_signer"`
	MaxIssuanceAmount      string                 `json:"max_issuance_amount"`
	InitialPreissuedAmount string                 `json:"initial_preissued_amount"`
	Details                map[string]interface{} `json:"details"`
}

type AssetUpdateRequest struct {
	Asset    string                 `json:"asset"`
	Policies int32                  `json:"policies"`
	Details  map[string]interface{} `json:"details"`
}

type PreIssuanceRequest struct {
	Asset     string `json:"asset"`
	Amount    string `json:"amount"`
	Signature string `json:"signature"`
	Reference string `json:"reference"`
}

type IssuanceRequest struct {
	Asset           string                 `json:"asset"`
	Amount          string                 `json:"amount"`
	Receiver        string                 `json:"receiver"`
	ExternalDetails map[string]interface{} `json:"external_details"`
}

type WithdrawalRequest struct {
	BalanceID              string                 `json:"balance_id"`
	Amount                 string                 `json:"amount"`
	FixedFee               string                 `json:"fixed_fee"`
	PercentFee             string                 `json:"percent_fee"`
	ExternalDetails        map[string]interface{} `json:"external_details"`
	ReviewerDetails        map[string]interface{} `json:"reviewer_details"`
	PreConfirmationDetails map[string]interface{} `json:"pre_confirmation_details"`
}

type SaleRequest struct {
	BaseAsset           string                     `json:"base_asset"`
	DefaultQuoteAsset   string                     `json:"quote_asset"`
	StartTime           time.Time                  `json:"start_time"`
	EndTime             time.Time                  `json:"end_time"`
	SoftCap             string                     `json:"soft_cap"`
	HardCap             string                     `json:"hard_cap"`
	Details             map[string]interface{}     `json:"details"`
	QuoteAssets         []regources.SaleQuoteAsset `json:"quote_assets"`
	SaleType            xdr.SaleType               `json:"sale_type"`
	BaseAssetForHardCap string                     `json:"base_asset_for_hard_cap"`
}

type LimitsUpdateRequest struct {
	DocumentHash string                 `json:"document_hash"`
	Details      map[string]interface{} `json:"details"`
}

type AmlAlertRequest struct {
	BalanceID string `json:"balance_id"`
	Amount    string `json:"amount"`
	Reason    string `json:"reason"`
}

type ChangeRoleRequest struct {
	DestinationAccount string                 `json:"destination_account"`
	AccountRoleToSet   uint64                 `json:"account_role_to_set"`
	KYCData            map[string]interface{} `json:"kyc_data"`
	SequenceNumber     uint32                 `json:"sequence_number"`
}

type UpdateSaleDetailsRequest struct {
	SaleID     uint64                 `json:"sale_id"`
	NewDetails map[string]interface{} `json:"new_details"`
}

type InvoiceRequest struct {
	Asset           string                 `json:"receiver_balance_id"`
	Amount          uint64                 `json:"amount"`
	ContractID      *int64                 `json:"contract_id"`
	Details         map[string]interface{} `json:"details"`
	PayerBalance    string                 `json:"payer_balance"`
	ReceiverBalance string                 `json:"receiver_balance"`
}

type ContractRequest struct {
	Escrow    string                 `json:"escrow"`
	Details   map[string]interface{} `json:"details"`
	StartTime time.Time              `json:"start_time"`
	EndTime   time.Time              `json:"end_time"`
}

type PromotionUpdateRequest struct {
	SaleID           uint64      `json:"sale_id"`
	NewPromotionData SaleRequest `json:"new_promotion_data"`
}

type AtomicSwapBidCreation struct {
	BaseBalance string                 `json:"base_balance"`
	BaseAmount  uint64                 `json:"base_amount"`
	Details     map[string]interface{} `json:"details"`
	QuoteAssets []regources.AssetPrice `json:"quote_assets"`
}

type AtomicSwap struct {
	BidID      uint64 `json:"bid_id"`
	BaseAmount uint64 `json:"base_amount"`
	QuoteAsset string `json:"quote_asset"`
}
