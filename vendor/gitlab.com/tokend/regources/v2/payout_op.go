package regources

//PayoutAttrs - details of corresponding op
type Payout struct {
	Key
	Attributes PayoutAttrs `json:"attributes"`
}

//PayoutAttrs - details of corresponding op
type PayoutAttrs struct {
	SourceAccountAddress string `json:"source_account_address"`
	SourceBalanceAddress string `json:"source_balance_address"`
	Asset                string `json:"asset"`
	MaxPayoutAmount      Amount `json:"max_payout_amount"`
	MinAssetHolderAmount Amount `json:"min_asset_holder_amount"`
	MinPayoutAmount      Amount `json:"min_payout_amount"`
	ExpectedFee          Fee    `json:"expected_fee"`
	ActualFee            Fee    `json:"actual_fee"`
	ActualPayoutAmount   Amount `json:"actual_payout_amount"`
}
