package history

import (
	"time"

	"database/sql/driver"

	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/horizon/db2"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/regources"
)

type ReviewableRequestDetails struct {
	AssetCreation            *AssetCreationRequest     `json:"asset_create,omitempty"`
	AssetUpdate              *AssetUpdateRequest       `json:"asset_update,omitempty"`
	PreIssuanceCreate        *PreIssuanceRequest       `json:"pre_issuance_create,omitempty"`
	IssuanceCreate           *IssuanceRequest          `json:"issuance_create,omitempty"`
	Withdraw                 *WithdrawalRequest        `json:"withdraw,omitempty"`
	TwoStepWithdraw          *WithdrawalRequest        `json:"two_step_withdrawal"`
	Sale                     *SaleRequest              `json:"sale,omitempty"`
	LimitsUpdate             *LimitsUpdateRequest      `json:"limits_update"`
	AmlAlert                 *AmlAlertRequest          `json:"aml_alert"`
	UpdateKYC                *UpdateKYCRequest         `json:"update_kyc,omitempty"`
	UpdateSaleDetails        *UpdateSaleDetailsRequest `json:"update_sale_details"`
	UpdateSaleEndTimeRequest *UpdateSaleEndTimeRequest `json:"update_sale_end_time_request"`
	PromotionUpdate          *PromotionUpdateRequest   `json:"promotion_update"`
	Invoice                  *InvoiceRequest           `json:"invoice"`
	Contract                 *ContractRequest          `json:"contract"`
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
	DestAssetCode          string                 `json:"dest_asset_code"`
	DestAssetAmount        string                 `json:"dest_asset_amount"`
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
	State               xdr.SaleState              `json:"state"`
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

type UpdateKYCRequest struct {
	AccountToUpdateKYC string                   `json:"updated_account_id"`
	AccountTypeToSet   xdr.AccountType          `json:"account_type_to_set"`
	KYCLevel           uint32                   `json:"kyc_level"`
	KYCData            map[string]interface{}   `json:"kyc_data"`
	AllTasks           uint32                   `json:"all_tasks"`
	PendingTasks       uint32                   `json:"pending_tasks"`
	SequenceNumber     uint32                   `json:"sequence_number"`
	ExternalDetails    []map[string]interface{} `json:"external_details"`
}

type UpdateSaleDetailsRequest struct {
	SaleID     uint64                 `json:"sale_id"`
	NewDetails map[string]interface{} `json:"new_details"`
}

type InvoiceRequest struct {
	Asset      string                 `json:"receiver_balance_id"`
	Amount     uint64                 `json:"amount"`
	ContractID *int64                 `json:"contract_id"`
	Details    map[string]interface{} `json:"details"`
}

type ContractRequest struct {
	Escrow    string                 `json:"escrow"`
	Details   map[string]interface{} `json:"details"`
	StartTime time.Time              `json:"start_time"`
	EndTime   time.Time              `json:"end_time"`
}

type UpdateSaleEndTimeRequest struct {
	SaleID     uint64    `json:"sale_id"`
	NewEndTime time.Time `json:"new_end_time"`
}

type PromotionUpdateRequest struct {
	SaleID           uint64      `json:"sale_id"`
	NewPromotionData SaleRequest `json:"new_promotion_data"`
}
