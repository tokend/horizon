package regources

type CreateAtomicSwapRequestRelationships struct {
	Bid        *Relation `json:"bid,omitempty"`
	QuoteAsset *Relation `json:"quote_asset,omitempty"`
}
