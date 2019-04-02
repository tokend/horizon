package regources

import (
	"time"

	"gitlab.com/tokend/go/xdr"
)

type SaleAttributes struct {
	Details Details `json:"details"`
	// time when the sale expires
	EndTime time.Time `json:"end_time"`
	// state of sale
	SaleState SaleState `json:"sale_state"`
	// type of sale
	SaleType xdr.SaleType `json:"sale_type"`
	// time when the sale starts
	StartTime time.Time `json:"start_time"`
}
