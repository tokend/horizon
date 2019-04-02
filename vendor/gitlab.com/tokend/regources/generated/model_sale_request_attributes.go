package regources

import (
	"time"

	"gitlab.com/tokend/go/xdr"
)

type SaleRequestAttributes struct {
	Details *Details `json:"details,omitempty"`
	// End time of the sale
	EndTime *time.Time `json:"end_time,omitempty"`
	// * 1 - open * 2 - closed * 3 - cancelled
	SaleState *SaleState `json:"sale_state,omitempty"`
	// * 1 - basic sale * 2 - crowdfunding sale * 3 - fixed price sale
	SaleType *xdr.SaleType `json:"sale_type,omitempty"`
	// Start time of the sale
	StartTime *time.Time `json:"start_time,omitempty"`
}
