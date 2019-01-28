package regources

// AtomicSwapRequest - represents details of the `atomic swap` reviewable request
type AtomicSwapRequest struct {
	Key
	Attributes    AtomicSwapRequestAttrs     `json:"attributes"`
	Relationships AtomicSwapRequestRelations `json:"relationships"`
}

// AtomicSwapRequestAttrs - attributes of the `atomic swap` reviewable request
type AtomicSwapRequestAttrs struct {
	BaseAmount uint64 `json:"base_amount"`
}

// AtomicSwapRequestRelations - relationships of the `atomic swap` reviewable request
type AtomicSwapRequestRelations struct {
	Bid        *Relation `json:"bid"`
	QuoteAsset *Relation `json:"quote_asset"`
}
