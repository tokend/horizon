package history

type Balance struct {
	ID           int64  `db:"id"`
	BalanceID    string `db:"balance_id"`
	AccountID    string `db:"account_id"`
	ExchangeID   string `db:"exchange_id"`
	Asset        string `db:"asset"`
	ExchangeName string `db:"exchange_name"`
}
