package resources

import (
	"github.com/spf13/cast"
	"gitlab.com/tokend/horizon/db2/core2"
	"gitlab.com/tokend/regources/rgenerated"
)

// NewLimits creates new instance of Limits from provided one
func NewLimits(limits core2.Limits) *rgenerated.Limits {
	newLimits := &rgenerated.Limits{
		Key: rgenerated.Key{
			ID:   cast.ToString(limits.ID),
			Type: rgenerated.LIMITS,
		},
		Attributes: rgenerated.LimitsAttributes{
			StatsOpType:     limits.StatsOpType,
			IsConvertNeeded: limits.IsConvertNeeded,
			DailyOut:        rgenerated.Amount(limits.DailyOut),
			WeeklyOut:       rgenerated.Amount(limits.WeeklyOut),
			MonthlyOut:      rgenerated.Amount(limits.MonthlyOut),
			AnnualOut:       rgenerated.Amount(limits.AnnualOut),
		},
		Relationships: rgenerated.LimitsRelationships{
			Asset: NewAssetKey(limits.AssetCode).AsRelation(),
		},
	}

	if limits.AccountId != nil {
		newLimits.Relationships.Account = NewAccountKey(*limits.AccountId).AsRelation()
	}

	if limits.AccountType != nil {
		newLimits.Relationships.AccountRole = NewAccountRoleKey(*limits.AccountType).AsRelation()
	}

	return newLimits
}
