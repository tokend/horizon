package statslimits

import (
	"gitlab.com/tokend/horizon/db2/core2"
)

//Table is used to built complete limits & stats overview,
//using different level limits

type (
	Group struct {
		AssetCode   string
		StatsOpType int32
	}
	LimitsWithStats struct {
		Limits core2.Limits
		Stats  core2.StatisticsEntry
	}
	Table map[Group]LimitsWithStats
)

func NewTable(limits []core2.Limits, stats []core2.StatisticsEntry) (lt Table) {
	lt = Table{}
	for _, entry := range limits {
		key := Group{
			AssetCode:   entry.AssetCode,
			StatsOpType: entry.StatsOpType,
		}

		limitsEntry := LimitsWithStats{
			Limits: entry,
		}

		lt[key] = limitsEntry
	}

	for _, entry := range stats {
		key := Group{
			AssetCode:   entry.AssetCode,
			StatsOpType: entry.StatsOpType,
		}

		limitsWithStatsEntry := LimitsWithStats{
			Limits: lt[key].Limits,
			Stats:  entry,
		}

		lt[key] = limitsWithStatsEntry
	}

	return lt
}

func (lt Table) UpdateLimits(limits []core2.Limits) {
	for _, v := range limits {
		key := Group{
			AssetCode:   v.AssetCode,
			StatsOpType: v.StatsOpType,
		}

		lt[key] = LimitsWithStats{
			Limits: v,
			Stats:  lt[key].Stats,
		}
	}
}
