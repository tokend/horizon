package core2

type LimitsV2 struct {
	Id              uint64  `db:"id"`
	AccountType     *int32  `db:"account_type"`
	AccountId       *string `db:"account_id"`
	StatsOpType     int32   `db:"stats_op_type"`
	AssetCode       string  `db:"asset_code"`
	IsConvertNeeded bool    `db:"is_convert_needed"`
	DailyOut        uint64  `db:"daily_out"`
	WeeklyOut       uint64  `db:"weekly_out"`
	MonthlyOut      uint64  `db:"monthly_out"`
	AnnualOut       uint64  `db:"annual_out"`
}
