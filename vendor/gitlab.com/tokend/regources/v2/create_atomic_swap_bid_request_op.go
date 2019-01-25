package regources

//CreateAtomicSwapBidRequestAttrs - details of corresponding op
type CreateAtomicSwapBidRequest struct {
	Key
	Attributes CreateAtomicSwapBidRequestAttrs `json:"attributes"`
}

//CreateAtomicSwapBidRequestAttrs - details of corresponding op
type CreateAtomicSwapBidRequestAttrs struct {
	Amount         Amount           `json:"amount"`
	BaseBalance    string           `json:"base_balance"`
	QuoteAssets    []SaleQuoteAsset `json:"quote_assets"`
	Details        Details          `json:"details"`
	RequestDetails Request          `json:"request_details"`
}
