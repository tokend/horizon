package history2

import (
	"database/sql/driver"

	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2"
	"gitlab.com/tokend/regources"
)

// OperationDetails - stores details of the operation performed in union switch form.
// Only one value must be selected at a type
type OperationDetails struct {
	Type                       xdr.OperationType                  `json:"type"`
	CreateAccount              *CreateAccountDetails              `json:"create_account,omitempty"`
	ManageAccount              *ManageAccountDetails              `json:"manage_account,omitempty"`
	ManageBalance              *ManageBalanceDetails              `json:"manage_balance,omitempty"`
	ManageKeyValue             *ManageKeyValueDetails             `json:"manage_key_value,omitempty"`
	ManageAsset                *ManageAssetDetails                `json:"manage_asset,omitempty"`
	ManageAssetPair            *ManageAssetPairDetails            `json:"manage_asset_pair,omitempty"`
	ManageLimits               *ManageLimitsDetails               `json:"manage_limits,omitempty"`
	ManageOffer                *ManageOfferDetails                `json:"manage_offer,omitempty"`
	ManageExternalSystemPool   *ManageExternalSystemPoolDetails   `json:"manage_external_system_pool"`
	BindExternalSystemAccount  *BindExternalSystemAccountDetails  `json:"bind_external_system_account"`
	Payment                    *PaymentDetails                    `json:"payment,omitempty"`
	Payout                     *PayoutDetails                     `json:"payout,omitempty"`
	SetFee                     *SetFeeDetails                     `json:"set_fee,omitempty"`
	CancelAtomicSwapBid        *CancelAtomicSwapBidDetails        `json:"cancel_atomic_swap_bid,omitempty"`
	CheckSaleState             *CheckSaleStateDetails             `json:"check_sale_state,omitempty"`
	CreateChangeRoleRequest    *CreateChangeRoleRequestDetails    `json:"create_change_role_request"`
	CreatePreIssuanceRequest   *CreatePreIssuanceRequestDetails   `json:"create_pre_issuance_request,omitempty"`
	CreateIssuanceRequest      *CreateIssuanceRequestDetails      `json:"create_issuance_request,omitempty"`
	CreateManageLimitsRequest  *CreateManageLimitsRequestDetails  `json:"create_manage_limits_request,omitempty"`
	CreateWithdrawRequest      *CreateWithdrawRequestDetails      `json:"create_withdraw_request,omitempty"`
	CreateAMLAlertRequest      *CreateAMLAlertRequestDetails      `json:"create_aml_alert_request,omitempty"`
	CreateSaleRequest          *CreateSaleRequestDetails          `json:"create_sale_request,omitempty"`
	CreateAtomicSwapBidRequest *CreateAtomicSwapBidRequestDetails `json:"create_atomic_swap_bid_request,omitempty"`
	CreateAtomicSwapRequest    *CreateAtomicSwapRequestDetails    `json:"create_atomic_swap_request,omitempty"`
	ReviewRequest              *ReviewRequestDetails              `json:"review_request,omitempty"`
	ManageSale                 *ManageSaleDetails                 `json:"manage_sale,omitempty"`
}

func (r OperationDetails) Value() (driver.Value, error) {
	result, err := db2.DriverValue(r)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal operation details")
	}

	return result, nil
}

func (r *OperationDetails) Scan(src interface{}) error {
	err := db2.DriveScan(src, r)
	if err != nil {
		return errors.Wrap(err, "failed to scan operation details")
	}

	return nil
}

// CreateAccountDetails - stores details of create account operation
type CreateAccountDetails struct {
	AccountAddress string          `json:"account_address"`
	AccountType    xdr.AccountType `json:"account_type"`
}

//ManageBalanceDetails - details of ManageBalanceOp
type ManageBalanceDetails struct {
	DestinationAccount string                  `json:"destination_account"`
	Action             xdr.ManageBalanceAction `json:"action"`
	Asset              xdr.AssetCode           `json:"asset"`
	BalanceAddress     string                  `json:"balance_address"`
}

//ManageAccountDetails - details of ManageAccountOp
type ManageAccountDetails struct {
	AccountAddress       string `json:"account_address"`
	BlockReasonsToAdd    int32  `json:"block_reasons_to_add"`
	BlockReasonsToRemove int32  `json:"block_reasons_to_remove"`
}

//ManageKeyValueDetails - details of ManageKeyValueOp
type ManageKeyValueDetails struct {
	Key    string                  `json:"key"`
	Action xdr.ManageKvAction      `json:"action"`
	Value  *xdr.KeyValueEntryValue `json:"value,omitempty"`
}

//SetFeeDetails - details of SetFeeOp
type SetFeeDetails struct {
	AssetCode      xdr.AssetCode    `json:"asset_code"`
	FixedFee       string           `json:"fixed_fee"`
	PercentFee     string           `json:"percent_fee"`
	FeeType        xdr.FeeType      `json:"fee_type"`
	AccountAddress *string          `json:"account_address,omitempty"`
	AccountType    *xdr.AccountType `json:"account_type,omitempty"`
	Subtype        int64            `json:"subtype"`
	LowerBound     string           `json:"lower_bound"`
	UpperBound     string           `json:"upper_bound"`
	// FeeAsset deprecated
}

//CreateWithdrawRequestDetails - details of corresponding op
type CreateWithdrawRequestDetails struct {
	BalanceAddress  string                 `json:"balance_address"`
	Amount          string                 `json:"amount"`
	FixedFee        string                 `json:"fixed_fee"`
	PercentFee      string                 `json:"percent_fee"`
	ExternalDetails map[string]interface{} `json:"external_details"`
}

//CreateManageLimitsRequestDetails - details of corresponding op
type CreateManageLimitsRequestDetails struct {
	Data      map[string]interface{} `json:"data"`
	RequestID int64                  `json:"request_id"`
}

//ManageLimitsDetails - details of corresponding op
type ManageLimitsDetails struct {
	Action   xdr.ManageLimitsAction       `json:"action"`
	Creation *ManageLimitsCreationDetails `json:"creation_details,omitempty"`
	Removal  *ManageLimitsRemovalDetails  `json:"removal_details,omitempty"`
}

//ManageLimitsCreationDetails - details of corresponding op
type ManageLimitsCreationDetails struct {
	AccountType     *xdr.AccountType `json:"account_type,omitempty"`
	AccountAddress  string           `json:"account_address,omitempty"`
	StatsOpType     xdr.StatsOpType  `json:"stats_op_type"`
	AssetCode       xdr.AssetCode    `json:"asset_code"`
	IsConvertNeeded bool             `json:"is_convert_needed"`
	DailyOut        string           `json:"daily_out"`
	WeeklyOut       string           `json:"weekly_out"`
	MonthlyOut      string           `json:"monthly_out"`
	AnnualOut       string           `json:"annual_out"`
}

//ManageLimitsRemovalDetails - details of corresponding op
type ManageLimitsRemovalDetails struct {
	ID int64 `json:"id"`
}

//ManageAssetPairDetails - details of corresponding op
type ManageAssetPairDetails struct {
	BaseAsset               xdr.AssetCode `json:"base_asset"`
	QuoteAsset              xdr.AssetCode `json:"quote_asset"`
	PhysicalPrice           string        `json:"physical_price"`
	PhysicalPriceCorrection string        `json:"physical_price_correction"`
	MaxPriceStep            string        `json:"max_price_step"`
	PoliciesI               int32         `json:"policies_i"`
}

//ManageOfferDetails - details of corresponding op
type ManageOfferDetails struct {
	OfferID     int64         `json:"offer_id,omitempty"`
	OrderBookID int64         `json:"order_book_id"`
	BaseAsset   xdr.AssetCode `json:"base_asset"`
	QuoteAsset  xdr.AssetCode `json:"quote_asset"`
	Amount      string        `json:"base_amount"`
	Price       string        `json:"price"`
	IsBuy       bool          `json:"is_buy"`
	Fee         string        `json:"fee"`
	IsDeleted   bool          `json:"is_deleted"`
}

//ReviewRequestDetails - details of corresponding op
type ReviewRequestDetails struct {
	Action          xdr.ReviewRequestOpAction         `json:"action"`
	Reason          string                            `json:"reason"`
	RequestHash     string                            `json:"request_hash"`
	RequestID       int64                             `json:"request_id"`
	RequestType     xdr.ReviewableRequestType         `json:"request_type"`
	RequestDetails  xdr.ReviewRequestOpRequestDetails `json:"request_details"`
	IsFulfilled     bool                              `json:"is_fulfilled"`
	AddedTasks      uint32                            `json:"added_tasks"`
	RemovedTasks    uint32                            `json:"removed_tasks"`
	ExternalDetails string                            `json:"external_details"`
}

//ManageAssetDetails - details of corresponding op
type ManageAssetDetails struct {
	AssetCode         xdr.AssetCode          `json:"asset_code,omitempty"`
	RequestID         int64                  `json:"request_id"`
	Action            xdr.ManageAssetAction  `json:"action"`
	Policies          *int32                 `json:"policies,omitempty"`
	Details           map[string]interface{} `json:"details,omitempty"`
	PreissuedSigner   string                 `json:"preissued_signer,omitempty"`
	MaxIssuanceAmount string                 `json:"max_issuance_amount,omitempty"`
}

//CreatePreIssuanceRequestDetails - details of corresponding op
type CreatePreIssuanceRequestDetails struct {
	AssetCode   xdr.AssetCode `json:"asset_code"`
	Amount      string        `json:"amount"`
	RequestID   int64         `json:"request_id"`
	IsFulfilled bool          `json:"is_fulfilled"`
}

//CreateIssuanceRequestDetails - details of corresponding op
type CreateIssuanceRequestDetails struct {
	FixedFee               string                 `json:"fixed_fee"`
	PercentFee             string                 `json:"percent_fee"`
	Reference              string                 `json:"reference"`
	Amount                 string                 `json:"amount"`
	Asset                  xdr.AssetCode          `json:"asset"`
	ReceiverAccountAddress string                 `json:"receiver_account_address"`
	ReceiverBalanceAddress string                 `json:"receiver_balance_address"`
	ExternalDetails        map[string]interface{} `json:"external_details"`
	AllTasks               *int64                 `json:"all_tasks,omitempty"`
	RequestDetails         RequestDetails         `json:"request_details"`
}

//CreateSaleRequestDetails -details of corresponding op
type CreateSaleRequestDetails struct {
	RequestID         int64                      `json:"request_id"`
	BaseAsset         xdr.AssetCode              `json:"base_asset"`
	DefaultQuoteAsset xdr.AssetCode              `json:"default_quote_asset"`
	StartTime         int64                      `json:"start_time"`
	EndTime           int64                      `json:"end_time"`
	SoftCap           string                     `json:"soft_cap"`
	HardCap           string                     `json:"hard_cap"`
	QuoteAssets       []regources.SaleQuoteAsset `json:"quote_assets"`
	Details           map[string]interface{}     `json:"details"`
}

//CheckSaleStateDetails - details of corresponding op
type CheckSaleStateDetails struct {
	SaleID int64 `json:"sale_id"`
	Effect xdr.CheckSaleStateEffect
}

//CreateAtomicSwapBidRequestDetails - details of corresponding op
type CreateAtomicSwapBidRequestDetails struct {
	Amount         string                     `json:"amount"`
	BaseBalance    string                     `json:"base_balance"`
	QuoteAssets    []regources.SaleQuoteAsset `json:"quote_assets"`
	Details        map[string]interface{}     `json:"details"`
	RequestDetails RequestDetails             `json:"request_details"`
}

//CreateAtomicSwapRequestDetails - details of corresponding op
type CreateAtomicSwapRequestDetails struct {
	BidID          int64          `json:"bid_id"`
	BaseAmount     string         `json:"base_amount"`
	QuoteAsset     xdr.AssetCode  `json:"quote_asset"`
	RequestDetails RequestDetails `json:"request_details"`
}

//CancelAtomicSwapBidDetails - details of corresponding op
type CancelAtomicSwapBidDetails struct {
	BidID int64 `json:"bid_id"`
}

//CreateAMLAlertRequestDetails - details of corresponding op
type CreateAMLAlertRequestDetails struct {
	Amount         string `json:"amount"`
	BalanceAddress string `json:"balance_address"`
	Reason         string `json:"reason"`
}

// PaymentDetails - stores details of payment operation
type PaymentDetails struct {
	AccountFrom             string        `json:"account_from"`
	AccountTo               string        `json:"account_to"`
	BalanceFrom             string        `json:"balance_from"`
	BalanceTo               string        `json:"balance_to"`
	Amount                  string        `json:"amount"`
	Asset                   xdr.AssetCode `json:"asset"`
	SourceFeeData           FeeData       `json:"source_fee_data"`
	DestinationFeeData      FeeData       `json:"destination_fee_data"`
	SourcePayForDestination bool          `json:"source_pay_for_destination"`
	Subject                 string        `json:"subject"`
	Reference               string        `json:"reference"`
	UniversalAmount         string        `json:"universal_amount"`
}

//PayoutDetails - details of corresponding op
type PayoutDetails struct {
	SourceAccountAddress string        `json:"source_account_address"`
	SourceBalanceAddress string        `json:"source_balance_address"`
	Asset                xdr.AssetCode `json:"asset"`
	MaxPayoutAmount      string        `json:"max_payout_amount"`
	MinAssetHolderAmount string        `json:"min_asset_holder_amount"`
	MinPayoutAmount      string        `json:"min_payout_amount"`
	ExpectedFixedFee     string        `json:"expected_fixed_fee"`
	ExpectedPercentFee   string        `json:"expected_percent_fee"`
	ActualFixedFee       string        `json:"actual_fixed_fee"`
	ActualPercentFee     string        `json:"actual_percent_fee"`
	ActualPayoutAmount   string        `json:"actual_payout_amount"`
}

//FeeData - details of fees for the op
type FeeData struct {
	FixedFee  string `json:"fixed_fee"`
	ActualFee string `json:"actual_fee"`
}

//RequestDetails - details of the request created or reviewed via op
type RequestDetails struct {
	RequestID   int64 `json:"request_id,omitempty"`
	IsFulfilled bool  `json:"is_fulfilled"`
}

//CreateChangeRoleRequestDetails - details of corresponding op
type CreateChangeRoleRequestDetails struct {
	DestinationAccount string                 `json:"destination_account"`
	AccountRoleToSet   uint64                 `json:"account_role_to_set"`
	KYCData            map[string]interface{} `json:"kyc_data"`
	AllTasks           *uint32                `json:"all_tasks"`
	RequestDetails     RequestDetails         `json:"request_details"`
}

//ManageExternalSystemPoolDetails - details of corresponding op
type ManageExternalSystemPoolDetails struct {
	Action xdr.ManageExternalSystemAccountIdPoolEntryAction `json:"action"`
	Create *CreateExternalSystemPoolDetails                 `json:"create"`
	Remove *RemoveExternalSystemPoolDetails                 `json:"remove"`
}

//CreateExternalSystemPoolDetails - details of corresponding op
type CreateExternalSystemPoolDetails struct {
	PoolID             uint64 `json:"pool_id"`
	Data               string `json:"data"`
	Parent             uint64 `json:"parent"`
	ExternalSystemType int32  `json:"external_system_type"`
}

//RemoveExternalSystemPoolDetails - details of corresponding op
type RemoveExternalSystemPoolDetails struct {
	PoolID uint64 `json:"pool_id"`
}

//BindExternalSystemAccountDetails - details of corresponding op
type BindExternalSystemAccountDetails struct {
	ExternalSystemType int32 `json:"external_system_type"`
}

//ManageSaleDetails - details of corresponding op
type ManageSaleDetails struct {
	SaleID uint64               `json:"sale_id"`
	Action xdr.ManageSaleAction `json:"action"`
}
