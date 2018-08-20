package horizon

import (
	"gitlab.com/swarmfund/horizon/db2"
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/swarmfund/horizon/render/hal"
	"gitlab.com/swarmfund/horizon/render/problem"
	"gitlab.com/swarmfund/horizon/resource/reviewablerequest"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/regources"
)

// ReviewableRequestIndexAction renders slice of reviewable requests
type ReviewableRequestIndexAction struct {
	Action
	CustomFilter       func(action *ReviewableRequestIndexAction)
	CustomCheckAllowed func(action *ReviewableRequestIndexAction)
	q                  history.ReviewableRequestQI
	Reviewer           string
	Requestor          string
	State              *int64
	UpdatedAfter       *int64
	Records            []history.ReviewableRequest
	Count              regources.RequestsCount

	RequestTypes []xdr.ReviewableRequestType

	PagingParams db2.PageQuery
	Page         hal.Page
}

// JSON is a method for actions.JSON
func (action *ReviewableRequestIndexAction) JSON() {
	action.Do(
		action.EnsureHistoryFreshness,
		action.loadParams,
		action.checkAllowed,
		action.loadRecord,
		action.loadPage,
		func() {
			hal.Render(action.W, action.Page)
		},
	)
}

func (action *ReviewableRequestIndexAction) loadParams() {
	action.PagingParams = action.GetPageQuery()
	action.Reviewer = action.GetString("reviewer")
	action.Requestor = action.GetString("requestor")
	action.State = action.GetOptionalInt64("state")
	action.UpdatedAfter = action.GetOptionalInt64("updated_after")
	action.Page.Filters = map[string]string{
		"reviewer":      action.Reviewer,
		"requestor":     action.Requestor,
		"state":         action.GetString("state"),
		"updated_after": action.GetString("updated_after"),
	}
}

func (action *ReviewableRequestIndexAction) checkAllowed() {
	if action.CustomCheckAllowed != nil {
		action.CustomCheckAllowed(action)
		return
	}
	action.IsAllowed(action.Requestor, action.Reviewer)
}

func (action *ReviewableRequestIndexAction) loadRecord() {
	action.q = action.HistoryQ().ReviewableRequests()

	if action.Reviewer != "" {
		action.q = action.q.ForReviewer(action.Reviewer)
	}

	if action.Requestor != "" {
		action.q = action.q.ForRequestor(action.Requestor)
	}

	if action.State != nil {
		action.q = action.q.ForState(*action.State)
	}

	if action.UpdatedAfter != nil {
		action.q = action.q.UpdatedAfter(*action.UpdatedAfter)
	}

	if action.CustomFilter != nil {
		action.CustomFilter(action)
	}

	if action.Err != nil {
		return
	}

	action.q = action.q.ForTypes(action.RequestTypes).Page(action.PagingParams)
	var err error
	action.Records, err = action.q.Select()
	if err != nil {
		action.Log.WithError(err).Error("failed to load reviewable requests")
		action.Err = &problem.ServerError
		return
	}

	q := action.HistoryQ().ReviewableRequests().CountQuery().ForTypes(action.RequestTypes)
	if action.Err != nil {
		return
	}

	action.Count.Approved, err = q.ForState(int64(history.ReviewableRequestStateApproved)).Count()
	if err != nil {
		action.Log.WithError(err).Error("failed to load count of approved reviewable requests")
		action.Err = &problem.ServerError
		return
	}

	action.Count.Pending, err = q.ForState(int64(history.ReviewableRequestStatePending)).Count()
	if err != nil {
		action.Log.WithError(err).Error("failed to load count of pending reviewable requests")
		action.Err = &problem.ServerError
		return
	}
}

func (action *ReviewableRequestIndexAction) loadPage() {
	for i := range action.Records {
		res, err := reviewablerequest.PopulateReviewableRequest(&action.Records[i])
		if err != nil {
			action.Log.WithError(err).Error("Failed to populate reviewable request")
			action.Err = &problem.ServerError
			return
		}
		action.Page.Add(res)
	}

	action.Page.Embedded.Meta = &hal.PageMeta{
		Count: &action.Count,
	}
	action.Page.BaseURL = action.BaseURL()
	action.Page.BasePath = action.Path()
	action.Page.Limit = action.PagingParams.Limit
	action.Page.Cursor = action.PagingParams.Cursor
	action.Page.Order = action.PagingParams.Order
	action.Page.PopulateLinks()
}
