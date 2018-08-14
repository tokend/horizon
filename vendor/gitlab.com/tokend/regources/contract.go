package regources

import (
	"time"

	"gitlab.com/tokend/regources/reviewablerequest2"
	"gitlab.com/tokend/regources/valueflag"
)

type Contract struct {
	ID            string                                 `json:"id"`
	PT            string                                 `json:"paging_token"`
	Contractor    string                                 `json:"contractor"`
	Customer      string                                 `json:"customer"`
	Escrow        string                                 `json:"escrow"`
	Disputer      string                                 `json:"disputer,omitempty"`
	StartTime     time.Time                              `json:"start_time"`
	EndTime       time.Time                              `json:"end_time"`
	Details       []map[string]interface{}               `json:"details"`
	Invoices      []reviewablerequest2.ReviewableRequest `json:"invoices,omitempty"`
	DisputeReason map[string]interface{}                 `json:"dispute_reason,omitempty"`
	State         []valueflag.Flag                       `json:"state"`
}

func (c Contract) PagingToken() string {
	return c.PT
}
