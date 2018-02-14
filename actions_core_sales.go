package horizon

import (
	"gitlab.com/swarmfund/horizon/db2/core"
	"gitlab.com/swarmfund/horizon/render/hal"
	"gitlab.com/swarmfund/horizon/render/problem"
)

// CoreSalesAction exists for non-obvious reasons and used by PSIM to close sales
type CoreSalesAction struct {
	Action
	Records []core.Sale
}

func (action *CoreSalesAction) JSON() {
	action.Do(
		action.loadRecords,
		action.checkAllowed,
		func() {
			hal.Render(action.W, action.Records)
		},
	)
}

func (action *CoreSalesAction) checkAllowed() {
	action.IsAllowed("")
}

func (action *CoreSalesAction) loadRecords() {
	records, err := action.CoreQ().Sales().Select()
	if err != nil {
		action.Log.WithError(err).Error("Failed to get sales from core DB")
		action.Err = &problem.ServerError
		return
	}
	action.Records = records
}
