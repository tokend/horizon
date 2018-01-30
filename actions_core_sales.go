package horizon

import (
	"gitlab.com/swarmfund/horizon/db2/core"
	"gitlab.com/swarmfund/horizon/render/problem"
	"gitlab.com/swarmfund/horizon/render/hal"
)

type CoreSalesAction struct {
	Action
	CoreRecords []core.Sale
}

func (action *CoreSalesAction) JSON() {
	action.Do(
		action.loadRecords,
		func() {
			hal.Render(action.W, action.CoreRecords)
		},
	)
}

func (action *CoreSalesAction) loadRecords() {
	q := action.CoreQ().Sales()
	err := q.Select(&action.CoreRecords)
	if err != nil {
		action.Log.WithError(err).Error("Failed to get sales from core DB")
		action.Err = &problem.ServerError
		return
	}
}
