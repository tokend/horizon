package core2

import (
	sq "github.com/masterminds/squirrel"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2"
)

type ExternalSystemIDsQ struct {
	repo     *db2.Repo
	selector sq.SelectBuilder
}

func NewExternalSystemIDsQ(repo *db2.Repo) ExternalSystemIDsQ {
	return ExternalSystemIDsQ{
		repo: repo,
		selector: sq.Select(
			"ext_pool.id, " +
				"ext.external_system_type, " +
				"ext.data, " +
				"ext.account_id, " +
				"ext_pool.is_deleted, " +
				"ext_pool.expires_at, " +
				"ext_pool.binded_at").
			From("external_system_account_id ext").
			LeftJoin("external_system_account_id_pool ext_pool ON " +
				"ext.account_id = ext_pool.account_id AND ext.external_system_type = ext_pool.external_system_type"),
	}
}

func (esid ExternalSystemIDsQ) FilterByAccount(accountID string) ExternalSystemIDsQ {
	esid.selector = esid.selector.Where("ext.account_id = ?", accountID)
	return esid
}

// Select - loads a rows from `external_system_account_id` left joining
// with table external_system_account_id_pool on account_id
// returns nil, nil - if external system IDs for particular account does not exists
func (l2 ExternalSystemIDsQ) Select() ([]ExternalSystemID, error) {
	var result []ExternalSystemID
	err := l2.repo.Select(&result, l2.selector)
	if err != nil {
		if l2.repo.NoRows(err) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to select external system IDs")
	}

	return result, nil
}
