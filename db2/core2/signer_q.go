package core2

import (
	sq "github.com/lann/squirrel"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2"
)

//SignerQ - helper struct to load signers from db
type SignerQ struct {
	repo     *db2.Repo
	selector sq.SelectBuilder
}

//NewSignerQ - returns new instance of SignerQ with empty filter
func NewSignerQ(repo *db2.Repo) SignerQ {
	return SignerQ{
		repo: repo,
		selector: sq.Select("signers.accountid", "signers.publickey", "signers.weight", "signers.signer_type",
			"signers.identity_id", "signers.signer_name").From("signers"),
	}
}

//FilterByAccountAddress - return new instance of SignerQ with filter by account address
func (q SignerQ) FilterByAccountAddress(address string) SignerQ {
	q.selector = q.selector.Where("signers.accountid = ?", address)
	return q
}

//Select - selects slice of signers from db using specified filters. Returns nil, nil if none exists
func (q SignerQ) Select() ([]Signer, error) {
	var result []Signer
	err := q.repo.Select(&result, q.selector)
	if err != nil {
		if q.repo.NoRows(err) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to select signers")
	}

	return result, nil
}
