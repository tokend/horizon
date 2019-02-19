package core2

import (
	sq "github.com/lann/squirrel"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2"
)

// AccountRuleQ is a helper struct to aid in configuring queries that loads
// accountRule structs.
type AccountRuleQ struct {
	repo     *db2.Repo
	selector sq.SelectBuilder
}

// NewAccountRuleQ - creates new instance of AccountRuleQ
func NewAccountRuleQ(repo *db2.Repo) AccountRuleQ {
	return AccountRuleQ{
		repo: repo,
		selector: sq.Select("ar.id",
			"ar.resource",
			"ar.action",
			"ar.is_forbid",
			"ar.details",
		).From("account_rules ar").Join("account_role_rules arr on arr.rule_id = ar.id"),
	}
}

//FilterByRole - filter rules by role ID
func (q AccountRuleQ) FilterByRole(roleID uint64) AccountRuleQ {
	q.selector = q.selector.Where("arr.role_id = ?", roleID)
	return q
}

// Get - loads a row from `account_rules`
// returns nil, nil - if account does not exists
// returns error if more than one AccountRule found
func (q AccountRuleQ) Get() (*AccountRule, error) {
	var result AccountRule
	err := q.repo.Get(&result, q.selector)
	if err != nil {
		if q.repo.NoRows(err) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load account rule")
	}

	return &result, nil
}

//Select - selects slice from the db, if no account rules found - returns nil, nil
func (q AccountRuleQ) Select() ([]AccountRule, error) {
	var result []AccountRule
	err := q.repo.Select(&result, q.selector)
	if err != nil {
		if q.repo.NoRows(err) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load account rules")
	}

	return result, nil
}
