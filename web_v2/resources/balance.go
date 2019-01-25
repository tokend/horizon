package resources

import (
	core "gitlab.com/tokend/horizon/db2/core2"
	"gitlab.com/tokend/regources/v2"
)

//NewBalance - creates new instance of balance using core balance
func NewBalance(record *core.Balance) *regources.Balance {
	return &regources.Balance{
		Key: NewBalanceKey(record.BalanceAddress),
	}
}

//NewBalanceKey - creates new instance of balance key
func NewBalanceKey(balanceAddress string) regources.Key {
	return regources.Key{
		Type: regources.TypeBalances,
		ID:   balanceAddress,
	}
}

//NewBalanceState - returns new balance state
func NewBalanceState(record *core.Balance) *regources.BalanceState {
	return &regources.BalanceState{
		Key: regources.Key{
			ID:   record.BalanceAddress,
			Type: regources.TypeBalancesState,
		},
		Attributes: regources.BalanceStateAttr{
			Locked:    regources.Amount(record.Locked),
			Available: regources.Amount(record.Amount),
		},
	}
}
