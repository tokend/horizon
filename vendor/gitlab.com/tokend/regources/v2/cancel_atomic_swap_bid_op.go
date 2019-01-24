package regources

//CancelAtomicSwapBidAttrs - details of corresponding op
type CancelAtomicSwapBid struct {
	Key
	Attributes CancelAtomicSwapBidAttrs `json:"attributes"`
}

//CancelAtomicSwapBidAttrs - details of corresponding op
type CancelAtomicSwapBidAttrs struct {
	BidID int64 `json:"bid_id"`
}
