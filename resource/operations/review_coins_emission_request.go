package operations

// ReviewCoinsEmissionRequest is the json resource representing a single operation whose type
// is ManageCoinsEmissionRequestOp.
type ReviewCoinsEmissionRequest struct {
	Base
	RequestID  uint64 `json:"request_id"`
	Issuer     string `json:"issuer,omitempty"`
	Amount     string `json:"amount"`
	IsApproved bool   `json:"approved"`
	Reason     string `json:"reason,omitempty"`
	Asset      string `json:"asset"`
}
