package regources

// Balance - resource object representing BalanceEntry
type Balance struct {
	ID    string        `jsonapi:"primary,balances"`
	Asset *Asset        `jsonapi:"relation,asset,omitempty"`
	State *BalanceState `jsonapi:"relation,state,omitempty"`
}

//BalanceState - resource represents balance state
type BalanceState struct {
	ID        string `jsonapi:"primary,balance_states"`
	Available Amount `jsonapi:"attr,available"`
	Locked    Amount `jsonapi:"attr,locked"`
}