package horizon

import (
	"fmt"

	"gitlab.com/swarmfund/horizon/render/hal"
)

type AccountsBalancesReportAction struct {
	Action
	threshold int64
	Resource  struct {
		Lol string `json:"kek"`
	}
}

func (action *AccountsBalancesReportAction) JSON() {
	action.Do(
		action.setupThreshold,
		action.loadRecords,
		func() {
			hal.Render(action.W, &action.Resource)
		},
	)
}

func (action *AccountsBalancesReportAction) setupThreshold() {
	action.threshold = 100
}

func (action *AccountsBalancesReportAction) loadRecords() {
	accounts, err := action.CoreQ().Accounts().WithBalance().First()
	if err != nil {
		action.Log.WithError(err).Warn("failed to filter accounts with balance")
	}
	action.Log.Info("got accounts, first is ", fmt.Sprint(accounts))
}
