package resource

import (
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/horizon/db2/core"
)

type Limits struct {
	DailyOut   string `json:"daily_out"`
	WeeklyOut  string `json:"weekly_out"`
	MonthlyOut string `json:"monthly_out"`
	AnnualOut  string `json:"annual_out"`
}

// Populate fills out the fields of the signer, using one of an account's
// secondary signers.
func (s *Limits) Populate(row core.Limits) {
	s.DailyOut = amount.String(row.DailyOut)
	s.WeeklyOut = amount.String(row.WeeklyOut)
	s.MonthlyOut = amount.String(row.MonthlyOut)
	s.AnnualOut = amount.String(row.AnnualOut)
}
