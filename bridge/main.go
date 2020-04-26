package bridge

import (
	"database/sql"
	"github.com/lann/squirrel"
	"gitlab.com/StepanTita/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3"
)

type TransactionFunc pgdb.TransactionFunc

// Struct which contains pgdb.DB instance, as well as other instances to substitute db2 fields
type Mediator struct {
	*pgdb.DB
	Log *logan.Entry
}

func (m Mediator) Clone() *Mediator {
	return &Mediator{
		DB: m.DB.Clone(),
	}
}

func Open(databaseURL string) (*Mediator, error) {
	repo, err := pgdb.Open(pgdb.Opts{
		URL:                databaseURL,
		MaxOpenConnections: 12,
		MaxIdleConnections: 4,
	})

	return &Mediator{
		DB: repo,
	}, err
}

func (m Mediator) Transaction(transactionFunc TransactionFunc) error {
	return m.DB.Transaction(pgdb.TransactionFunc(transactionFunc))
}

func (m Mediator) NoRows(err error) bool {
	return err == sql.ErrNoRows
}

func (m Mediator) Exec(query squirrel.Sqlizer) error {
	return m.DB.Exec(query)
}

func (m Mediator) Select(dest interface{}, query squirrel.Sqlizer) error {
	return m.DB.Select(dest, query)
}

func (m Mediator) GetRepo() *pgdb.DB {
	return m.DB
}

func (m Mediator) Get(dest interface{}, query squirrel.Sqlizer) error {
	return m.DB.Get(dest, query)
}

func (m *Mediator) ExecRawWithResult(query string, args ...interface{}) (sql.Result, error) {
	return m.DB.ExecRawWithResult(query, args...)
}
