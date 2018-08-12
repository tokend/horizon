package horizon

import (
		"gitlab.com/swarmfund/horizon/render/hal"
	"gitlab.com/swarmfund/horizon/render/problem"
	"gitlab.com/tokend/regources"
	"gitlab.com/swarmfund/horizon/resource"
)

type ContractShowAction struct {
	Action
	ContractID     int64
	ContractRecord regources.Contract
}

// JSON is a method for actions.JSON
func (action *ContractShowAction) JSON() {
	action.Do(
		action.EnsureHistoryFreshness,
		action.checkAllowed,
		action.loadParams,
		action.loadRecords,
		func() {
			hal.Render(action.W, action.ContractRecord)
		},
	)
}

func (action *ContractShowAction) checkAllowed() {
	action.IsAllowed("")
}

func (action *ContractShowAction) loadParams() {
	action.ContractID = action.GetInt64("id")
}

func (action *ContractShowAction) loadRecords() {

	contract, err := action.HistoryQ().Contracts().ByID(action.ContractID)
	if err != nil {
		action.Log.WithError(err).Error("Failed to get contract record")
		action.Err = &problem.ServerError
		return
	}

	action.ContractRecord = resource.PopulateContract(contract)
}
