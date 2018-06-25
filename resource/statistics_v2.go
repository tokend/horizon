package resource

import (
	"gitlab.com/swarmfund/horizon/db2/core"
)

type StatisticsV2 struct {
	Id              uint64 		`json:"id"`
	AccountId       string     `json:"account_id"`
	StatsOpType     int32     	`json:"stats_op_type"`
	AssetCode       string      `json:"asset_code"`
	IsConvertNeeded bool        `json:"is_convert_needed"`
	DailyOutcome   	uint64 		`json:"daily_outcome"`
	WeeklyOutcome  	uint64 		`json:"weekly_outcome"`
	MonthlyOutcome 	uint64 		`json:"monthly_outcome"`
	AnnualOutcome  	uint64 		`json:"annual_outcome"`
	UpdatedAt  		int64 		`json:"updated_at"`
}

// Populate fills out the fields of the signer, using one of an account's
// secondary signers.
func (s *StatisticsV2) Populate(row core.StatisticsV2Entry) {
	s.Id = row.Id
	s.AccountId = row.AccountId
	s.StatsOpType = row.StatsOpType
	s.AssetCode = row.AssetCode
	s.IsConvertNeeded = row.IsConvertNeeded
	s.DailyOutcome = row.DailyOutcome
	s.WeeklyOutcome = row.WeeklyOutcome
	s.MonthlyOutcome = row.MonthlyOutcome
	s.AnnualOutcome = row.AnnualOutcome
	s.UpdatedAt = row.UpdatedAt
}
