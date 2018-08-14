package reviewablerequest2

type IssuanceRequest struct {
	Asset           string                 `json:"asset"`
	Amount          string                 `json:"amount"`
	Receiver        string                 `json:"receiver"`
	ExternalDetails map[string]interface{} `json:"external_details"`
}
