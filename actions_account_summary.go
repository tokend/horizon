package horizon

import (
	"time"

	"errors"

	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/swarmfund/horizon/render/hal"
	"gitlab.com/swarmfund/horizon/render/problem"
	"gitlab.com/swarmfund/horizon/resource"
)

type AccountSummaryAction struct {
	Action

	AccountID string
	Since     *time.Time
	To        *time.Time

	Records  []history.BalanceSummary
	Resource resource.AccountSummary
}

func (action *AccountSummaryAction) JSON() {
	action.Do(
		action.loadParams,
		action.loadRecords,
		action.loadResource,
		func() {
			hal.Render(action.W, action.Resource)
		},
	)
}

func (action *AccountSummaryAction) loadParams() {
	action.AccountID = action.GetNonEmptyString("id")

	action.Since = action.GetTime("since")
	if action.Since != nil && action.Since.After(time.Now()) {
		action.SetInvalidField("since", errors.New("time travel is forbidden"))
		return
	}

	action.To = action.TryGetTime("to")
	if action.Err == nil && action.To == nil {
		now := time.Now()
		action.To = &now
	}
}

func (action *AccountSummaryAction) loadRecords() {
	summary, err := action.HistoryQ().AccountSummary(action.AccountID, action.Since, action.To)
	if err != nil {
		action.Log.WithError(err).Error("failed to get account summary")
		action.Err = &problem.ServerError
		return
	}
	action.Records = summary
}

func (action *AccountSummaryAction) loadResource() {
	action.Resource.Populate(action.Records)
}
