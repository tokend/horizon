package rgenerated

type ParticularBalanceChangeEffect struct {
	Amount         Amount `json:"amount"`
	AssetCode      string `json:"asset_code"`
	BalanceAddress string `json:"balance_address"`
	Fee            Fee    `json:"fee"`
}
