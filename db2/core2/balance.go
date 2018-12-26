package core2

type Balance struct {
	BalanceAddress string `json:"balance_id"`
	Asset          string `json:"asset"`
	AccountAddress string `json:"account_id"`
	BalanceSeqID   uint64 `json:"sequence_id"`
}
