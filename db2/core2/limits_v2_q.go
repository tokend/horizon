package core2

import (
	sq "github.com/lann/squirrel"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2"
)

type LimitsV2Q struct {
	repo     *db2.Repo
	selector sq.SelectBuilder
}

// NewLimitsV2Q - default constructor for LimitsV2Q which
// creates LimitsV2Q with given db2.Repo and default selector
func NewLimitsV2Q(repo *db2.Repo) LimitsV2Q {
	return LimitsV2Q{
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
func (l2 LimitsV2Q) FilterByAccountID(accountID string) LimitsV2Q {
	l2.selector = l2.selector.Where("limits.account_id = ?", accountID)
	return l2
}

// Select - loads a rows from `limits_v2`
// returns nil, nil - if limits for particular account does not exists
func (l2 LimitsV2Q) Select() ([]LimitsV2, error) {
	var result []LimitsV2
	err := l2.repo.Select(&result, l2.selector)
	if err != nil {
		if l2.repo.NoRows(err) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to select limits")
	}

	return result, nil
}
