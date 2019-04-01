package rgenerated

type CreateAtomicSwapBidRequestOpRelationships struct {
	BaseBalance *Relation           `json:"base_balance,omitempty"`
	QuoteAssets *RelationCollection `json:"quote_assets,omitempty"`
	Request     *Relation           `json:"request,omitempty"`
}
