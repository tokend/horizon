package statslimits

import (
	"gitlab.com/tokend/horizon/db2/core2"
	"math"
)

//Table is used to built complete limits & stats overview,
//using different level limits

var maxLimits = core2.Limits{
	DailyOut:   math.MaxUint64,
	WeeklyOut:  math.MaxUint64,
	MonthlyOut: math.MaxUint64,
	AnnualOut:  math.MaxUint64,
}

type (
	Group struct {
		AssetCode   string
		StatsOpType int32
	}
	LimitsWithStats struct {
		Limits core2.Limits
		Stats  core2.Statistics
	}
	Table map[Group]LimitsWithStats
)

func NewTable(limits []core2.Limits, stats []core2.Statistics) (lt Table) {
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

func (lt Table) Update(limits []core2.Limits) {
	for _, v := range limits {
		key := Group{
			AssetCode:   v.AssetCode,
			StatsOpType: v.StatsOpType,
		}

		entry, ok := lt[key]
		if !ok {
			entry = LimitsWithStats{
				Stats: core2.Statistics{},
			}
		}

		lt[key] = LimitsWithStats{
			Limits: v,
			Stats:  entry.Stats,
		}
	}
}

func (lt Table) FulfillEmptyLimits() {
	for k, v := range lt {
		if v.Limits.ID != 0 {
			continue
		}
		lt[k] = LimitsWithStats{
			Limits: maxLimits,
			Stats:  v.Stats,
		}
	}
}
