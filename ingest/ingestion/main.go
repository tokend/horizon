package ingestion

import (
	sq "github.com/lann/squirrel"
	"gitlab.com/swarmfund/horizon/db2"
	"gitlab.com/swarmfund/horizon/db2/core"
	"gitlab.com/swarmfund/horizon/db2/history"
)

// Ingestion receives write requests from a Session
type Ingestion struct {
	// DB is the sql repo to be used for writing any rows into the horizon
	// database.
	DB     *db2.Repo
	CoreDB *db2.Repo

	CoreQ    core.QInterface
	HistoryQ history.QInterface

	ledgers                  sq.InsertBuilder
	transactions             sq.InsertBuilder
	transaction_participants sq.InsertBuilder
	operations               sq.InsertBuilder
	operation_participants   sq.InsertBuilder
	recovery_requests        sq.InsertBuilder
	payment_requests         sq.InsertBuilder
	balances                 sq.InsertBuilder
	trades                   sq.InsertBuilder
	priceHistory             sq.InsertBuilder
}
