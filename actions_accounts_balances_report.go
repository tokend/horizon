package horizon

import (
	"gitlab.com/swarmfund/horizon/render/hal"
	"gitlab.com/swarmfund/horizon/render/problem"
)

type AccountsBalancesReportAction struct {
	Action

	threshold int64
}

func (action *AccountsBalancesReportAction) JSON() {
	action.Do(
		action.loadRecords,
		func() {
			hal.Render(action.W, action.Resource)
		},
	)
}

func (action *AccountsBalancesReportAction) loadRecords() {
	action.Err = &problem.ServerError
}
