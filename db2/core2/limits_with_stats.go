package core2

type LimitsWithStatsEntry struct {
	ID          string
	AccountID   string
	AssetCode   string
	StatsOpType int32
	Statistics  Statistics
	Limits      *Limits
}
