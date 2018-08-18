package horizon

import (
	"gitlab.com/swarmfund/horizon/render/hal"
	"gitlab.com/swarmfund/horizon/render/problem"
	"gitlab.com/tokend/go/xdr"
)

// ReviewableRequestIndexAction renders slice of reviewable requests
type ReviewableRequestCountAction struct {
	Action
	State        *int64
	RequestTypes []xdr.ReviewableRequestType
	Record       int64
}

// JSON is a method for actions.JSON
func (action *ReviewableRequestCountAction) JSON() {
	action.Do(
		action.EnsureHistoryFreshness,
		action.loadParams,
		action.checkAllowed,
		action.loadRecord,
		func() {
			hal.Render(action.W, action.Record)
		},
	)
}

func (action *ReviewableRequestCountAction) loadParams() {
	action.State = action.GetOptionalInt64("state")
}

func (action *ReviewableRequestCountAction) checkAllowed() {
	action.IsAllowed("")
}

func (action *ReviewableRequestCountAction) loadRecord() {
	q := action.HistoryQ().ReviewableRequests().CountQuery().ForTypes(action.RequestTypes)
	if action.State != nil {
		q = q.ForState(*action.State)
	}
	if action.Err != nil {
		return
	}

	var err error
	action.Record, err = q.Count()
	if err != nil {
		action.Log.WithError(err).Error("failed to load count of reviewable requests")
		action.Err = &problem.ServerError
		return
	}
}
