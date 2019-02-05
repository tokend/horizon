package regources

type Limits struct {
	Key
	Attributes LimitsAttr `json:"attributes"`
}

type LimitsAttr struct {
	DailyOut   int64 `json:"daily_out"`
	WeeklyOut  int64 `json:"weekly_out"`
	MonthlyOut int64 `json:"monthly_out"`
	AnnualOut  int64 `json:"annual_out"`
}
