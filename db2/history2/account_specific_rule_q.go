package history2

import (
	sq "github.com/lann/squirrel"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2"
)

// AccountSpecificRulesQ is a helper struct to aid in configuring queries that loads accounts
type AccountSpecificRulesQ struct {
	repo     *db2.Repo
	selector sq.SelectBuilder
}

// NewAccountSpecificRulesQ - creates new instance of AccountSpecificRulesQ
func NewAccountSpecificRulesQ(repo *db2.Repo) AccountSpecificRulesQ {
	return AccountSpecificRulesQ{
		repo: repo,
		selector: sq.Select(
			"sr.id",
			"sr.address",
			"sr.entry_type",
			"sr.forbids",
			"sr.key",
		).From("account_specific_rules sr"),
	}
}

// ByAddress loads a row from `accounts`, by address
// returns nil, nil - if account does not exists
func (q AccountSpecificRulesQ) ByAddress(address string) AccountSpecificRulesQ {
	q.selector = q.selector.Where("sr.address = ?", address)
	return q
}

//Permission returns q with filter by `forbids`
func (q AccountSpecificRulesQ) Permission(forbids bool) AccountSpecificRulesQ {
	q.selector = q.selector.Where("sr.forbids = ?", forbids)
	return q
}

//ForSale - returns q with filter by sale id
func (q AccountSpecificRulesQ) ForSale(saleID uint64) AccountSpecificRulesQ {
	q.selector = q.selector.Where("sr.key#>>'{sale,saleID}' = ?", saleID)
	return q
}

//Get - selects account from db, returns nil, nil if one does not exists
func (q AccountSpecificRulesQ) Get() (*AccountSpecificRule, error) {
	var result AccountSpecificRule
	err := q.repo.Get(&result, q.selector)
	if err != nil {
		if q.repo.NoRows(err) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load account specific rule")
	}

	return &result, nil
}

// Page - returns Q with specified limit and offset params
func (q AccountSpecificRulesQ) Page(params db2.CursorPageParams) AccountSpecificRulesQ {
	q.selector = params.ApplyTo(q.selector, "sr.id")
	return q
}

//Get - selects account from db, returns nil, nil if one does not exists
func (q AccountSpecificRulesQ) Select() ([]AccountSpecificRule, error) {
	var result []AccountSpecificRule
	err := q.repo.Select(&result, q.selector)
	if err != nil {
		if q.repo.NoRows(err) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load account specific rules")
	}

	return result, nil
}
