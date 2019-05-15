package db2

import (
	"database/sql"
	"fmt"
	"reflect"
	"time"

	"github.com/jmoiron/sqlx"
	sq "github.com/lann/squirrel"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/log"
)

// Begin binds this repo to a new transaction.
func (r *Repo) Begin() error {
	if r.tx != nil {
		return errors.New("already in transaction")
	}

	tx, err := r.DB.Beginx()
	if err != nil {
		return errors.Wrap(err, "failed to begin tx")
	}
	r.logBegin()

	r.tx = tx
	return nil
}

// Clone clones the receiver, returning a new instance backed by the same
// context and db. The result will not be bound to any transaction that the
// source is currently within.
func (r *Repo) Clone() *Repo {
	return &Repo{
		DB: r.DB,
	}
}

// Commit commits the current transaction
func (r *Repo) Commit() error {
	if r.tx == nil {
		return errors.New("not in transaction")
	}

	err := r.tx.Commit()
	r.logCommit()
	r.tx = nil
	return err
}

func (r *Repo) DeleteRange(
	start, end int64,
	table string,
	idCol string,
) error {
	del := sq.Delete(table).Where(
		fmt.Sprintf("%s >= ? AND %s < ?", idCol, idCol),
		start,
		end,
	)
	_, err := r.Exec(del)
	return err
}

// Get runs `query`, setting the first result found on `dest`, if
// any.
func (r *Repo) Get(dest interface{}, query sq.Sqlizer) error {
	sql, args, err := r.build(query)
	if err != nil {
		return err
	}
	return r.GetRaw(dest, sql, args...)
}

// GetRaw runs `query` with `args`, setting the first result found on
// `dest`, if any.
func (r *Repo) GetRaw(dest interface{}, query string, args ...interface{}) error {
	query = r.conn().Rebind(query)
	start := time.Now()
	err := r.conn().Get(dest, query, args...)
	r.log("get", start, query, args)

	if err == nil {
		return nil
	}

	if r.NoRows(err) {
		return err
	}

	return errors.Wrap(err, "failed to exec query", logan.F{
		"query": query,
		"args":  args,
	})
}

// Exec runs `query`
func (r *Repo) Exec(query sq.Sqlizer) (sql.Result, error) {
	sql, args, err := r.build(query)
	if err != nil {
		return nil, err
	}
	return r.ExecRaw(sql, args...)
}

// ExecAll runs all sql commands in `script` against `r` within a single
// transaction.
func (r *Repo) ExecAll(script string) error {
	err := r.Begin()
	if err != nil {
		return err
	}

	defer r.Rollback()

	for _, cmd := range AllStatements(script) {
		_, err = r.ExecRaw(cmd)
		if err != nil {
			return err
		}
	}

	return r.Commit()
}

// ExecRaw runs `query` with `args`
func (r *Repo) ExecRaw(query string, args ...interface{}) (sql.Result, error) {
	query = r.conn().Rebind(query)
	start := time.Now()
	result, err := r.conn().Exec(query, args...)
	r.log("exec", start, query, args)

	if err == nil {
		return result, nil
	}

	if r.NoRows(err) {
		return nil, err
	}

	return nil, errors.Wrap(err, "failed to exec", logan.F{
		"query": query,
		"args":  args,
	})
}

// NoRows returns true if the provided error resulted from a query that found
// no results.
func (r *Repo) NoRows(err error) bool {
	return err == sql.ErrNoRows
}

// Query runs `query`, returns a *sqlx.Rows instance
func (r *Repo) Query(query sq.Sqlizer) (*sqlx.Rows, error) {
	sql, args, err := r.build(query)
	if err != nil {
		return nil, err
	}
	return r.QueryRaw(sql, args...)
}

// QueryRaw runs `query` with `args`
func (r *Repo) QueryRaw(query string, args ...interface{}) (*sqlx.Rows, error) {
	query = r.conn().Rebind(query)
	start := time.Now()
	result, err := r.conn().Queryx(query, args...)
	r.log("query", start, query, args)

	if err == nil {
		return result, nil
	}

	if r.NoRows(err) {
		return nil, err
	}

	return nil, errors.Wrap(err, "failed to query raw", logan.F{
		"query": query,
		"args":  args,
	})
}

// Rollback rolls back the current transaction
func (r *Repo) Rollback() error {
	if r.tx == nil {
		return errors.New("not in transaction")
	}

	err := r.tx.Rollback()
	r.logRollback()
	r.tx = nil
	return err
}

// Select runs `query`, setting the results found on `dest`.
func (r *Repo) Select(dest interface{}, query sq.Sqlizer) error {
	sql, args, err := r.build(query)
	if err != nil {
		return err
	}
	return r.SelectRaw(dest, sql, args...)
}

// SelectRaw runs `query` with `args`, setting the results found on `dest`.
func (r *Repo) SelectRaw(
	dest interface{},
	query string,
	args ...interface{},
) error {
	r.clearSliceIfPossible(dest)
	query = r.conn().Rebind(query)
	start := time.Now()
	err := r.conn().Select(dest, query, args...)
	r.log("select", start, query, args)

	if err == nil {
		return nil
	}

	if r.NoRows(err) {
		return err
	}

	return errors.Wrap(err, "failed to select raw", logan.F{
		"query": query,
		"args":  args,
	})
}

// build converts the provided sql builder `b` into the sql and args to execute
// against the raw database connections.
func (r *Repo) build(b sq.Sqlizer) (sql string, args []interface{}, err error) {
	sql, args, err = b.ToSql()

	if err != nil {
		err = errors.Wrap(err, "failed to build sql request")
	}
	return
}

// clearSliceIfPossible is a utility function that clears a slice if the
// provided interface wraps one. In the event that `dest` is not a pointer to a
// slice this func will fail with a warning, this allowing the forthcoming db
// select fail more concretely due to an incompatible destination.
func (r *Repo) clearSliceIfPossible(dest interface{}) {
	v := reflect.ValueOf(dest)
	vt := v.Type()

	if vt.Kind() != reflect.Ptr {
		log.Warn("cannot clear slice: dest is not pointer")
		return
	}

	if vt.Elem().Kind() != reflect.Slice {
		log.Warn("cannot clear slice: dest is a pointer, but not to a slice")
		return
	}

	reflect.Indirect(v).SetLen(0)
}

func (r *Repo) conn() Conn {
	if r.tx != nil {
		return r.tx
	}

	return r.DB
}

func (r *Repo) getLog() *logan.Entry {
	if r.Log != nil {
		return r.Log
	}

	return &log.DefaultLogger.Entry
}

func (r *Repo) log(typ string, start time.Time, query string, args []interface{}) {
	dur := time.Since(start)

	lEntry := r.getLog()
	fields := logan.F{
		"args": args,
		"sql":  query,
		"dur":  dur.String(),
	}
	lEntry.WithFields(fields).Debugf("sql: %s", typ)

	if dur > log.SlowQueryBound {
		lEntry.WithField("type", typ).WithFields(fields).Warn("too slow sql")
	}
}

func (r *Repo) logBegin() {
	r.getLog().Debug("sql: begin")
}

func (r *Repo) logCommit() {
	r.getLog().Debug("sql: commit")
}

func (r *Repo) logRollback() {
	r.getLog().Debug("sql: rollback")
}
