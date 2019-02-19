package core2

import (
	sq "github.com/lann/squirrel"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2"
)

// AccountRoleQ is a helper struct to aid in configuring queries that loads
// accountRole structs.
type SignerRuleQ struct {
	repo     *db2.Repo
	selector sq.SelectBuilder
}

// NewAccountRoleQ - creates new instance of AccountRoleQ
func NewSignerRuleQ(repo *db2.Repo) SignerRuleQ {
	return SignerRuleQ{
		repo: repo,
		selector: sq.Select("sr.id",
			"sr.resource",
			"sr.action",
			"sr.is_forbid",
			"sr.is_default",
			"sr.owner_id",
			"sr.details",
		).From("signer_rules sr").Join("signer_role_rules srr on srr.rule_id = sr.id"),
	}
}

//FilterByRole - filter rules by role ID
func (q SignerRuleQ) FilterByRole(roleID uint64) SignerRuleQ {
	q.selector = q.selector.Where("srr.role_id = ?", roleID)
	return q
}

// Get - loads a row from `account_roles`
// returns nil, nil - if account does not exists
// returns error if more than one AccountRole found
func (q SignerRuleQ) Get() (*SignerRule, error) {
	var result SignerRule
	err := q.repo.Get(&result, q.selector)
	if err != nil {
		if q.repo.NoRows(err) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load signer rule")
	}

	return &result, nil
}

//Select - selects slice from the db, if no account rules found - returns nil, nil
func (q SignerRuleQ) Select() ([]SignerRule, error) {
	var result []SignerRule
	err := q.repo.Select(&result, q.selector)
	if err != nil {
		if q.repo.NoRows(err) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load signer rules")
	}

	return result, nil
}
