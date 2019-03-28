package resources

import (
	core "gitlab.com/tokend/horizon/db2/core2"
	"gitlab.com/tokend/regources/rgenerated"
)

//NewBalance - creates new instance of balance using core balance
func NewBalance(record *core.Balance) *rgenerated.Balance {
	return &rgenerated.Balance{
		Key: NewBalanceKey(record.BalanceAddress),
	}
}

//NewBalanceKey - creates new instance of balance key
func NewBalanceKey(balanceAddress string) rgenerated.Key {
	return rgenerated.Key{
		Type: rgenerated.BALANCES,
		ID:   balanceAddress,
	}
}

//NewBalanceState - returns new balance state
func NewBalanceState(record *core.Balance) *rgenerated.BalanceState {
	return &rgenerated.BalanceState{
		Key: rgenerated.Key{
			ID:   record.BalanceAddress,
			Type: rgenerated.BALANCES_STATE,
		},
		Attributes: &rgenerated.BalanceStateAttributes{
			Locked:    rgenerated.Amount(record.Locked),
			Available: rgenerated.Amount(record.Amount),
		},
	}
}

//NewBalanceStateKey - creates new balance state key using balance address
func NewBalanceStateKey(balanceAddress string) rgenerated.Key {
	return rgenerated.Key{
		ID:   balanceAddress,
		Type: rgenerated.BALANCES_STATE,
	}
}
