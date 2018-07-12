package history

import (
	"gitlab.com/swarmfund/horizon/db2"
	sq "github.com/lann/squirrel"
)

var selectLedgerChanges = sq.Select(
	"hlc.tx_id",
	"hlc.op_id",
	"hlc.entry_index",
	"hlc.effect",
	"hlc.entry_type",
	"hlc.ledger_key",
).From("history_ledger_changes hlc")

// LedgerChangesQ is a helper struct to aid in configuring queries that loads
// slices of LedgerChanges structs.
type LedgerChangesQ struct {
	Err    error
	parent *Q
	sql    sq.SelectBuilder
}

type LedgerChangesQI interface {
	Page(page db2.PageQuery) LedgersQI
	Select(dest interface{}) error
}

// LedgerChanges provides a helper to filter rows from the `history_ledger_changes`
// table with pre-defined filters. See `LedgerChangesQ` methods for the available filters.
func (q *Q) LedgerChanges() LedgerChangesQI {
	return &LedgerChangesQ{
		parent: q,
		sql:    selectLedgerChanges,
	}
}

// Page specifies the paging constraints for the query being built by `q`.
func (q *LedgerChangesQ) Page(page db2.PageQuery) LedgersQI {
	if q.Err != nil {
		return q
	}

	q.sql, q.Err = page.ApplyTo(q.sql, "hlc.tx_id")
	return q
}

// Select loads the results of the query specified by `q` into `dest`.
func (q *LedgerChangesQ) Select(dest interface{}) error {
	if q.Err != nil {
		return q.Err
	}

	q.Err = q.parent.Select(dest, q.sql)
	return q.Err
}
