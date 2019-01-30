package regources

// AtomicSwapBidRequest - represents details of the `atomic swap bid` reviewable request
type AtomicSwapBidRequest struct {
	Key
	Attributes    AtomicSwapBidRequestAttrs     `json:"attributes"`
	Relationships AtomicSwapBidRequestRelations `json:"relationships"`
}

// AtomicSwapBidRequestAttrs - attributes of the `atomic swap bid` reviewable request
type AtomicSwapBidRequestAttrs struct {
	BaseAmount Amount  `json:"base_amount"`
	Details    Details `json:"details"`
}

// AtomicSwapBidRequestRelations - relationships of the `atomic swap bid` reviewable request
type AtomicSwapBidRequestRelations struct {
	BaseBalance *Relation           `json:"base_balance"`
	QuoteAssets *RelationCollection `json:"quote_assets"`
}
