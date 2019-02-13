package resources

import (
	"github.com/spf13/cast"
	"gitlab.com/tokend/horizon/db2/core2"
	"gitlab.com/tokend/regources/v2"
)

// NewLimits creates new instance of Limits from provided one
func NewLimits(limits core2.Limits) *regources.Limits {
	return &regources.Limits{
		Key: regources.Key{
			ID:   cast.ToString(limits.ID),
			Type: regources.TypeLimits,
		},
		Attributes: regources.LimitsAttr{
			StatsOpType:     limits.StatsOpType,
			IsConvertNeeded: limits.IsConvertNeeded,
			DailyOut:        regources.Amount(limits.DailyOut),
			WeeklyOut:       regources.Amount(limits.WeeklyOut),
			MonthlyOut:      regources.Amount(limits.MonthlyOut),
			AnnualOut:       regources.Amount(limits.AnnualOut),
		},
	}
}
