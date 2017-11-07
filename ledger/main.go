// Package ledger provides useful utilities concerning ledgers within stellar,
// specifically a central location to store a cached snapshot of the state of
// both horizon's and stellar-core's view of the ledger.  This package is
// intended to be at the lowest levels of horizon's dependency tree, please keep
// it free of dependencies to other horizon packages.
package ledger

import (
	"fmt"
	"sync"
	"time"
)

// State represents a snapshot of both horizon's and stellar-core's view of the
// ledger.
type State struct {
	CoreLatest    int32 `db:"core_latest"`
	CoreElder     int32 `db:"core_elder"`
	HistoryLatest int32 `db:"history_latest"`
	HistoryElder  int32 `db:"history_elder"`
}

// CurrentState returns the cached snapshot of ledger state
func CurrentState() State {
	lock.RLock()
	ret := current
	lock.RUnlock()
	return ret
}

// SetState updates the cached snapshot of the ledger state
func SetState(next State) error {
	lock.Lock()
	defer lock.Unlock()

	currentLatest := current.CoreLatest
	current = next

	if next.CoreLatest > currentLatest {
		lastLedgerIncreaseTime = time.Now().UTC()
		return nil
	}

	if time.Now().UTC().Sub(lastLedgerIncreaseTime) > 30*time.Second {
		return fmt.Errorf("core latest ledger number doesn't increase. Latest ledger: %d", next.CoreLatest)
	}

	return nil
}

var current State
var lock sync.RWMutex

var lastLedgerIncreaseTime time.Time
