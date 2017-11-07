package horizon

import (
	"gitlab.com/distributed_lab/tokend/horizon/db2/core"
	"gitlab.com/distributed_lab/tokend/horizon/render/hal"
	"gitlab.com/distributed_lab/tokend/horizon/render/problem"
	"gitlab.com/distributed_lab/tokend/horizon/resource"
)

type BalanceTrustsAction struct {
	Action
	BalanceID   string
	CoreRecords []core.Trust
	Resource    resource.Trusts
}

// JSON is a method for actions.JSON
func (action *BalanceTrustsAction) JSON() {
	action.Do(
		action.loadParams,
		action.checkAllowed,
		action.loadRecords,
		func() {
			hal.Render(action.W, action.Resource)
		},
	)
}

func (action *BalanceTrustsAction) loadParams() {
	action.BalanceID = action.GetString("balance_id")
}

func (action *BalanceTrustsAction) checkAllowed() {
	action.IsAllowed("")
}

func (action *BalanceTrustsAction) loadRecords() {
	var err error
	err = action.CoreQ().Trusts().ForBalance(action.BalanceID).Select(&action.CoreRecords)

	if err != nil {
		action.Log.WithError(err).Error("Failed to get trusts from core DB")
		action.Err = &problem.ServerError
		return
	}

	action.Resource.Populate(action.CoreRecords)

}
