package core

import (
	sq "github.com/lann/squirrel"
)

var _ AccountKYCQI = &AccountKYCQ{}

type AccountKYCQI interface {
	// returns nil, nil if account not found
	ByAddress(address string) (*AccountKYC, error)
	// performs select with specified filters
	Select(destination interface{}) error
}

type AccountKYCQ struct {
	Err    error
	parent *Q
	sql    sq.SelectBuilder
}

func (q *Q) AccountKYCs() AccountKYCQI {
	return &AccountKYCQ{
		parent: q,
		sql:    selectAccountKYC,
	}
}

func (q *AccountKYCQ) ByAddress(address string) (*AccountKYC, error) {
	result := new(AccountKYC)
	query := selectAccountKYC.Where("accountid = ?", address)
	err := q.parent.Get(result, query)
	if q.parent.NoRows(err) {
		return nil, nil
	}
	return result, err
}

func (q *AccountKYCQ) Select(destination interface{}) error {
	if q.Err != nil {
		return q.Err
	}
	return q.parent.Repo.Select(destination, q.sql)
}

var selectAccountKYC = sq.Select(
	"ak.accountid",
	"ak.KYC_data",
).From("account_KYC ak")
