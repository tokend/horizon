package ingestion

import (
	sq "github.com/lann/squirrel"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/tokend/horizon/db2/history"
)

// Ingestion receives write requests from a Session
type Ingestion struct {
	// DB is the sql repo to be used for writing any rows into the horizon
	// database.
	DB     *pgdb.DB
	CoreDB *pgdb.DB

	ledgers                  sq.InsertBuilder
	transactions             sq.InsertBuilder
	transaction_participants sq.InsertBuilder
	operations               sq.InsertBuilder
	operation_participants   sq.InsertBuilder
	trades                   sq.InsertBuilder
	priceHistory             sq.InsertBuilder
	ledger_changes           sq.InsertBuilder
	contracts                sq.InsertBuilder
	contractsDetails         sq.InsertBuilder
	contractsDisputes        sq.InsertBuilder
}

func (i *Ingestion) HistoryQ() history.QInterface {
	return &history.Q{DB: i.DB}
}
