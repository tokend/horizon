package test

import (
	"gitlab.com/tokend/horizon/db2"
	"gitlab.com/tokend/horizon/ledger"
)

// CoreRepo returns a db2.Repo instance pointing at the stellar core test database
func (t *T) CoreRepo() *db2.Repo {
	return &db2.Repo{
		DB:  t.CoreDB,
		Ctx: t.Ctx,
	}
}

// Finish finishes the test, logging any accumulated horizon logs to the logs
// output
func (t *T) Finish() {
	RestoreLogger()
	// Reset cached ledger state
	ledger.SetState(ledger.State{})

	if t.LogBuffer.Len() > 0 {
		t.T.Log("\n" + t.LogBuffer.String())
	}
}

// HistoryRepo returns a db2.Repo instance pointing at the horizon test database
func (t *T) HorizonRepo() *db2.Repo {
	return &db2.Repo{
		DB:  t.HorizonDB,
		Ctx: t.Ctx,
	}
}

// UpdateLedgerState updates the cached ledger state (or panicing on failure).
func (t *T) UpdateLedgerState() {
	var next ledger.State

	err := t.CoreRepo().GetRaw(&next, `
		SELECT
			COALESCE(MIN(ledgerseq), 0) as core_elder,
			COALESCE(MAX(ledgerseq), 0) as core_latest
		FROM ledgerheaders
	`)

	if err != nil {
		panic(err)
	}

	err = t.HorizonRepo().GetRaw(&next, `
			SELECT
				COALESCE(MIN(sequence), 0) as history_elder,
				COALESCE(MAX(sequence), 0) as history_latest
			FROM history_ledgers
		`)

	if err != nil {
		panic(err)
	}

	ledger.SetState(next)
	return
}
