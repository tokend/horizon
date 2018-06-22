package resource

import (
	"gitlab.com/swarmfund/horizon/db2/core"
)

type LimitsV2 struct {
	Id              uint64 		`json:"id"`
	AccountType     *int32    	`json:"account_type"`
	AccountId       *string     `json:"account_id"`
	StatsOpType     int32     	`json:"stats_op_type"`
	AssetCode       string      `json:"asset_code"`
	IsConvertNeeded bool        `json:"is_convert_needed"`
	DailyOut   		uint64 		`json:"daily_out"`
	WeeklyOut 		uint64 		`json:"weekly_out"`
	MonthlyOut 		uint64 		`json:"monthly_out"`
	AnnualOut  		uint64 		`json:"annual_out"`
}

// Populate fills out the fields of the signer, using one of an account's
// secondary signers.
func (s *LimitsV2) Populate(row core.LimitsV2Entry) {
	s.Id = row.Id
	s.AccountType = row.AccountType
	s.AccountId = row.AccountId
	s.StatsOpType = row.StatsOpType
	s.AssetCode = row.AssetCode
	s.IsConvertNeeded = row.IsConvertNeeded
	s.DailyOut = row.DailyOut
	s.WeeklyOut = row.WeeklyOut
	s.MonthlyOut = row.MonthlyOut
	s.AnnualOut = row.AnnualOut
}
