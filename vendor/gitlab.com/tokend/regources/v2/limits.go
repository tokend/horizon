package regources

type Limits struct {
	Key
	Attributes LimitsAttr `json:"attributes"`
}

type LimitsAttr struct {
	AccountType     *int32  `json:"account_type"`
	AccountID       *string `json:"account_id"`
	StatsOpType     int32   `json:"stats_op_type"`
	AssetCode       string  `json:"asset_code"`
	IsConvertNeeded bool    `json:"is_convert_needed"`
	DailyOut        int64   `json:"daily_out"`
	WeeklyOut       int64   `json:"weekly_out"`
	MonthlyOut      int64   `json:"monthly_out"`
	AnnualOut       int64   `json:"annual_out"`
}
