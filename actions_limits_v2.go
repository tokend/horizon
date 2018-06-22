package horizon

import (
	"gitlab.com/swarmfund/horizon/db2/core"
	"gitlab.com/swarmfund/horizon/render/hal"
	"gitlab.com/swarmfund/horizon/render/problem"
	"gitlab.com/swarmfund/horizon/resource"
)

// This file contains the actions:
//
// AccountTypeLimitsShowAction: renders AccountTypeLimits for operationType
type LimitsV2ShowAction struct {
	Action
	StatsOpType	int32
	AccountID 	*string
	AccountType	*int32
	LimitsV2	[]resource.LimitsV2
}

// JSON is a method for actions.JSON
func (action *LimitsV2ShowAction) JSON() {
	action.Do(
		action.loadParams,
		action.loadData,
		func() {
			hal.Render(action.W, action.LimitsV2)
		},
	)
}
func (action *LimitsV2ShowAction) loadParams() {
	action.StatsOpType = action.GetInt32("stats_op_type")
	action.AccountID = nil
	action.AccountType = nil
	accountID := action.GetString("account_id")
	if accountID != "" {
		action.AccountID = &accountID
	}
	accountType := action.GetInt32("account_type")
	if accountType != 0 {
		action.AccountType = &accountType
	}
}

func (action *LimitsV2ShowAction) loadData() {
	if (action.AccountID != nil) && (action.AccountType != nil){
		action.Log.Error("Unexpected state. Expected accountID or accountType, not both")
		action.Err = &problem.ServerError
		return
	}

	var result []core.LimitsV2Entry
	var err error

	if action.AccountID != nil {
		result, err = action.CoreQ().LimitsV2().ForAccountByStatsOpType(action.StatsOpType, *action.AccountID)
		if err != nil {
			action.Log.WithError(err).Error("Failed to load limits")
			action.Err = &problem.ServerError
			return
		}
	} else if action.AccountType != nil {
		result, err = action.CoreQ().LimitsV2().ForAccountTypeByStatsOpType(action.StatsOpType, *action.AccountType)
		if err != nil {
			action.Log.WithError(err).Error("Failed to load limits")
			action.Err = &problem.ServerError
			return
		}
	} else {
		result, err = action.CoreQ().LimitsV2().Select(action.StatsOpType)
		if err != nil {
			action.Log.WithError(err).Error("Failed to load limits")
			action.Err = &problem.ServerError
			return
		}
	}
	action.populateLimitsV2(result)
}

func (action *LimitsV2ShowAction) populateLimitsV2(limitsV2Records []core.LimitsV2Entry) {
	for i, limitsV2 := range limitsV2Records {
		action.LimitsV2 = append(action.LimitsV2, resource.LimitsV2{})
		action.LimitsV2[i].Populate(limitsV2)
	}
}