package history2

import "gitlab.com/tokend/go/xdr"

type Balance struct {
	ID             int64  `db:"id"`
	AccountID      int64  `db:"account_id"`
	BalanceAddress string `db:"address"`
	AssetCode      string `db:"asset_code"`
}

func NewBalance(balanceID int64, accountID int64, balance xdr.BalanceEntry) Balance {
	return Balance{
		ID:             balanceID,
		AccountID:      accountID,
		BalanceAddress: balance.BalanceId.AsString(),
		AssetCode:      string(balance.Asset),
	}
}
