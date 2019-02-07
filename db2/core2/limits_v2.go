package core2

type LimitsV2 struct {
	ID              uint64  `db:"id"`
	AccountType     *int32  `db:"account_type"`
	AccountId       *string `db:"account_id"`
	StatsOpType     int32   `db:"stats_op_type"`
	AssetCode       string  `db:"asset_code"`
	IsConvertNeeded bool    `db:"is_convert_needed"`
	DailyOut        int64   `db:"daily_out"`
	WeeklyOut       int64   `db:"weekly_out"`
	MonthlyOut      int64   `db:"monthly_out"`
	AnnualOut       int64   `db:"annual_out"`
}
