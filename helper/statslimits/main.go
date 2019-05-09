package statslimits

import (
	"github.com/spf13/cast"
	"gitlab.com/tokend/horizon/db2/core2"
)

func ToCoreStatsLimitsUnitList(table Table) []core2.LimitsWithStatsEntry {
	res := make([]core2.LimitsWithStatsEntry, 0, len(table))

	for g, coreUnit := range table {
		var limitsID string

		if coreUnit.Limits != nil {
			limitsID = cast.ToString(coreUnit.Limits.ID)
		}

		res = append(res, core2.LimitsWithStatsEntry{
			ID: g.AssetCode + ":" + // todo rm before send to review
				limitsID + ":" +
				cast.ToString(coreUnit.Stats.ID),
			AccountID:   coreUnit.Stats.AccountID,
			StatsOpType: g.StatsOpType,
			AssetCode:   g.AssetCode,
			Limits:      coreUnit.Limits,
			Statistics:  coreUnit.Stats,
		})

	}

	return res
}
