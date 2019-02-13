package core2

import (
	sq "github.com/masterminds/squirrel"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2"
)

type ExternalSystemIDsQ struct {
	repo     *db2.Repo
	Selector sq.SelectBuilder
}

// NewExternalSystemIDsQ - default constructor for ExternalSystemIDsQ which
// creates ExternalSystemIDsQ with given db2.Repo and default Selector
func NewExternalSystemIDsQ(repo *db2.Repo) ExternalSystemIDsQ {
	return ExternalSystemIDsQ{
		repo: repo,
		Selector: sq.Select(
			"ext_pool.id",
			"ext_pool.external_system_type",
			"ext_pool.data",
			"ext_pool.account_id",
			"ext_pool.is_deleted",
			"to_timestamp(ext_pool.expires_at) at time zone 'UTC' expires_at",
			"to_timestamp(ext_pool.binded_at) at time zone 'UTC' binded_at").
			From("external_system_account_id_pool ext_pool"),
	}
}

// FilterByAccount - adds accountID filter for query to external system IDs table
func (esid ExternalSystemIDsQ) FilterByAccount(accountID string) ExternalSystemIDsQ {
	esid.Selector = esid.Selector.Where("ext_pool.account_id = ?", accountID)
	return esid
}

// Select - loads a rows from `external_system_account_id` left joining
// with table external_system_account_id_pool on account_id
// returns nil, nil - if external system IDs for particular account does not exists
func (esid ExternalSystemIDsQ) Select() ([]ExternalSystemID, error) {
	var result []ExternalSystemID
	err := esid.repo.Select(&result, esid.Selector)
	if err != nil {
		if esid.repo.NoRows(err) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to select external system IDs")
	}

	return result, nil
}
