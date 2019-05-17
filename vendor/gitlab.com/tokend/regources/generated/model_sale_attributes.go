/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package regources

import "gitlab.com/tokend/go/xdr"
import (
	"time"
)

type SaleAttributes struct {
	// indicates access definition type of the sale
	AccessDefinitionType string `json:"access_definition_type"`
	// Hard cap calculated in base asset
	BaseHardCap Amount  `json:"base_hard_cap"`
	Details     Details `json:"details"`
	// time when the sale expires
	EndTime time.Time `json:"end_time"`
	// state of sale
	SaleState SaleState `json:"sale_state"`
	// type of sale
	SaleType xdr.SaleType `json:"sale_type"`
	// time when the sale starts
	StartTime time.Time `json:"start_time"`
}
