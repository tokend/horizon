package core

import (
	"database/sql"

	"gitlab.com/tokend/horizon/db2"

	sq "github.com/lann/squirrel"
	"gitlab.com/tokend/go/xdr"
)

var _ AccountQI = &AccountQ{}

type AccountQI interface {
	// returns nil, nil if account not found
	ByAddress(address string) (*Account, error)
	// filters by account type
	ForTypes(types []xdr.Uint64) AccountQI
	// performs select with specified filters
	Select(destination interface{}) error
	// filters by account ids
	ForAddresses(addresses ...string) AccountQI
	// filters by referrer
	ForReferrer(referrer string) AccountQI
	// Selects first element from filtered
	First() (*Account, error)
	// joins account KYC
	WithAccountKYC() AccountQI
	// applies paging params
	PageV2(page db2.PageQueryV2) AccountQI
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

func (q *AccountQ) ForTypes(types []xdr.Uint64) AccountQI {
	if q.Err != nil {
		return q
	}
	q.sql = q.sql.Where(sq.Eq{"account_type": types})
	return q
}

func (q *AccountQ) ForReferrer(referrer string) AccountQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("referrer = ?", referrer)
	return q
}

func (q *AccountQ) WithAccountKYC() AccountQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.
		LeftJoin("account_KYC ak on (ak.accountid = a.accountid)").
		Columns("ak.KYC_data as account_kyc_data")
	return q
}

func (q *AccountQ) PageV2(page db2.PageQueryV2) AccountQI {
	if q.Err != nil {
		return q
	}

	q.sql, q.Err = page.ApplyTo(q.sql)
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
	q.sql = q.sql.Where(sq.Eq{"a.accountid": addresses})
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
	"a.recoveryid",
	"a.thresholds",
	"a.account_type",
	"a.block_reasons",
	"a.referrer",
	"a.policies",
	"a.kyc_level",
).From("accounts a")
