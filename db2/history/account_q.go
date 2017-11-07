package history

import (
	"time"

	"fmt"

	"bullioncoin.githost.io/development/horizon/db2"
	sq "github.com/lann/squirrel"
)

var selectAccount = sq.Select("ha.*").From("history_accounts ha")
var selectAccountsByOperation = sq.Select("ops.id as operation_id, acc.address").
	From("history_operations ops").
	LeftJoin("history_operation_participants opsparts ON opsparts.history_operation_id = ops.id").
	LeftJoin("history_accounts acc on acc.id = opsparts.history_account_id")

// AccountsQ is a helper struct to aid in configuring queries that loads
// slices of account structs.
type AccountsQ struct {
	Err    error
	parent *Q
	sql    sq.SelectBuilder
}

type AccountsQI interface {
	AccountsByOperationIDs(operationIDs []int64) (map[int64][]Account, error)
	Page(page db2.PageQuery) AccountsQI
	Select(dest interface{}) error
}

// Accounts provides a helper to filter rows from the `history_accounts` table
// with pre-defined filters.  See `AccountsQ` methods for the available filters.
func (q *Q) Accounts() AccountsQI {
	return &AccountsQ{
		parent: q,
		sql:    selectAccount,
	}
}

// AccountSummary will return all balance updates within bounds of `since` and `to`,
// including ones that does not change amount. There are couple of reasons for this:
// * idea of keeping ingest as simple/fast as possible and also avoiding of
// 	 abstraction levels mixing too much
// * filtering it out with SQL while reading is possible, but query will become
//	 unmaintainable before it even starts to work.
// If any additional filtering is required save your time and just do it inline
func (q *Q) AccountSummary(address string, since, to *time.Time) ([]BalanceSummary, error) {
	var result []BalanceSummary
	sincef := since.Format(time.RFC3339)
	tof := to.Format(time.RFC3339)

	amountBefore := fmt.Sprintf(`
	coalesce((
    	SELECT hbu.amount
		FROM history_balance_updates hbu
		WHERE hbu.balance_id = bid.balance_id
		  AND hbu.updated_at < '%s'::timestamp
		ORDER BY hbu.id DESC
		LIMIT 1
	), 0) AS amount_before`, sincef)

	updates := fmt.Sprintf(`
   (SELECT json_agg(hbu ORDER BY hbu.id)
   	FROM history_balance_updates hbu
    WHERE hbu.balance_id = bid.balance_id
      AND hbu.updated_at >= '%s'::timestamp
	  AND hbu.updated_at < '%s'::timestamp) AS updates`, sincef, tof)

	balances := fmt.Sprintf(`
	SELECT DISTINCT hbu.balance_id
    FROM history_balance_updates hbu
    JOIN history_balances hb ON hb.balance_id = hbu.balance_id
    WHERE hb.account_id = '%s'
	`, address)

	stmt := fmt.Sprintf(`
	SELECT bid.balance_id, %s, %s
	FROM (%s) AS bid`, amountBefore, updates, balances)

	err := q.SelectRaw(&result, stmt)
	return result, err
}

// AccountByAddress loads a row from `history_accounts`, by address
func (q *Q) AccountByAddress(dest interface{}, addy string) error {
	sql := selectAccount.Limit(1).Where("ha.address = ?", addy)
	return q.Get(dest, sql)
}

// AccountByID loads a row from `history_accounts`, by id
func (q *Q) AccountByID(dest interface{}, id int64) error {
	sql := selectAccount.Limit(1).Where("ha.id = ?", id)
	return q.Get(dest, sql)
}

// Page specifies the paging constraints for the query being built by `q`.
func (q *AccountsQ) Page(page db2.PageQuery) AccountsQI {
	if q.Err != nil {
		return q
	}

	q.sql, q.Err = page.ApplyTo(q.sql, "ha.id")
	return q
}

// Select loads the results of the query specified by `q` into `dest`.
func (q *AccountsQ) Select(dest interface{}) error {
	if q.Err != nil {
		return q.Err
	}

	q.Err = q.parent.Select(dest, q.sql)
	return q.Err
}

type OperationAccount struct {
	OperationID int64 `db:"operation_id"`
	Account
}

func (q *AccountsQ) AccountsByOperationIDs(operationIDs []int64) (map[int64][]Account, error) {
	if q.Err != nil {
		return nil, q.Err
	}

	if len(operationIDs) == 0 {
		return map[int64][]Account{}, nil
	}

	q.sql = selectAccountsByOperation.Where(sq.Eq{"ops.id": operationIDs})
	var operationUsers []OperationAccount
	err := q.parent.Select(&operationUsers, q.sql)
	if err != nil {
		return nil, err
	}

	result := make(map[int64][]Account)
	for i := range operationUsers {
		operationAccount := operationUsers[i]
		if result[operationAccount.OperationID] == nil {
			result[operationAccount.OperationID] = make([]Account, 0, 1)
		}

		result[operationAccount.OperationID] = append(result[operationAccount.OperationID], operationAccount.Account)
	}

	return result, nil
}
