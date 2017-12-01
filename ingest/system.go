package ingest

import (
	"gitlab.com/swarmfund/horizon/db2/core"
	"gitlab.com/swarmfund/horizon/errors"
	"gitlab.com/swarmfund/horizon/ledger"
	"gitlab.com/swarmfund/horizon/log"
	err2 "github.com/pkg/errors"
)

// Tick triggers the ingestion system to ingest any new ledger data, provided
// that there currently is not an import session in progress.
func (i *System) Tick() *Session {
	i.lock.Lock()
	if i.current != nil {
		log.Info("ingest: already in progress")
		i.lock.Unlock()
		return nil
	}

	is := i.newTickSession()
	i.current = is
	i.lock.Unlock()

	i.runOnce()
	return is
}

// newTickSession creates an unverified new ingestion session that reflects the
// current cached ledger state.
func (i *System) newTickSession() *Session {
	var (
		start int32
		ls    = ledger.CurrentState()
	)

	if ls.HistoryLatest == 0 {
		start = ls.CoreElder
	} else {
		start = ls.HistoryLatest + 1
	}

	end := ls.CoreLatest

	return NewSession(ls.HistoryElder > 1 || ls.CoreElder > 1, start, end, i)
}

// run causes the importer to check stellar-core to see if we can import new
// data.
func (i *System) runOnce() {
	defer func() {
		if rec := recover(); rec != nil {
			err := errors.FromPanic(rec)
			log.Errorf("import session panicked: %s", err)
		}
	}()

	ls := ledger.CurrentState()

	// 1. stash a copy of the current ingestion session (assigned from the tick)
	// 2. output "initial ingestion" message if the
	// 3. import until none available

	// 1.
	i.lock.Lock()
	is := i.current
	i.lock.Unlock()

	defer func() {
		i.lock.Lock()
		i.current = nil
		i.lock.Unlock()
	}()

	if is == nil {
		log.Warn("ingest: runOnce ran with a nil current session")
		return
	}

	if is.Cursor.FirstLedger > is.Cursor.LastLedger {
		return
	}

	// 2.
	if ls.HistoryLatest == 0 {
		log.Infof(
			"history db is empty, starting ingestion from ledger %d",
			is.Cursor.FirstLedger,
		)
	}

	if is.Cursor.FirstLedger != ls.CoreElder {
		err := i.validateLedgerChain(is.Cursor.FirstLedger)
		if err != nil {
			log.
				WithField("start", is.Cursor.FirstLedger).
				Errorf("ledger gap detected (possible db corruption): %s", err)
			return
		}
	}

	// 3.
	is.Run()

	if is.Err != nil {
		log.Errorf("import session failed: %s", is.Err)
	}

	return
}

// validateLedgerChain helps to ensure the chain of ledger entries is contiguous
// within horizon.  It ensures the ledger at `seq` is a child of `seq - 1`.
func (i *System) validateLedgerChain(seq int32) error {
	var (
		cur  core.LedgerHeader
		prev core.LedgerHeader
	)

	q := core.NewQ(i.CoreDB)

	err := q.LedgerHeaderBySequence(&cur, seq)
	if err != nil {
		return err2.Wrap(err, "validateLedgerChain: failed to load cur ledger")
	}

	err = q.LedgerHeaderBySequence(&prev, seq-1)
	if err != nil {
		return err2.Wrap(err, "validateLedgerChain: failed to load prev ledger")
	}

	if cur.PrevHash != prev.LedgerHash {
		return err2.New("cur and prev ledger hashes don't match")
	}

	return nil
}
