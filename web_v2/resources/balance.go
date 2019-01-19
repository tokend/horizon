package resources

import (
	core "gitlab.com/tokend/horizon/db2/core2"
	"gitlab.com/tokend/regources/v2"
)

//NewBalance - creates new instance of balance using core balance
func NewBalance(record *core.Balance) *regources.Balance {
	return &regources.Balance{
		ID:    record.BalanceAddress,
		Asset: NewAsset(record.Asset),
	}
}

//NewBalanceState - returns new balance state
func NewBalanceState(record *core.Balance) *regources.BalanceState {
	return &regources.BalanceState{
		ID:        record.BalanceAddress,
		Locked:    regources.Amount(record.Locked),
		Available: regources.Amount(record.Amount),
	}
}
