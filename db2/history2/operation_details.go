package history2

import (
	"gitlab.com/tokend/go/xdr"
)

// OperationDetails - stores details of the operation performed in union switch form. Only one value must be selected at
// a type
type OperationDetails struct {
	Type                      xdr.OperationType                 `json:"type"`
	CreateAccount             *CreateAccountDetails             `json:"create_account,omitempty"`
	ManageAccount             *ManageAccountDetails             `json:"manage_account,omitempty"`
	ManageBalance             *ManageBalanceDetails             `json:"manage_balance,omitempty"`
	Payment                   *PaymentDetails                   `json:"payment,omitempty"`
	ManageKeyValue            *ManageKeyValueDetails            `json:"manage_key_value,omitempty"`
	SetFee                    *SetFeeDetails                    `json:"set_fee,omitempty"`
	ManageAssetPair           *ManageAssetPairDetails           `json:"manage_asset_pair"`
	ManageLimits              *ManageLimitsDetails              `json:"manage_limits,omitempty"`
	ManageOffer               *ManageOfferDetails               `json:"manage_offer"`
	CreateManageLimitsRequest *CreateManageLimitsRequestDetails `json:"create_manage_limits_request"`
	CreateWithdrawRequest     *CreateWithdrawRequestDetails     `json:"create_withdraw_request,omitempty"`
	ManageInvoiceRequest      *ManageInvoiceRequestDetails      `json:"manage_invoice_request"`
}

// CreateAccountDetails - stores details of create account operation
type CreateAccountDetails struct {
	AccountType xdr.AccountType `json:"account_type"`
}

type ManageBalanceDetails struct {
	DestinationAccount int64                   `json:"destination_account"`
	Action             xdr.ManageBalanceAction `json:"action"`
	Asset              xdr.AssetCode           `json:"asset"`
	BalanceID          int64                   `json:"balance_id"`
}

type ManageAccountDetails struct {
	AccountID            int64 `json:"account_id"`
	BlockReasonsToAdd    int32 `json:"block_reasons_to_add"`
	BlockReasonsToRemove int32 `json:"block_reasons_to_remove"`
}

type ManageKeyValueDetails struct {
	Key    string                  `json:"key"`
	Action xdr.ManageKvAction      `json:"action"`
	Value  *xdr.KeyValueEntryValue `json:"value,omitempty"`
}

type SetFeeDetails struct {
	AssetCode   xdr.AssetCode    `json:"asset_code"`
	FixedFee    string           `json:"fixed_fee"`
	PercentFee  string           `json:"percent_fee"`
	FeeType     xdr.FeeType      `json:"fee_type"`
	AccountID   *int64           `json:"account_id,omitempty"`
	AccountType *xdr.AccountType `json:"account_type,omitempty"`
	Subtype     int64            `json:"subtype"`
	LowerBound  string           `json:"lower_bound"`
	UpperBound  string           `json:"upper_bound"`
	// FeeAsset deprecated
}

type CreateWithdrawRequestDetails struct {
	BalanceID         int64                  `json:"balance_id"`
	Amount            string                 `json:"amount"`
	FixedFee          string                 `json:"fixed_fee"`
	PercentFee        string                 `json:"percent_fee"`
	ExternalDetails   map[string]interface{} `json:"external_details"`
	DestinationAsset  xdr.AssetCode          `json:"destination_asset"`
	DestinationAmount string                 `json:"destination_amount"`
}

type CreateManageLimitsRequestDetails struct {
	Data      map[string]interface{} `json:"data"`
	RequestID int64                  `json:"request_id"`
}

type ManageInvoiceRequestDetails struct {
	Action xdr.ManageInvoiceRequestAction `json:"action"`
	Create *CreateInvoiceRequestDetails   `json:"create,omitempty"`
	Remove *RemoveInvoiceRequestDetails   `json:"remove,omitempty"`
}
type CreateInvoiceRequestDetails struct {
	Amount     string                 `json:"amount"`
	Sender     int64                  `json:"sender"`
	RequestID  int64                  `json:"request_id"`
	Asset      xdr.AssetCode          `json:"asset"`
	ContractID *int64                 `json:"contract_id,omitempty"`
	Details    map[string]interface{} `json:"details"`
}
type RemoveInvoiceRequestDetails struct {
	RequestID int64 `json:"request_id"`
}

type ManageContractRequestDetails struct {
	Action xdr.ManageContractRequestAction
	/*Create
	Remove*/
}
type CreateContractRequestDetails struct {
	Escrow    int64                  `json:"escrow"`
	Details   map[string]interface{} `json:"details"`
	StartTime int64                  `json:"start_time"`
	EndTime   int64                  `json:"end_time"`
	RequestID int64                  `json:"request_id"`
}
type RemoveContractReqeustDetails struct {
	RequestID int64 `json:"request_id"`
}

type ManageLimitsDetails struct {
	Action   xdr.ManageLimitsAction       `json:"action"`
	Creation *ManageLimitsCreationDetails `json:"creation_details,omitempty"`
	Removal  *ManageLimitsRemovalDetails  `json:"removal_details,omitempty"`
}
type ManageLimitsCreationDetails struct {
	AccountType     *xdr.AccountType `json:"account_type,omitempty"`
	AccountID       *int64           `json:"account_id,omitempty"`
	StatsOpType     xdr.StatsOpType  `json:"stats_op_type"`
	AssetCode       xdr.AssetCode    `json:"asset_code"`
	IsConvertNeeded bool             `json:"is_convert_needed"`
	DailyOut        string           `json:"daily_out"`
	WeeklyOut       string           `json:"weekly_out"`
	MonthlyOut      string           `json:"monthly_out"`
	AnnualOut       string           `json:"annual_out"`
}
type ManageLimitsRemovalDetails struct {
	ID int64 `json:"id"`
}

type ManageAssetPairDetails struct {
	BaseAsset               xdr.AssetCode `json:"base_asset"`
	QuoteAsset              xdr.AssetCode `json:"quote_asset"`
	PhysicalPrice           string        `json:"physical_price"`
	PhysicalPriceCorrection string        `json:"physical_price_correction"`
	MaxPriceStep            string        `json:"max_price_step"`
	PoliciesI               int32         `json:"policies_i"`
}

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

// PaymentDetails - stores details of payment operation
type PaymentDetails struct {
}
