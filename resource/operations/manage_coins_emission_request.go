package operations

// ManageCoinsEmissionRequest is the json resource representing a single operation whose type
// is ManageCoinsEmissionRequestOp.
type ManageCoinsEmissionRequest struct {
	Base
	RequestID int64  `json:"request_id"`
	Amount    string `json:"amount"`
	Asset     string `json:"asset"`
}
