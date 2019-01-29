package pgdb

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type DB struct {
	Queryer
	db *sqlx.DB
}

// Transaction is generic helper method for specific Q's to implement Transaction capabilities
func (db *DB) Transaction(fn TransactionFunc) (err error) {
	tx, err := db.db.BeginTxx(context.TODO(), nil)
	if err != nil {
		return errors.Wrap(err, "failed to begin tx")
	}
	// swallowing rollback err, should not affect data consistency
	defer tx.Rollback()

	if err = fn(newQueryer(tx)); err != nil {
		return errors.Wrap(err, "failed to execute statements")
	}

	if err = tx.Commit(); err != nil {
		return errors.Wrap(err, "failed to commit tx")
	}

	return
}
