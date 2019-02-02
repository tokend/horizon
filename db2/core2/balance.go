package core2

// Balance - the db representation of balance
type Balance struct {
	BalanceAddress string `db:"balance_id"`
	AssetCode      string `db:"asset"`
	AccountAddress string `db:"account_id"`
	BalanceSeqID   uint64 `db:"sequential_id"`
	Amount         int64  `db:"amount"`
	Locked         int64  `db:"locked"`
	*Asset         `db:"assets"`
}
