package horizon

import (
	"gitlab.com/swarmfund/horizon/db2/core"
	"gitlab.com/swarmfund/horizon/render/hal"
	"gitlab.com/swarmfund/horizon/render/problem"
)

type CoreReferencesAction struct {
	Action

	accountID string
	reference string

	record core.Reference
}

func (action *CoreReferencesAction) JSON() {
	action.Do(
		action.loadParams,
		action.loadRecords,
		func() {
			response := map[string]interface{}{
				"data": action.record,
			}
			hal.Render(action.W, response)
		},
	)
}

func (action *CoreReferencesAction) loadParams() {
	action.accountID = action.GetNonEmptyString("account_id")
	action.reference = action.GetString("reference")
}

func (action *CoreReferencesAction) checkAllowed() {
	action.IsAllowed(action.accountID)
}

func (action *CoreReferencesAction) loadRecords() {
	q := action.CoreQ().References().ForAccount(action.accountID)

	if action.reference != "" {
		q = q.ByReference(action.reference)
	}

	records, err := q.Select()
	if err != nil {
		action.Log.WithError(err).Error("Failed to get References from core DB")
		action.Err = &problem.ServerError
		return
	}

	if len(records) == 0 {
		action.Log.Error("No records in DB matching request")
		action.Err = &problem.NotFound
		return
	}

	action.record = records[0]
}
