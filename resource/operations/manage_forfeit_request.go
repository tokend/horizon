package operations

type ManageForfeitRequest struct {
	Base
	Action      int32              `json:"action"`
	RequestID   uint64             `json:"request_id"`
	Amount      string             `json:"amount"`
	Asset       string             `json:"asset"`
	UserDetails string             `json:"user_details,omitempty"`
	TotalFee	string			   `json:"total_fee"`
}
