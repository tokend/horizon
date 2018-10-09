package history

import (
	"time"

	"github.com/lib/pq"
	"gitlab.com/tokend/horizon/db2"
)

type Contract struct {
	db2.TotalOrderID
	Contractor      string        `db:"contractor"`
	Customer        string        `db:"customer"`
	Escrow          string        `db:"escrow"`
	StartTime       time.Time     `db:"start_time"`
	EndTime         time.Time     `db:"end_time"`
	InitialDetails  db2.Details   `db:"initial_details"`
	CustomerDetails db2.Details   `db:"customer_details"`
	Invoices        pq.Int64Array `db:"invoices"`
	State           int32         `db:"state"`
}

type ContractDetails struct {
	ContractID int64       `db:"contract_id"`
	Details    db2.Details `db:"details"`
	Author     string      `db:"author"`
	CreatedAt  time.Time   `db:"created_at"`
}

type ContractDispute struct {
	ContractID int64       `db:"contract_id"`
	Reason     db2.Details `db:"reason"`
	Author     string      `db:"author"`
	CreatedAt  time.Time   `db:"created_at"`
}
