package history2

import (
	"database/sql/driver"

	"time"

	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2"
	"gitlab.com/tokend/regources/rgenerated"
)

// OperationDetails - stores details of the operation performed in union switch form.
// Only one value must be selected at a type
type OperationDetails struct {
	//NOTE: omitempty MUST be specified for each switch value
	Type                       xdr.OperationType                  `json:"type"`
	CreateAccount              *CreateAccountDetails              `json:"create_account,omitempty"`
	ManageAccountRule          *ManageAccountRuleDetails          `json:"manage_account_rule,omitempty"`
	ManageAccountRole          *ManageAccountRoleDetails          `json:"manage_account_role,omitempty"`
	ManageSigner               *ManageSignerDetails               `json:"manage_signer"`
	ManageSignerRule           *ManageSignerRuleDetails           `json:"manage_signer_rule"`
	ManageSignerRole           *ManageSignerRoleDetails           `json:"manage_signer_role"`
	ManageBalance              *ManageBalanceDetails              `json:"manage_balance,omitempty"`
	ManageKeyValue             *ManageKeyValueDetails             `json:"manage_key_value,omitempty"`
	ManageAsset                *ManageAssetDetails                `json:"manage_asset,omitempty"`
	ManageAssetPair            *ManageAssetPairDetails            `json:"manage_asset_pair,omitempty"`
	ManageLimits               *ManageLimitsDetails               `json:"manage_limits,omitempty"`
	ManageOffer                *ManageOfferDetails                `json:"manage_offer,omitempty"`
	ManageExternalSystemPool   *ManageExternalSystemPoolDetails   `json:"manage_external_system_pool,omitempty"`
	BindExternalSystemAccount  *BindExternalSystemAccountDetails  `json:"bind_external_system_account,omitempty"`
	Payment                    *PaymentDetails                    `json:"payment,omitempty"`
	Payout                     *PayoutDetails                     `json:"payout,omitempty"`
	SetFee                     *SetFeeDetails                     `json:"set_fee,omitempty"`
	CancelAtomicSwapBid        *CancelAtomicSwapBidDetails        `json:"cancel_atomic_swap_bid,omitempty"`
	CheckSaleState             *CheckSaleStateDetails             `json:"check_sale_state,omitempty"`
	CreateChangeRoleRequest    *CreateChangeRoleRequestDetails    `json:"create_change_role_request,omitempty"`
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
	License                    *LicenseDetails                    `json:"license,omitempty"`
	Stamp                      *StampDetails                      `json:"stamp,omitempty"`
	ManageCreatePollRequest    *ManageCreatePollRequestDetails    `json:"manage_create_poll_request,omitempty"`
	ManagePoll                 *ManagePollDetails                 `json:"manage_poll,omitempty"`
	ManageVote                 *ManageVoteDetails                 `json:"manage_vote,omitempty"`
}

//Value - converts operation details into jsonb
func (r OperationDetails) Value() (driver.Value, error) {
	result, err := db2.DriverValue(r)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal operation details")
	}

	return result, nil
}

//Scan - converts jsonb into OperationDetails
func (r *OperationDetails) Scan(src interface{}) error {
	err := db2.DriveScan(src, r)
	if err != nil {
		return errors.Wrap(err, "failed to scan operation details")
	}

	return nil
}

// CreateAccountDetails - stores details of create account operation
type CreateAccountDetails struct {
	AccountAddress string `json:"account_address"`
	AccountRole    uint64 `json:"account_role"`
}

//ManageBalanceDetails - details of ManageBalanceOp
type ManageBalanceDetails struct {
	DestinationAccount string                  `json:"destination_account"`
	Action             xdr.ManageBalanceAction `json:"action"`
	Asset              string                  `json:"asset"`
	BalanceAddress     string                  `json:"balance_address"`
}

// ManageAccountRuleDetails - details of ManageAccountRuleOp
type ManageAccountRuleDetails struct {
	Action        xdr.ManageAccountRuleAction `json:"action"`
	RuleID        uint64                      `json:"rule_id"`
	CreateDetails *UpdateAccountRuleDetails   `json:"create_details,omitempty"`
	UpdateDetails *UpdateAccountRuleDetails   `json:"update_details,omitempty"`
}

// UpdateAccountRuleDetails - details of new or updated rule
type UpdateAccountRuleDetails struct {
	Resource xdr.AccountRuleResource `json:"resource"`
	Action   xdr.AccountRuleAction   `json:"action"`
	IsForbid bool                    `json:"is_forbid"`
	Details  rgenerated.Details      `json:"details"`
}

// ManageAccountRoleDetails - details of ManageAccountRoleOp
type ManageAccountRoleDetails struct {
	Action        xdr.ManageAccountRoleAction `json:"action"`
	RoleID        uint64                      `json:"role_id"`
	CreateDetails *UpdateAccountRoleDetails   `json:"create_details,omitempty"`
	UpdateDetails *UpdateAccountRoleDetails   `json:"update_details,omitempty"`
}

// UpdateAccountRoleDetails - details of new or updated role
type UpdateAccountRoleDetails struct {
	RuleIDs []uint64           `json:"rule_ids"`
	Details rgenerated.Details `json:"details"`
}

// ManageSignerDetails - details op manage signer operation
type ManageSignerDetails struct {
	Action        xdr.ManageSignerAction `json:"action"`
	PublicKey     string                 `json:"public_key"`
	CreateDetails *UpdateSignerDetails   `json:"create_details,omitempty"`
	UpdateDetails *UpdateSignerDetails   `json:"update_details,omitempty"`
}

// UpdateSignerDetails - details of new or updated signer
type UpdateSignerDetails struct {
	RoleID   uint64             `json:"role_id"`
	Weight   uint32             `json:"weight"`
	Identity uint32             `json:"identity"`
	Details  rgenerated.Details `json:"details"`
}

// ManageSignerRuleDetails - details of ManageAccountRuleOp
type ManageSignerRuleDetails struct {
	Action        xdr.ManageSignerRuleAction `json:"action"`
	RuleID        uint64                     `json:"rule_id"`
	CreateDetails *CreateSignerRuleDetails   `json:"create_details,omitempty"`
	UpdateDetails *UpdateSignerRuleDetails   `json:"update_details,omitempty"`
}

// CreateSignerRuleDetails - details of new or updated rule
type CreateSignerRuleDetails struct {
	Resource   xdr.SignerRuleResource `json:"resource"`
	Action     xdr.SignerRuleAction   `json:"action"`
	IsForbid   bool                   `json:"is_forbid"`
	IsDefault  bool                   `json:"is_default"`
	IsReadOnly bool                   `json:"is_read_only"`
	Details    rgenerated.Details     `json:"details"`
}

// UpdateSignerRuleDetails - details of new or updated rule
type UpdateSignerRuleDetails struct {
	Resource  xdr.SignerRuleResource `json:"resource"`
	Action    xdr.SignerRuleAction   `json:"action"`
	IsForbid  bool                   `json:"is_forbid"`
	IsDefault bool                   `json:"is_default"`
	Details   rgenerated.Details     `json:"details"`
}

// ManageSignerRoleDetails - details of ManageAccountRoleOp
type ManageSignerRoleDetails struct {
	Action        xdr.ManageSignerRoleAction `json:"action"`
	RoleID        uint64                     `json:"role_id"`
	CreateDetails *CreateSignerRoleDetails   `json:"create_details,omitempty"`
	UpdateDetails *UpdateSignerRoleDetails   `json:"update_details,omitempty"`
}

// UpdateSignerRoleDetails - details of new or updated role
type CreateSignerRoleDetails struct {
	RuleIDs    []uint64           `json:"rule_ids"`
	IsReadOnly bool               `json:"is_read_only"`
	Details    rgenerated.Details `json:"details"`
}

// UpdateSignerRoleDetails - details of new or updated role
type UpdateSignerRoleDetails struct {
	RuleIDs []uint64           `json:"rule_ids"`
	Details rgenerated.Details `json:"details"`
}

//ManageKeyValueDetails - details of ManageKeyValueOp
type ManageKeyValueDetails struct {
	Action xdr.ManageKvAction             `json:"action"`
	Key    string                         `json:"key"`
	Value  *rgenerated.KeyValueEntryValue `json:"value,omitempty"`
}

//SetFeeDetails - details of SetFeeOp
type SetFeeDetails struct {
	AccountAddress *string           `json:"account_address,omitempty"`
	AccountRole    *xdr.Uint64       `json:"account_role,omitempty"`
	AssetCode      string            `json:"asset_code"`
	FeeType        xdr.FeeType       `json:"fee_type"`
	FixedFee       rgenerated.Amount `json:"fixed_fee"`
	IsDelete       bool              `json:"is_delete"`
	LowerBound     rgenerated.Amount `json:"lower_bound"`
	PercentFee     rgenerated.Amount `json:"percent_fee"`
	Subtype        int64             `json:"subtype"`
	UpperBound     rgenerated.Amount `json:"upper_bound"`
	// FeeAsset deprecated
}

//CreateWithdrawRequestDetails - details of corresponding op
type CreateWithdrawRequestDetails struct {
	BalanceAddress string             `json:"balance_address"`
	Amount         rgenerated.Amount  `json:"amount"`
	Fee            rgenerated.Fee     `json:"fee"`
	CreatorDetails rgenerated.Details `json:"creator_details"`
}

//CreateManageLimitsRequestDetails - details of corresponding op
type CreateManageLimitsRequestDetails struct {
	CreatorDetails rgenerated.Details `json:"creator_details"`
	RequestID      int64              `json:"request_id"`
}

//ManageLimitsDetails - details of corresponding op
type ManageLimitsDetails struct {
	Action   xdr.ManageLimitsAction       `json:"action"`
	Creation *ManageLimitsCreationDetails `json:"creation_details,omitempty"`
	Removal  *ManageLimitsRemovalDetails  `json:"removal_details,omitempty"`
}

//ManageLimitsCreationDetails - details of corresponding op
type ManageLimitsCreationDetails struct {
	AccountRole     *xdr.Uint64       `json:"account_role,omitempty"`
	AccountAddress  string            `json:"account_address,omitempty"`
	StatsOpType     xdr.StatsOpType   `json:"stats_op_type"`
	AssetCode       string            `json:"asset_code"`
	IsConvertNeeded bool              `json:"is_convert_needed"`
	DailyOut        rgenerated.Amount `json:"daily_out"`
	WeeklyOut       rgenerated.Amount `json:"weekly_out"`
	MonthlyOut      rgenerated.Amount `json:"monthly_out"`
	AnnualOut       rgenerated.Amount `json:"annual_out"`
}

//ManageLimitsRemovalDetails - details of corresponding op
type ManageLimitsRemovalDetails struct {
	LimitsID int64 `json:"limits_id"`
}

//ManageAssetPairDetails - details of corresponding op
type ManageAssetPairDetails struct {
	BaseAsset               string              `json:"base_asset"`
	QuoteAsset              string              `json:"quote_asset"`
	PhysicalPrice           rgenerated.Amount   `json:"physical_price"`
	PhysicalPriceCorrection rgenerated.Amount   `json:"physical_price_correction"`
	MaxPriceStep            rgenerated.Amount   `json:"max_price_step"`
	Policies                xdr.AssetPairPolicy `json:"policies"`
}

//ManageOfferDetails - details of corresponding op
type ManageOfferDetails struct {
	OfferID     int64             `json:"offer_id,omitempty"`
	OrderBookID int64             `json:"order_book_id"`
	BaseAsset   string            `json:"base_asset"`
	QuoteAsset  string            `json:"quote_asset"`
	Amount      rgenerated.Amount `json:"base_amount"`
	Price       rgenerated.Amount `json:"price"`
	IsBuy       bool              `json:"is_buy"`
	Fee         rgenerated.Fee    `json:"fee"`
	IsDeleted   bool              `json:"is_deleted"`
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
	ExternalDetails rgenerated.Details                `json:"external_details"`
}

//ManageAssetDetails - details of corresponding op
type ManageAssetDetails struct {
	AssetCode         string                `json:"asset_code,omitempty"`
	Type              uint64                `json:"type,omitempty"`
	RequestID         int64                 `json:"request_id"`
	Action            xdr.ManageAssetAction `json:"action"`
	Policies          *xdr.AssetPolicy      `json:"policies,omitempty"`
	CreatorDetails    rgenerated.Details    `json:"creator_details,omitempty"`
	PreissuedSigner   string                `json:"preissued_signer,omitempty"`
	MaxIssuanceAmount rgenerated.Amount     `json:"max_issuance_amount,omitempty"`
}

//CreatePreIssuanceRequestDetails - details of corresponding op
type CreatePreIssuanceRequestDetails struct {
	AssetCode      string             `json:"asset_code"`
	Amount         rgenerated.Amount  `json:"amount"`
	CreatorDetails rgenerated.Details `json:"creator_details"`
	RequestID      int64              `json:"request_id"`
	IsFulfilled    bool               `json:"is_fulfilled"`
}

//CreateIssuanceRequestDetails - details of corresponding op
type CreateIssuanceRequestDetails struct {
	Fee                    rgenerated.Fee     `json:"fee"`
	Reference              string             `json:"reference"`
	Amount                 rgenerated.Amount  `json:"amount"`
	Asset                  string             `json:"asset"`
	ReceiverAccountAddress string             `json:"receiver_account_address"`
	ReceiverBalanceAddress string             `json:"receiver_balance_address"`
	CreatorDetails         rgenerated.Details `json:"creator_details"`
	AllTasks               *int64             `json:"all_tasks,omitempty"`
	RequestDetails         RequestDetails     `json:"request_details"`
}

//CreateSaleRequestDetails -details of corresponding op
type CreateSaleRequestDetails struct {
	RequestID         int64                   `json:"request_id"`
	BaseAsset         string                  `json:"base_asset"`
	DefaultQuoteAsset string                  `json:"default_quote_asset"`
	StartTime         time.Time               `json:"start_time"`
	EndTime           time.Time               `json:"end_time"`
	SoftCap           rgenerated.Amount       `json:"soft_cap"`
	HardCap           rgenerated.Amount       `json:"hard_cap"`
	QuoteAssets       []rgenerated.AssetPrice `json:"quote_assets"`
	CreatorDetails    rgenerated.Details      `json:"creator_details"`
}

//CheckSaleStateDetails - details of corresponding op
type CheckSaleStateDetails struct {
	SaleID int64 `json:"sale_id"`
	Effect xdr.CheckSaleStateEffect
}

//CreateAtomicSwapBidRequestDetails - details of corresponding op
type CreateAtomicSwapBidRequestDetails struct {
	Amount         rgenerated.Amount       `json:"amount"`
	BaseBalance    string                  `json:"base_balance"`
	QuoteAssets    []rgenerated.AssetPrice `json:"quote_assets"`
	Details        rgenerated.Details      `json:"details"`
	RequestDetails RequestDetails          `json:"request_details"`
}

//CreateAtomicSwapRequestDetails - details of corresponding op
type CreateAtomicSwapRequestDetails struct {
	BidID          int64             `json:"bid_id"`
	BaseAmount     rgenerated.Amount `json:"base_amount"`
	QuoteAsset     string            `json:"quote_asset"`
	RequestDetails RequestDetails    `json:"request_details"`
}

//CancelAtomicSwapBidDetails - details of corresponding op
type CancelAtomicSwapBidDetails struct {
	BidID int64 `json:"bid_id"`
}

//CreateAMLAlertRequestDetails - details of corresponding op
type CreateAMLAlertRequestDetails struct {
	Amount         rgenerated.Amount  `json:"amount"`
	BalanceAddress string             `json:"balance_address"`
	CreatorDetails rgenerated.Details `json:"creator_details"`
}

// PaymentDetails - stores details of payment operation
type PaymentDetails struct {
	AccountFrom             string            `json:"account_from"`
	AccountTo               string            `json:"account_to"`
	BalanceFrom             string            `json:"balance_from"`
	BalanceTo               string            `json:"balance_to"`
	Amount                  rgenerated.Amount `json:"amount"`
	Asset                   string            `json:"asset"`
	SourceFee               rgenerated.Fee    `json:"source_fee"`
	DestinationFee          rgenerated.Fee    `json:"destination_fee"`
	SourcePayForDestination bool              `json:"source_pay_for_destination"`
	Subject                 string            `json:"subject"`
	Reference               string            `json:"reference"`
}

//PayoutDetails - details of corresponding op
type PayoutDetails struct {
	SourceAccountAddress string            `json:"source_account_address"`
	SourceBalanceAddress string            `json:"source_balance_address"`
	Asset                string            `json:"asset"`
	MaxPayoutAmount      rgenerated.Amount `json:"max_payout_amount"`
	MinAssetHolderAmount rgenerated.Amount `json:"min_asset_holder_amount"`
	MinPayoutAmount      rgenerated.Amount `json:"min_payout_amount"`
	ExpectedFee          rgenerated.Fee    `json:"expected_fee"`
	ActualFee            rgenerated.Fee    `json:"actual_fee"`
	ActualPayoutAmount   rgenerated.Amount `json:"actual_payout_amount"`
}

//RequestDetails - details of the request created or reviewed via op
type RequestDetails struct {
	RequestID   int64 `json:"request_id,omitempty"`
	IsFulfilled bool  `json:"is_fulfilled"`
}

//CreateChangeRoleRequestDetails - details of corresponding op
type CreateChangeRoleRequestDetails struct {
	DestinationAccount string             `json:"destination_account"`
	AccountRoleToSet   uint64             `json:"account_role_to_set"`
	CreatorDetails     rgenerated.Details `json:"creator_details"`
	AllTasks           *uint32            `json:"all_tasks"`
	RequestDetails     RequestDetails     `json:"request_details"`
}

//ManageExternalSystemPoolDetails - details of corresponding op
type ManageExternalSystemPoolDetails struct {
	Action xdr.ManageExternalSystemAccountIdPoolEntryAction `json:"action"`
	Create *CreateExternalSystemPoolDetails                 `json:"create"`
	Remove *RemoveExternalSystemPoolDetails                 `json:"remove"`
}

//CreateExternalSystemPoolDetails - details of corresponding op
type CreateExternalSystemPoolDetails struct {
	Data               string `json:"data"`
	ExternalSystemType int32  `json:"external_system_type"`
	Parent             uint64 `json:"parent"`
	PoolId             uint64 `json:"pool_id"`
}

//RemoveExternalSystemPoolDetails - details of corresponding op
type RemoveExternalSystemPoolDetails struct {
	PoolId uint64 `json:"pool_id"`
}

//BindExternalSystemAccountDetails - details of corresponding op
type BindExternalSystemAccountDetails struct {
	ExternalSystemType int32 `json:"external_system_type"`
}

//ManageSaleDetails - details of corresponding op
type ManageSaleDetails struct {
	Action xdr.ManageSaleAction `json:"action"`
	SaleId uint64               `json:"sale_id"`
}

type LicenseDetails struct {
	PrevLicenseHash string    `json:"prev_license_hash"`
	AdminCount      uint64    `json:"admin_count"`
	DueDate         time.Time `json:"due_date"`
	LedgerHash      string    `json:"ledger_hash"`
	Signatures      []string  `json:"signatures"`
}

type StampDetails struct {
	LicenseHash string `json:"license_hash"`
	LedgerHash  string `json:"ledger_hash"`
}

type CreatePollRequestDetails struct {
	PollType                 xdr.PollType      `json:"poll_type"`
	PermissionType           uint32            `json:"permission_type"`
	NumberOfChoices          uint32            `json:"number_of_choices"`
	CreatorDetails           regources.Details `json:"creator_details"`
	StartTime                time.Time         `json:"start_time"`
	EndTime                  time.Time         `json:"end_time"`
	ResultProviderID         string            `json:"result_provider_id"`
	VoteConfirmationRequired bool              `json:"vote_confirmation_required"`
	PollData                 xdr.PollData      `json:"poll_data"`
	AllTasks                 *uint32           `json:"all_tasks,omitempty"`
	RequestDetails           RequestDetails    `json:"request_details"`
}

type CancelCreatePollRequestDetails struct {
	RequestID int64 `json:"request_id"`
}

type ManageCreatePollRequestDetails struct {
	Action        xdr.ManageCreatePollRequestAction `json:"action"`
	CreateDetails *CreatePollRequestDetails         `json:"create_details,omitempty"`
	CancelDetails *CancelCreatePollRequestDetails   `json:"cancel_details,omitempty"`
}

type ClosePollData struct {
	PollResult xdr.PollResult    `json:"poll_result"`
	Details    regources.Details `json:"details"`
}

type ManagePollDetails struct {
	Action    xdr.ManagePollAction `json:"action"`
	PollID    int64                `json:"poll_id"`
	ClosePoll *ClosePollData       `json:"close_poll,omitempty"`
}

type ManageVoteDetails struct {
	PollID   int64                `json:"poll_id"`
	Action   xdr.ManageVoteAction `json:"action"`
	VoteData *xdr.VoteData        `json:"vote_data,omitmepty"`
}
