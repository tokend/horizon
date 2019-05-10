package statslimits

import (
	"fmt"
	"github.com/spf13/cast"
	"gitlab.com/tokend/horizon/db2/core2"
)

func ToCoreStatsLimitsUnitList(table Table) []core2.LimitsWithStats {
	res := make([]core2.LimitsWithStats, 0, len(table))

	for g, coreUnit := range table {
		var limitsID, statsID string
		if coreUnit.Limits.ID != 0 {
			limitsID = cast.ToString(coreUnit.Limits.ID)
		}
		if coreUnit.Stats.ID != 0 {
			statsID = cast.ToString(coreUnit.Stats.ID)
		}

		res = append(res, core2.LimitsWithStats{
			ID:          fmt.Sprintf("%s:%s", limitsID, statsID),
			AccountID:   coreUnit.Stats.AccountID,
			StatsOpType: g.StatsOpType,
			AssetCode:   g.AssetCode,
			Limits:      coreUnit.Limits,
			Statistics:  coreUnit.Stats,
		})

	}

	return res
}
