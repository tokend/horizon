package ingestion

import (
	sq "github.com/lann/squirrel"
	"gitlab.com/swarmfund/horizon/db2"
	"gitlab.com/swarmfund/horizon/db2/history"
)

// Ingestion receives write requests from a Session
type Ingestion struct {
	// DB is the sql repo to be used for writing any rows into the horizon
	// database.
	DB     *db2.Repo
	CoreDB *db2.Repo

	ledgers                  sq.InsertBuilder
	transactions             sq.InsertBuilder
	transaction_participants sq.InsertBuilder
	operations               sq.InsertBuilder
	operation_participants   sq.InsertBuilder
	recovery_requests        sq.InsertBuilder
	balances                 sq.InsertBuilder
	trades                   sq.InsertBuilder
	priceHistory             sq.InsertBuilder
	ledger_changes           sq.InsertBuilder
	contracts                sq.InsertBuilder
	contractsDetails         sq.InsertBuilder
	contractsDisputes        sq.InsertBuilder
}

func (i *Ingestion) HistoryQ() history.QInterface {
	return &history.Q{Repo: i.DB}
}
