package regources

import "gitlab.com/tokend/go/xdr"

//ManageLimitsCreation - details of corresponding op
type ManageLimits struct {
	Key
	Attributes ManageLimitsAttributes `json:"attributes"`
}

//ManageLimitsAttributes - details of the manage limits op
type ManageLimitsAttributes struct {
	Action xdr.ManageLimitsAction
	Create *ManageLimitsCreation `json:"create"`
	Remove *ManageLimitsRemoval  `json:"remove"`
}

//ManageLimitsCreation - details of corresponding op
type ManageLimitsCreation struct {
	AccountType     *xdr.AccountType `json:"account_type,omitempty"`
	AccountAddress  string           `json:"account_address,omitempty"`
	StatsOpType     xdr.StatsOpType  `json:"stats_op_type"`
	AssetCode       string           `json:"asset_code"`
	IsConvertNeeded bool             `json:"is_convert_needed"`
	DailyOut        Amount           `json:"daily_out"`
	WeeklyOut       Amount           `json:"weekly_out"`
	MonthlyOut      Amount           `json:"monthly_out"`
	AnnualOut       Amount           `json:"annual_out"`
}

//ManageLimitsRemoval - details of corresponding op
type ManageLimitsRemoval struct {
	LimitsID int64 `json:"limits_id"`
}
