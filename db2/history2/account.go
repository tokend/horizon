package history2

import "gitlab.com/tokend/go/xdr"

// Account is a row of data from the `history_accounts` table
type Account struct {
	ID          int64  `db:"id"`
	Address     string `db:"address"`
	AccountType int32  `db:"account_type"`
}

func NewAccount(account xdr.AccountEntry, ledgerSeq int32, ledgerOperationSeq int32) Account {
	return Account{
		ID:          accountID(ledgerSeq, ledgerOperationSeq),
		Address:     account.AccountId.Address(),
		AccountType: int32(account.AccountType),
	}
}

func accountID(ledgerSeq int32, ledgerOperationSeq int32) int64 {
	// we should never end up in situation when there is more than 16777216 in one ledger
	// we should never have operation which creates more than one account
	return int64(ledgerSeq)<<24 | int64(ledgerOperationSeq)
}
