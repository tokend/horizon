package regources

//CreateAtomicSwapRequest - details of corresponding op
type CreateAtomicSwapRequest struct {
	Key
	Attributes CreateAtomicSwapRequestAttrs `json:"attributes"`
}

//CreateAtomicSwapRequestAttrs - details of corresponding op
type CreateAtomicSwapRequestAttrs struct {
	BidID          int64   `json:"bid_id"`
	BaseAmount     Amount  `json:"base_amount"`
	QuoteAsset     string  `json:"quote_asset"`
	RequestDetails Request `json:"request_details"`
}
