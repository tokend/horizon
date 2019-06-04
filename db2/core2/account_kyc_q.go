package core2

import (
	sq "github.com/lann/squirrel"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2"
)

// AccountsKycQ is a helper struct to aid in configuring queries that loads
// account structs.
type AccountsKycQ struct {
	repo     *db2.Repo
	selector sq.SelectBuilder
}

// NewAccountsKYCQ - creates new instance of AccountsKycQ
func NewAccountsKycQ(repo *db2.Repo) AccountsKycQ {
	return AccountsKycQ{
		repo: repo,
		selector: sq.Select("kyc.accountid",
			"kyc.kyc_data",
		).From("account_kyc kyc"),
	}
}

// GetByAddress loads a row from `accounts`, by address
// returns nil, nil - if account does not exists
func (q AccountsKycQ) GetByAddress(address string) (*AccountKYC, error) {
	return q.FilterByAddress(address).Get()
}

//FilterByAddress - returns q with filter by address
func (q AccountsKycQ) FilterByAddress(address string) AccountsKycQ {
	q.selector = q.selector.Where("kyc.accountid = ?", address)
	return q
}

// Get - loads a row from `accounts`
// returns nil, nil - if account does not exists
// returns error if more than one AccountKYC found
func (q AccountsKycQ) Get() (*AccountKYC, error) {
	var result AccountKYC
	err := q.repo.Get(&result, q.selector)
	if err != nil {
		if q.repo.NoRows(err) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load account")
	}

	return &result, nil
}
