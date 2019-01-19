package pgdb_test

import (
	"io"
	"os"
	"testing"
	"time"

	"github.com/masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"gitlab.com/distributed_lab/kit/pgdb"
)

var prefix = time.Now().UnixNano()

func getDB(t *testing.T) *pgdb.DB {
	t.Helper()
	connURL := os.Getenv("PGDB_CONN_URL")
	if connURL == "" {
		t.Skip("skipping, PGDB_CONN_URL not set")
	}
	db, err := pgdb.Open(pgdb.Opts{
		URL:                connURL,
		MaxOpenConnections: 1,
	})
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	return db
}

// TODO test db open

func TestQueryer(t *testing.T) {
	db := getDB(t)

	t.Run("exec raw", func(t *testing.T) {
		err := db.ExecRaw("create temporary table q (i int)")
		assert.NoError(t, err)
	})

	t.Run("exec", func(t *testing.T) {
		stmt := squirrel.Insert("q").Values(1).Values(2)
		err := db.Exec(stmt)
		assert.NoError(t, err)
	})

	t.Run("select raw", func(t *testing.T) {
		type ts []struct {
			I int
		}
		var got ts
		expected := ts{
			{1}, {2},
		}
		err := db.SelectRaw(&got, "select * from q")
		assert.NoError(t, err)
		assert.Len(t, got, 2)
		assert.Equal(t, expected, got)
	})

	t.Run("select", func(t *testing.T) {
		type ts []struct {
			I int
		}
		var got ts
		expected := ts{
			{1}, {2},
		}
		stmt := squirrel.Select("*").From("q")
		err := db.Select(&got, stmt)
		assert.NoError(t, err)
		assert.Len(t, got, 2)
		assert.Equal(t, expected, got)
	})

	t.Run("get raw", func(t *testing.T) {
		var got struct {
			I int
		}
		err := db.GetRaw(&got, `select * from q order by i`)
		assert.NoError(t, err)
		assert.Equal(t, 1, got.I)
	})

	t.Run("get", func(t *testing.T) {
		var got struct {
			I int
		}
		stmt := squirrel.Select("*").From("q").OrderBy("i")
		err := db.Get(&got, stmt)
		assert.NoError(t, err)
		assert.Equal(t, 1, got.I)
	})
}

func TestDB_Transaction(t *testing.T) {
	db := getDB(t)

	if err := db.ExecRaw("create temporary table t (i int)"); err != nil {
		t.Fatalf("failed to prepare table: %v", err)
	}

	t.Run("successful commit", func(t *testing.T) {
		err := db.Transaction(func(q pgdb.Queryer) error {
			if err := q.ExecRaw(`insert into t values (1)`); err != nil {
				return err
			}

			if err := q.ExecRaw(`insert into t values (2)`); err != nil {
				return err
			}
			return nil
		})
		assert.NoError(t, err)

		// check committed

		var got int
		err = db.GetRaw(&got, `select count(1) from t`)
		assert.NoError(t, err)
		assert.Equal(t, 2, got)
	})

	t.Run("error", func(t *testing.T) {
		err := db.Transaction(func(q pgdb.Queryer) error {
			if err := q.ExecRaw(`insert into t values (1)`); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			return io.EOF
		})
		assert.Equal(t, io.EOF, errors.Cause(err))

		// check not committed

		var got int
		err = db.GetRaw(&got, `select count(1) from t`)
		assert.NoError(t, err)
		assert.Equal(t, 2, got)
	})

	t.Run("panic", func(t *testing.T) {
		assert.PanicsWithValue(t, io.EOF, func() {
			_ = db.Transaction(func(q pgdb.Queryer) error {
				if err := q.ExecRaw(`insert into t values (1)`); err != nil {
					t.Fatalf("unexpected error: %v", err)
				}

				panic(io.EOF)
			})
		})

		// check not committed

		var got int
		err := db.GetRaw(&got, `select count(1) from t`)
		assert.NoError(t, err)
		assert.Equal(t, 2, got)
	})

}
