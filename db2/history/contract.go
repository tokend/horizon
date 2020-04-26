package history

import (
	"gitlab.com/tokend/horizon/bridge"
	"time"

	"github.com/lib/pq"
)

type Contract struct {
	bridge.TotalOrderID
	Contractor      string         `db:"contractor"`
	Customer        string         `db:"customer"`
	Escrow          string         `db:"escrow"`
	StartTime       time.Time      `db:"start_time"`
	EndTime         time.Time      `db:"end_time"`
	InitialDetails  bridge.Details `db:"initial_details"`
	CustomerDetails bridge.Details `db:"customer_details"`
	Invoices        pq.Int64Array  `db:"invoices"`
	State           int32          `db:"state"`
}

type ContractDetails struct {
	ContractID int64          `db:"contract_id"`
	Details    bridge.Details `db:"details"`
	Author     string         `db:"author"`
	CreatedAt  time.Time      `db:"created_at"`
}

type ContractDispute struct {
	ContractID int64          `db:"contract_id"`
	Reason     bridge.Details `db:"reason"`
	Author     string         `db:"author"`
	CreatedAt  time.Time      `db:"created_at"`
}
