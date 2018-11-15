package history2

import "gitlab.com/tokend/go/xdr"

type Balance struct {
	ID        int64  `db:"id"`
	AccountID int64  `db:"account_id"`
	BalanceID string `db:"address"`
	AssetCode string `db:"asset_code"`
}

func NewBalance(ledgerSeq, balanceSeq int32, accountID int64, balance xdr.BalanceEntry) Balance {
	return Balance{
		ID:        balanceID(ledgerSeq, balanceSeq),
		AccountID: accountID,
		BalanceID: balance.BalanceId.AsString(),
		AssetCode: string(balance.Asset),
	}
}

func balanceID(ledgerSeq int32, balanceSeq int32) int64 {
	// we should never end up in situation when there is more than 16777216 new balances in one ledger
	return int64(ledgerSeq)<<24 | int64(balanceSeq)
}
