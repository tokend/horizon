package core2

import (
	sq "github.com/lann/squirrel"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2"
)

type LimitsQ struct {
	repo     *db2.Repo
	selector sq.SelectBuilder
}

// NewLimitsQ - default constructor for LimitsQ which
// creates LimitsQ with given db2.Repo and default Selector
func NewLimitsQ(repo *db2.Repo) LimitsQ {
	return LimitsQ{
		repo: repo,
		selector: sq.
			Select("limits.id, " +
				"limits.account_type, " +
				"limits.account_id, " +
				"limits.stats_op_type, " +
				"limits.asset_code," +
				"limits.is_convert_needed, " +
				"limits.daily_out, " +
				"limits.weekly_out, " +
				"limits.monthly_out, " +
				"limits.annual_out").
			From("limits_v2 limits"),
	}
}

// FilterByAccount - adds accountID filter for query to Limits table
func (q LimitsQ) FilterByAccountID(accountID string) LimitsQ {
	q.selector = q.selector.Where("limits.account_id = ?", accountID)
	return q
}

// Select - loads rows from `limits_v2`
// returns nil, nil - if limits for particular account does not exists
func (q LimitsQ) Select() ([]Limits, error) {
	var result []Limits
	err := q.repo.Select(&result, q.selector)
	if err != nil {
		if q.repo.NoRows(err) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to select limits")
	}

	return result, nil
}
