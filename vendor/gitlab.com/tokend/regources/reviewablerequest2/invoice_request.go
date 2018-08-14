package reviewablerequest2

type InvoiceRequest struct {
	Amount     string                 `json:"amount"`
	Asset      string                 `json:"asset"`
	ContractID string                 `json:"contract_id,omitempty"`
	Details    map[string]interface{} `json:"details"`
}
