package horizon

import (
	"gitlab.com/distributed_lab/tokend/horizon/db2/core"
	"gitlab.com/distributed_lab/tokend/horizon/render/hal"
	"gitlab.com/distributed_lab/tokend/horizon/render/problem"
	"gitlab.com/distributed_lab/tokend/horizon/resource"
	"math"
)

// This file contains the actions:
//
// AccountTypeLimitsShowAction: renders AccountTypeLimits for operationType
type AccountTypeLimitsShowAction struct {
	Action
	AccountType int32
	Limits      resource.Limits
}

// JSON is a method for actions.JSON
func (action *AccountTypeLimitsShowAction) JSON() {
	action.Do(
		action.loadParams,
		action.loadData,
		func() {
			hal.Render(action.W, action.Limits)
		},
	)
}
func (action *AccountTypeLimitsShowAction) loadParams() {
	action.AccountType = action.GetInt32("account_type")
}

func (action *AccountTypeLimitsShowAction) loadData() {
	result, err := action.CoreQ().LimitsByAccountType(action.AccountType)
	if err != nil {
		action.Log.WithError(err).Error("Failed to load limits by account type")
		action.Err = &problem.ServerError
		return
	}

	if result == nil {
		result = new(core.AccountTypeLimits)
		result.AccountType = action.AccountType
		result.DailyOut = math.MaxInt64
		result.WeeklyOut = math.MaxInt64
		result.MonthlyOut = math.MaxInt64
		result.AnnualOut = math.MaxInt64
	}

	action.Limits.Populate(result.Limits)
}
