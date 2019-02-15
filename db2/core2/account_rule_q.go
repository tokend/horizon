package core2

import (
	sq "github.com/lann/squirrel"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2"
)

// AccountRoleQ is a helper struct to aid in configuring queries that loads
// accountRole structs.
type AccountRuleQ struct {
	repo     *db2.Repo
	selector sq.SelectBuilder
}

// NewAccountRoleQ - creates new instance of AccountRoleQ
func NewAccountRuleQ(repo *db2.Repo) AccountRuleQ {
	return AccountRuleQ{
		repo: repo,
		selector: sq.Select("ar.id",
			"ar.resource",
			"ar.action",
			"ar.is_forbid",
			"ar.details",
		).From("account_rules ar"),
	}
}

// GetByAddress loads a row from `accounts`, by address
// returns nil, nil - if account does not exists
func (q AccountRuleQ) GetByID(id uint64) (*AccountRule, error) {
	return q.FilterByIDs(id).Get()
}

//FilterByAddress - returns q with filter by address
func (q AccountRuleQ) FilterByIDs(ids ...uint64) AccountRuleQ {
	q.selector = q.selector.Where(sq.Eq{"ar.id": ids})
	return q
}

// Get - loads a row from `account_roles`
// returns nil, nil - if account does not exists
// returns error if more than one AccountRole found
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

		return nil, errors.Wrap(err, "failed to load account account rules")
	}

	return result, nil
}
