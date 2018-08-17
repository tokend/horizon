package horizon

import (
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/swarmfund/horizon/render/hal"
	"gitlab.com/swarmfund/horizon/render/problem"
	"gitlab.com/swarmfund/horizon/resource/reviewablerequest"
)

// ReviewableRequestShowAction renders a reviewable request found by its ID.
type ReviewableRequestShowAction struct {
	Action
	RequestID   uint64
	RequestType *int64
	Record      *history.ReviewableRequest
}

// JSON is a method for actions.JSON
func (action *ReviewableRequestShowAction) JSON() {
	action.Do(
		action.EnsureHistoryFreshness,
		action.loadParams,
		action.loadRecord,
		action.checkAllowed,
		func() {
			res, err := reviewablerequest.PopulateReviewableRequest(action.Record)
			if err != nil {
				action.Log.WithError(err).Error("Failed to populate reviewable request")
				action.Err = &problem.ServerError
				return
			}
			if res != nil {
				hal.Render(action.W, *res)
			}
		},
	)
}

func (action *ReviewableRequestShowAction) loadParams() {
	action.RequestID = action.GetUInt64("id")
	action.RequestType = action.GetOptionalInt64("request_type")
}

func (action *ReviewableRequestShowAction) loadRecord() {
	var err error
	q := action.HistoryQ().ReviewableRequests()
	if action.RequestType != nil {
		q = q.ForType(*action.RequestType)
	}
	action.Record, err = q.ByID(action.RequestID)
	if err != nil {
		action.Log.WithError(err).WithField("request_id", action.RequestID).Error("failed to load reviewable request")
		action.Err = &problem.ServerError
		return
	}

	if action.Record == nil {
		action.Err = &problem.NotFound
		return
	}
}

func (action *ReviewableRequestShowAction) checkAllowed() {
	action.IsAllowed(action.Record.Requestor, action.Record.Reviewer)
}
