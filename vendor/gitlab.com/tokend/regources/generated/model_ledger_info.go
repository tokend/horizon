/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import (
	"time"
)

type LedgerInfo struct {
	// time at which latest known ledger has been increased last time
	LastLedgerIncreaseTime *time.Time `json:"last_ledger_increase_time,omitempty"`
	// latest known ledger
	Latest *uint64 `json:"latest,omitempty"`
	// sequence of oldest ledger available
	OldestOnStart *uint64 `json:"oldest_on_start,omitempty"`
}
