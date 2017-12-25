package operations

type CreateIssuanceRequest struct {
	Base
	Reference       string                 `json:"reference"`
	Amount          string                 `json:"amount"`
	Asset           string                 `json:"asset"`
	ExternalDetails map[string]interface{} `json:"external_details"`
}
