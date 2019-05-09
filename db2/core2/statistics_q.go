package core2

import (
	"database/sql"
	"gitlab.com/tokend/horizon/db2"
	"time"

	sq "github.com/lann/squirrel"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

var selectStatistics = sq.Select("st.id",
	"st.account_id",
	"st.stats_op_type",
	"st.asset_code",
	"st.is_convert_needed",
	"st.daily_out",
	"st.weekly_out",
	"st.monthly_out",
	"st.annual_out",
	"st.updated_at").From("statistics_v2 st")

type StatisticsQI interface {
	FilterByAccount(accountID string) ([]StatisticsEntry, error)
}

type StatisticsQ struct {
	repo     *db2.Repo
	selector sq.SelectBuilder
}

func (q *StatisticsQ) FilterByAccount(accountID string) ([]StatisticsEntry, error) {
	query := selectStatistics.Where("st.account_id = ?", accountID)
	var result []StatisticsEntry
	err := q.repo.Select(&result, query)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, errors.Wrap(err, "Failed to load statistics_v2 for account")
	}

	for _, statisticsV2 := range result {
		statisticsV2.ClearObsolete(time.Now().UTC())
	}

	return result, nil
}
