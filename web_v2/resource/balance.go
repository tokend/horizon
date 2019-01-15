package resource

import (
	"gitlab.com/tokend/go/amount"
	core "gitlab.com/tokend/horizon/db2/core2"
)

// Balance - resource object representing BalanceEntry
type Balance struct {
	ID    string        `jsonapi:"primary,balances"`
	Asset *Asset        `jsonapi:"relation,asset,omitempty"`
	State *BalanceState `jsonapi:"relation,state,omitempty"`
}

//NewBalance - creates new instance of balance using core balance
func NewBalance(record *core.Balance) *Balance {
	return &Balance{
		ID:    record.BalanceAddress,
		Asset: NewAsset(record.Asset),
	}
}

//BalanceState - resource represents balance state
type BalanceState struct {
	ID        string `jsonapi:"primary,balance_states"`
	Available string `jsonapi:"attr,available"`
	Locked    string `jsonapi:"attr,locked"`
}

//NewBalanceState - returns new balance state
func NewBalanceState(record *core.Balance) *BalanceState {
	return &BalanceState{
		ID:        record.BalanceAddress,
		Locked:    amount.String(record.Locked),
		Available: amount.String(record.Amount),
	}
}
