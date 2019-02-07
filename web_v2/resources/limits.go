package resources

import (
	"github.com/spf13/cast"
	"gitlab.com/tokend/horizon/db2/core2"
	"gitlab.com/tokend/regources/v2"
)

func NewLimits(limits core2.LimitsV2) *regources.Limits {
	return &regources.Limits{
		Key: NewLimitsKey(limits.ID),
		Attributes: regources.LimitsAttr{
			AccountType:     limits.AccountType,
			AccountID:       limits.AccountId,
			StatsOpType:     limits.StatsOpType,
			AssetCode:       limits.AssetCode,
			IsConvertNeeded: limits.IsConvertNeeded,
			DailyOut:        limits.DailyOut,
			WeeklyOut:       limits.WeeklyOut,
			MonthlyOut:      limits.MonthlyOut,
			AnnualOut:       limits.AnnualOut,
		},
	}
}

func NewLimitsKey(limitsID uint64) regources.Key {
	return regources.Key{
		ID:   cast.ToString(limitsID),
		Type: regources.TypeLimits,
	}
}
