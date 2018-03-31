package horizon

import (
	"gitlab.com/swarmfund/horizon/db2/core"
	"gitlab.com/swarmfund/horizon/render/hal"
	"gitlab.com/swarmfund/horizon/render/problem"
)

type CoreReferencesAction struct {
	Action

	accountID string
	records []core.Reference
}

func (action *CoreReferencesAction) JSON() {
	action.Do(
		action.loadAccountID,
		action.checkAllowed,
		action.loadRecords,
		func() {
			response := map[string]interface{} {
				"data": action.records,
			}
			hal.Render(action.W, response)
		},
	)
}

func (action *CoreReferencesAction) loadAccountID() {
	action.accountID = action.GetNonEmptyString("account_id")
}

func (action *CoreReferencesAction) checkAllowed() {
	action.IsAllowed(action.accountID)
}

func (action *CoreReferencesAction) loadRecords() {
	records, err := action.CoreQ().References().ForAccount(action.accountID).Select()
	if err != nil {
		action.Log.WithError(err).Error("Failed to get References from core DB")
		action.Err = &problem.ServerError
		return
	}

	action.records = records
}
