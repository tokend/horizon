package history

import (
	"time"
	"gitlab.com/swarmfund/horizon/db2"
)

type Contract struct {
	db2.TotalOrderID
	Contractor    string                   `db:"contractor"`
	Customer      string                   `db:"customer"`
	Escrow        string                   `db:"escrow"`
	Disputer      string                   `db:"disputer"`
	StartTime     time.Time                `db:"start_time"`
	EndTime       time.Time                `db:"end_time"`
	Details       []map[string]interface{} `db:"details"`
	Invoices      []int64                  `db:"invoices"`
	DisputeReason map[string]interface{}   `db:"dispute_reason"`
	State         int32                    `db:"state"`
}
