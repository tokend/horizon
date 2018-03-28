package horizon

import (
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/horizon/db2"
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/swarmfund/horizon/render/hal"
	"gitlab.com/swarmfund/horizon/render/problem"
	"gitlab.com/swarmfund/horizon/resource/reviewablerequest"
)

// ReviewableRequestIndexAction renders slice of reviewable requests
type ReviewableRequestIndexAction struct {
	Action
	Reviewer     string
	Requestor    string
	State        *int64
	UpdatedAfter *int64
	Records      []history.ReviewableRequest

	RequestTypes []xdr.ReviewableRequestType
	RequestSpecificFilters map[string]string

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
	action.Page.Filters = map[string]string{}
	for key := range action.RequestSpecificFilters {
		action.RequestSpecificFilters[key] = action.GetString(key)
		action.Page.Filters[key] = action.RequestSpecificFilters[key]
	}

	action.Page.Filters["reviewer"] = action.Reviewer
	action.Page.Filters["requestor"] = action.Requestor
	action.Page.Filters["state"] = action.GetString("state")
}

func (action *ReviewableRequestIndexAction) checkAllowed() {
	action.IsAllowed(action.Requestor, action.Reviewer)
}

func (action *ReviewableRequestIndexAction) loadRecord() {
	q := action.HistoryQ().ReviewableRequests()

	if action.Reviewer != "" {
		q = q.ForReviewer(action.Reviewer)
	}

	if action.Requestor != "" {
		q = q.ForRequestor(action.Requestor)
	}

	if action.State != nil {
		q = q.ForState(*action.State)
	}

	if action.UpdatedAfter != nil {
		q = q.UpdatedAfter(*action.UpdatedAfter)
	}

	for key, value := range action.RequestSpecificFilters {
		if value == "" {
			continue
		}

		q = q.ByDetails(key, value)
	}

	q = q.ForTypes(action.RequestTypes).Page(action.PagingParams)
	var err error
	action.Records, err = q.Select()
	if err != nil {
		action.Log.WithError(err).Error("failed to load reviewable requests")
		action.Err = &problem.ServerError
		return
	}
}

func (action *ReviewableRequestIndexAction) loadPage() {
	for i := range action.Records {
		var res reviewablerequest.ReviewableRequest
		err := res.Populate(&action.Records[i])
		if err != nil {
			action.Log.WithError(err).Error("Failed to populate reviewable request")
			action.Err = &problem.ServerError
			return
		}
		action.Page.Add(&res)
	}

	action.Page.BaseURL = action.BaseURL()
	action.Page.BasePath = action.Path()
	action.Page.Limit = action.PagingParams.Limit
	action.Page.Cursor = action.PagingParams.Cursor
	action.Page.Order = action.PagingParams.Order
	action.Page.PopulateLinks()
}