package regources

type BalanceRelationships struct {
	Asset *Relation `json:"asset,omitempty"`
	State *Relation `json:"state,omitempty"`
}
