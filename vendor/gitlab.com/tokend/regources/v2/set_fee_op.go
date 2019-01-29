package regources

import "gitlab.com/tokend/go/xdr"

//SetFee - stores details of create account operation
type SetFee struct {
	Key
	Attributes SetFeeAttrs `json:"attributes"`
}

//SetFeeAttrs - details of SetFeeOp
type SetFeeAttrs struct {
	AssetCode      string           `json:"asset_code"`
	FixedFee       Amount           `json:"fixed_fee"`
	PercentFee     Amount           `json:"percent_fee"`
	FeeType        xdr.FeeType      `json:"fee_type"`
	AccountAddress *string          `json:"account_address,omitempty"`
	AccountType    *xdr.AccountType `json:"account_type,omitempty"`
	Subtype        int64            `json:"subtype"`
	LowerBound     Amount           `json:"lower_bound"`
	UpperBound     Amount           `json:"upper_bound"`
	IsDelete       bool             `json:"is_delete"`
	// FeeAsset deprecated
}
