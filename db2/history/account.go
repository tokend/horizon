package history

// Account is a row of data from the `accounts` table
type Account struct {
	ID          int64  `db:"id"`
	AccountID   string `db:"account_id"`
	AccountType int32  `db:"account_type"`
}
