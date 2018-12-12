package history2

// Account is a row of data from the `history_accounts` table
type Account struct {
	ID          int64  `db:"id"`
	Address     string `db:"address"`
	AccountType int32  `db:"account_type"`
}

func NewAccount(address string, accountType int32, ledgerSeq int32, ledgerOperationSeq int32) Account {
	return Account{
		ID:          accountID(ledgerSeq, ledgerOperationSeq),
		Address:     address,
		AccountType: int32(accountType),
	}
}

func accountID(ledgerSeq int32, ledgerOperationSeq int32) int64 {
	// we should never end up in situation when there is more than 16777216 in one ledger
	// we should never have operation which creates more than one account
	return int64(ledgerSeq)<<24 | int64(ledgerOperationSeq)
}
