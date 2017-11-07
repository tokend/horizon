package operations

import "bullioncoin.githost.io/development/horizon/resource/base"

type ManageForfeitRequest struct {
	Base
	Action      int32              `json:"action"`
	RequestID   uint64             `json:"request_id"`
	Amount      string             `json:"amount"`
	Asset       string             `json:"asset"`
	UserDetails string             `json:"user_details,omitempty"`
	Items       []base.ForfeitItem `json:"items"`
	FixedFee    string             `json:"fixed_fee"`
	PercentFee  string             `json:"percent_fee"`
}
