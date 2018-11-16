package history2

import (
	"gitlab.com/tokend/go/xdr"
)

// OperationDetails - stores details of the operation performed in union switch form. Only one value must be selected at
// a type
type OperationDetails struct {
	Type                  xdr.OperationType             `json:"type"`
	CreateAccount         *CreateAccountDetails         `json:"create_account,omitempty"`
	Payment               *PaymentDetails               `json:"payment,omitempty"`
	ManageKeyValue        *ManageKeyValueDetails        `json:"manage_key_value,omitempty"`
	SetFee                *SetFeeDetails                `json:"set_fee,omitempty"`
	ManageAccount         *ManageAccountDetails         `json:"manage_account,omitempty"`
	CreateWithdrawRequest *CreateWithdrawRequestDetails `json:"create_withdraw_request,omitempty"`
}

// CreateAccountDetails - stores details of create account operation
type CreateAccountDetails struct {
	AccountType xdr.AccountType `json:"account_type"`
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
	FixedFee    int64            `json:"fixed_fee"`
	PercentFee  int64            `json:"percent_fee"`
	FeeType     xdr.FeeType      `json:"fee_type"`
	AccountID   *int64           `json:"account_id,omitempty"`
	AccountType *xdr.AccountType `json:"account_type,omitempty"`
	Subtype     int64            `json:"subtype"`
	LowerBound  int64            `json:"lower_bound"`
	UpperBound  int64            `json:"upper_bound"`
	// FeeAsset deprecated
}

type CreateWithdrawRequestDetails struct {
	BalanceID         int64                  `json:"balance_id"`
	Amount            int64                  `json:"amount"`
	FixedFee          int64                  `json:"fixed_fee"`
	PercentFee        int64                  `json:"percent_fee"`
	ExternalDetails   map[string]interface{} `json:"external_details"`
	DestinationAsset  xdr.AssetCode          `json:"destination_asset"`
	DestinationAmount int64                  `json:"destination_amount"`
}

// PaymentDetails - stores details of payment operation
type PaymentDetails struct {
}
