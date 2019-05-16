/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package regources

import (
	"time"
)

type AtomicSwapBidAttributes struct {
	// Amount that can be bought through atomic swap request
	AvailableAmount Amount `json:"available_amount"`
	// time when the atomic swap bid was created
	CreatedAt time.Time `json:"created_at"`
	// represents user-provided data
	Details Details `json:"details"`
	// defines whether creating atomic swap requests for this bid is forbidden
	IsCanceled bool `json:"is_canceled"`
	// Amount that that is being processed now through atomic swap requests
	LockedAmount Amount `json:"locked_amount"`
}
