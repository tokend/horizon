package core

import (
	sq "github.com/lann/squirrel"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

var selectAccountKYC = sq.Select(
	"ak.kyc_data as account_kyc_data",
).From("account_kyc ak")

type AccountKYCQI interface {
	ByAddress(address string) (*AccountKYC, error)
}

type AccountKYCQ struct {
	Err    error
	parent *Q
	sql    sq.SelectBuilder
}

func (q *Q) AccountKYC() AccountKYCQI {
	return &AccountKYCQ{
		parent: q,
		sql:    selectAccountKYC,
	}
}

func (q *AccountKYCQ) ByAddress(accountID string) (*AccountKYC, error) {
	sql := selectAccountKYC.Where("ak.accountid = ?", accountID)

	var result AccountKYC

	err := q.parent.Get(&result, sql)

	if q.parent.NoRows(err) {
		return nil, nil
	}

	if err != nil {
		return nil, errors.Wrap(err, "failed to load account_kyc", map[string]interface{}{
			"account_id": accountID,
		})
	}

	return &result, err
}
