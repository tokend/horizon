package resource

import (
	"gitlab.com/tokend/go/amount"
	core "gitlab.com/tokend/horizon/db2/core2"
)

// Balance - resource object representing BalanceEntry
type Balance struct {
	ID        string `jsonapi:"primary,balances"`
	Available string `jsonapi:"attr,available"`
	Locked    string `jsonapi:"attr,locked"`
	Asset     *Asset `jsonapi:"relation,asset,omitempty"`
}

//NewBalance - creates new instance of balance using core balance
func NewBalance(record *core.Balance) Balance {
	return Balance{
		ID:        record.BalanceAddress,
		Locked:    amount.String(record.Locked),
		Available: amount.String(record.Amount),
		Asset:     NewAsset(record.Asset),
	}
}
