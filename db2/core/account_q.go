package core

import (
	"database/sql"

	"bullioncoin.githost.io/development/go/xdr"
	sq "github.com/lann/squirrel"
)

var _ AccountQI = &AccountQ{}

type AccountQI interface {
	// returns nil, nil if account not found
	ByAddress(address string) (*Account, error)
	// filters by account type
	ForTypes(types []xdr.AccountType) AccountQI
	// joins statistics
	WithStatistics() AccountQI
	// filters by referrer
	ForReferrer(referrer string) AccountQI
	// performs select with specified filters
	Select(destination interface{}) error
	// filters by account ids
	ForAddresses(addresses ...string) AccountQI
	// Selects first element from filtered
	First() (*Account, error)
}

// AccountQ is a helper struct to aid in configuring queries that loads
// slices or entry of Account structs.
type AccountQ struct {
	Err    error
	parent *Q
	sql    sq.SelectBuilder
}

func (q *Q) Accounts() AccountQI {
	return &AccountQ{
		parent: q,
		sql:    selectAccount,
	}
}

func (q *AccountQ) ByAddress(address string) (*Account, error) {
	result := new(Account)
	query := selectAccount.Limit(1).Where("accountid = ?", address)
	err := q.parent.Get(result, query)
	if q.parent.NoRows(err) {
		return nil, nil
	}
	return result, err
}

func (q *AccountQ) ForTypes(types []xdr.AccountType) AccountQI {
	if q.Err != nil {
		return q
	}

	for _, t := range types {
		if t == xdr.AccountTypeExchange {
			// master is our secret exchange
			types = append(types, xdr.AccountTypeMaster)
			break
		}
	}
	q.sql = q.sql.Where(sq.Eq{"account_type": types})
	return q
}

func (q *AccountQ) WithStatistics() AccountQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.
		LeftJoin("statistics st on (st.account_id = a.accountid)").
		Columns(
			"st.daily_out as st_daily_out",
			"st.weekly_out as st_weekly_out",
			"st.monthly_out as st_monthly_out",
			"st.annual_out as st_annual_out",
			"st.updated_at as st_updated_at")
	return q
}

func (q *AccountQ) ForReferrer(referrer string) AccountQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("referrer = ?", referrer)
	return q
}

func (q *AccountQ) Select(destination interface{}) error {
	if q.Err != nil {
		return q.Err
	}
	return q.parent.Repo.Select(destination, q.sql)
}

func (q *AccountQ) ForAddresses(addresses ...string) AccountQI {
	if q.Err != nil {
		return q
	}
	q.sql = q.sql.Where(sq.Eq{"accountid": addresses})
	return q
}

func (q *AccountQ) First() (*Account, error) {
	if q.Err != nil {
		return nil, q.Err
	}
	var result Account
	err := q.parent.Repo.Get(&result, q.sql)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &result, err
}

var selectAccount = sq.Select(
	"a.accountid",
	"a.thresholds",
	"a.account_type",
	"a.block_reasons",
	"a.referrer",
	"a.share_for_referrer",
	"ex.name",
	"ex.require_review",
	"a.policies",
	"a.created_at",
).From("accounts a").
	LeftJoin("exchanges ex ON a.accountid=ex.account_id")
