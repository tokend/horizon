package resource

import (
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/horizon/db2/core"
)

type Statistics struct {
	DailyOutcome   string `json:"daily_outcome"`
	WeeklyOutcome  string `json:"weekly_outcome"`
	MonthlyOutcome string `json:"monthly_outcome"`
	AnnualOutcome  string `json:"annual_outcome"`
}

// Populate fills out the fields of the signer, using one of an account's
// secondary signers.
func (s *Statistics) Populate(row core.Statistics) {
	s.DailyOutcome = amount.String(row.DailyOutcome)
	s.WeeklyOutcome = amount.String(row.WeeklyOutcome)
	s.MonthlyOutcome = amount.String(row.MonthlyOutcome)
	s.AnnualOutcome = amount.String(row.AnnualOutcome)
}
