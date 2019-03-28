package resources

import (
	"github.com/spf13/cast"
	"gitlab.com/tokend/horizon/db2/core2"
	regources "gitlab.com/tokend/regources/v2/generated"
)

// NewLimits creates new instance of Limits from provided one
func NewLimits(limits core2.Limits) *regources.Limits {
	newLimits := &regources.Limits{
		Key: regources.Key{
			ID:   cast.ToString(limits.ID),
			Type: regources.LIMITS,
		},
		Attributes: regources.LimitsAttributes{
			StatsOpType:     limits.StatsOpType,
			IsConvertNeeded: limits.IsConvertNeeded,
			DailyOut:        regources.Amount(limits.DailyOut),
			WeeklyOut:       regources.Amount(limits.WeeklyOut),
			MonthlyOut:      regources.Amount(limits.MonthlyOut),
			AnnualOut:       regources.Amount(limits.AnnualOut),
		},
		Relationships: regources.LimitsRelationships{
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
