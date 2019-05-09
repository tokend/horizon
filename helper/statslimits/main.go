package statslimits

import (
	"github.com/spf13/cast"
	"gitlab.com/tokend/horizon/db2/core2"
)

func ToCoreStatsLimitsUnitList(table Table) []core2.LimitsWithStatsEntry {
	res := make([]core2.LimitsWithStatsEntry, 0, len(table))

	for g, unit := range table {
		res = append(res, core2.LimitsWithStatsEntry{
			ID: g.AssetCode + ":" +
				cast.ToString(unit.Limits.ID) + ":" +
				cast.ToString(unit.Stats.Id),
			AccountID:   unit.Stats.AccountId,
			StatsOpType: g.StatsOpType,
			AssetCode:   g.AssetCode,
			Limits:      unit.Limits,
			Statistics:  unit.Stats,
		})
	}

	return res
}
