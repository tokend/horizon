package horizon

import (
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/swarmfund/horizon/render/hal"
	"gitlab.com/swarmfund/horizon/render/problem"
	"gitlab.com/swarmfund/horizon/resource"
)

// CoinsEmissionRequestShowAction returns a coins emission request based upon the provided
// id
type CoinsEmissionRequestShowAction struct {
	Action
	RequestID string
	Record    history.CoinsEmissionRequest
	Resource  resource.CoinsEmissionRequest
}

// JSON is a method for actions.JSON
func (action *CoinsEmissionRequestShowAction) JSON() {
	action.Do(
		action.loadParams,
		action.loadRecords,
		action.loadPage,
	)
	action.Do(func() {
		hal.Render(action.W, action.Resource)
	})
}

func (action *CoinsEmissionRequestShowAction) loadParams() {
	action.RequestID = action.GetString("id")
}

func (action *CoinsEmissionRequestShowAction) loadRecords() {
	action.Err = action.HistoryQ().CoinsEmissionRequestByRequestID(&action.Record, action.RequestID)
}

func (action *CoinsEmissionRequestShowAction) loadPage() {
	action.Resource.Populate(&action.Record)
}

type CheckPreEmissionAction struct {
	Action
	SerialNumber string
}

// JSON is a method for actions.JSON
func (action *CheckPreEmissionAction) JSON() {
	action.Do(
		action.loadParams,
		action.checkIsAllowed,
		action.loadRecords,
	)

	action.Do(func() {
		hal.Render(action.W, problem.Success)
	})
}

func (action *CheckPreEmissionAction) loadParams() {
	action.SerialNumber = action.GetString("serial_number")
}

func (action *CheckPreEmissionAction) checkIsAllowed() {
	action.isAllowed("")
}

func (action *CheckPreEmissionAction) loadRecords() {
	coinEmission, err := action.CoreQ().CoinsEmissions().BySerialNumber(action.SerialNumber)
	if err != nil {
		action.Log.WithError(err).Error("Failed to get coins emission from core_db")
		action.Err = &problem.ServerError
		return
	}

	if coinEmission != nil {
		action.Err = &problem.Conflict
		return
	}
}
