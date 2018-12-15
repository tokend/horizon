package history2

// Account is a row of data from the `history_accounts` table
type Account struct {
	ID          int64  `db:"id"`
	Address     string `db:"address"`
	AccountType int32  `db:"account_type"`
}

func NewAccount(accountID int64, address string, accountType int32) Account {
	return Account{
		ID:          accountID,
		Address:     address,
		AccountType: int32(accountType),
	}
}


