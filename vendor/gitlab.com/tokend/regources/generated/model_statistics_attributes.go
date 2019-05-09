/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package regources

import (
	"time"
)

type StatisticsAttributes struct {
	AnnualOut       Amount    `json:"annual_out"`
	DailyOut        Amount    `json:"daily_out"`
	IsConvertNeeded bool      `json:"is_convert_needed"`
	MonthlyOut      Amount    `json:"monthly_out"`
	StatsOpType     int32     `json:"stats_op_type"`
	UpdatedAt       time.Time `json:"updated_at"`
	WeeklyOut       Amount    `json:"weekly_out"`
}
