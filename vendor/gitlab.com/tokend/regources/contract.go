package regources

import "time"

type Contract struct {
	ID            string                   `json:"id"`
	PT            string                   `json:"paging_token"`
	Contractor    string                   `json:"contractor"`
	Customer      string                   `json:"customer"`
	Escrow        string                   `json:"escrow"`
	Disputer      string                   `json:"disputer,omitempty"`
	StartTime     time.Time                `json:"start_time"`
	EndTime       time.Time                `json:"end_time"`
	Details       []map[string]interface{} `json:"details"`
	DisputeReason map[string]interface{}   `json:"dispute_reason,omitempty"`
	State         []Flag                   `json:"state"`
}

func (c Contract) PagingToken() string {
	return c.PT
}
