package resources

import (
	"github.com/spf13/cast"
	"gitlab.com/tokend/horizon/db2/core2"
	regources "gitlab.com/tokend/regources/generated"
)

func NewStatistics(stats core2.StatisticsEntry) *regources.Statistics {
	return &regources.Statistics{
		Key: NewStatisticKey(stats.Id),
		Attributes: regources.StatisticsAttributes{
			StatsOpType:     stats.StatsOpType,
			IsConvertNeeded: stats.IsConvertNeeded,
			DailyOut:        regources.Amount(stats.DailyOutcome),
			WeeklyOut:       regources.Amount(stats.WeeklyOutcome),
			MonthlyOut:      regources.Amount(stats.MonthlyOutcome),
			AnnualOut:       regources.Amount(stats.AnnualOutcome),
		},
		Relationships: regources.StatisticsRelationships{
			Account: NewAccountKey(stats.AccountId).AsRelation(),
			Asset:   NewAssetKey(stats.AssetCode).AsRelation(),
		},
	}
}

func NewStatisticKey(id uint64) regources.Key {
	return regources.Key{
		ID:   cast.ToString(id),
		Type: regources.STATISTICS,
	}
}
