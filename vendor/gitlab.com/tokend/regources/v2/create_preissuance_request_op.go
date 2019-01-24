package regources

//CreatePreIssuanceRequest - details of corresponding op
type CreatePreIssuanceRequest struct {
	Key
	Attributes CreatePreIssuanceRequestAttrs `json:"attributes"`
}

//CreatePreIssuanceRequestAttrs - details of corresponding op
type CreatePreIssuanceRequestAttrs struct {
	AssetCode   string `json:"asset_code"`
	Amount      Amount `json:"amount"`
	RequestID   int64  `json:"request_id"`
	IsFulfilled bool   `json:"is_fulfilled"`
}
