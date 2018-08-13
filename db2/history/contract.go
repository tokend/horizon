package history

import (
	"time"

	"github.com/lib/pq"
	"gitlab.com/swarmfund/horizon/db2"
)

type Contract struct {
	db2.TotalOrderID
	Contractor    string           `db:"contractor"`
	Customer      string           `db:"customer"`
	Escrow        string           `db:"escrow"`
	Disputer      string           `db:"disputer"`
	StartTime     time.Time        `db:"start_time"`
	EndTime       time.Time        `db:"end_time"`
	Details       db2.DetailsArray `db:"details"`
	Invoices      pq.Int64Array    `db:"invoices"`
	DisputeReason db2.Details      `db:"dispute_reason"`
	State         int32            `db:"state"`
}
