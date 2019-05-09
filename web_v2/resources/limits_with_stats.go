package resources

import (
	"gitlab.com/tokend/horizon/db2/core2"
	regources "gitlab.com/tokend/regources/generated"
)

func NewLimitsWithStats(record *core2.LimitsWithStats) *regources.LimitsWithStats {
	res := &regources.LimitsWithStats{
		Key: regources.Key{
			ID:   record.ID,
			Type: regources.LIMITS_WITH_STATS,
		},
		Attributes: regources.LimitsWithStatsAttributes{
			Statistics: NewStatistics(record.Statistics),
		},
		Relationships: regources.LimitsWithStatsRelationships{
			Account: NewAccountKey(record.AccountID).AsRelation(),
		},
	}
	if record.Limits != nil {
		res.Attributes.Limits = NewLimits(*record.Limits)
	}

	return res
}
