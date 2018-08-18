package horizon

import (
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/swarmfund/horizon/render/hal"
	"gitlab.com/swarmfund/horizon/render/problem"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/regources"
)

// ReviewableRequestIndexAction renders slice of reviewable requests
type ReviewableRequestCountAction struct {
	Action
	RequestTypes []xdr.ReviewableRequestType
	Record       regources.RequestsCount
}

// JSON is a method for actions.JSON
func (action *ReviewableRequestCountAction) JSON() {
	action.Do(
		action.EnsureHistoryFreshness,
		action.checkAllowed,
		action.loadRecord,
		func() {
			hal.Render(action.W, action.Record)
		},
	)
}

func (action *ReviewableRequestCountAction) checkAllowed() {
	action.IsAllowed("")
}

func (action *ReviewableRequestCountAction) loadRecord() {
	q := action.HistoryQ().ReviewableRequests().CountQuery().ForTypes(action.RequestTypes)
	if action.Err != nil {
		return
	}

	var err error
	action.Record.Approved, err = q.ForState(int64(history.ReviewableRequestStateApproved)).Count()
	if err != nil {
		action.Log.WithError(err).Error("failed to load count of approved reviewable requests")
		action.Err = &problem.ServerError
		return
	}

	action.Record.Pending, err = q.ForState(int64(history.ReviewableRequestStatePending)).Count()
	if err != nil {
		action.Log.WithError(err).Error("failed to load count of pending reviewable requests")
		action.Err = &problem.ServerError
		return
	}
}
