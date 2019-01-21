package regources

// Balance - Resource object representing BalanceEntry
type Balance struct {
	Key
	Relationships BalanceRelation `json:"relationships,omitempty"`
}

type BalanceRelation struct {
	Asset *Relation `json:"asset,omitempty"`
	State *Relation `json:"state,omitempty"`
}

//BalanceState - Resource represents balance state
type BalanceState struct {
	Key
	Attributes BalanceStateAttr `json:"attributes"`
}

type BalanceStateAttr struct {
	Available Amount `json:"available"`
	Locked    Amount `json:"locked"`
}
