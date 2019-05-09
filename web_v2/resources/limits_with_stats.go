package resources

import (
	"gitlab.com/tokend/horizon/db2/core2"
	regources "gitlab.com/tokend/regources/generated"
)

func NewLimitsWithStats(record *core2.LimitsWithStatsEntry) *regources.LimitsWithStats {
	res := &regources.LimitsWithStats{
		Key: regources.Key{
			ID:   record.ID,
			Type: regources.LIMITS_WITH_STATS,
		},
		Relationships: regources.LimitsWithStatsRelationships{
			Account:    NewAccountKey(record.AccountID).AsRelation(),
			Limits:     NewLimitsKey(record.Limits.ID).AsRelation(),
			Statistics: NewStatisticKey(record.Statistics.Id).AsRelation(),
		},
	}

	return res
}
