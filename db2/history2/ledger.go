package history2

import (
	"gitlab.com/tokend/horizon/db2"
	"time"
)

// Ledger is a row of data from the `history_ledgers` table
type Ledger struct {
	db2.TotalOrderID
	Sequence     int32     `db:"sequence"`
	Hash         string    `db:"hash"`
	PreviousHash string    `db:"previous_hash"`
	ClosedAt     time.Time `db:"closed_at"`
}
