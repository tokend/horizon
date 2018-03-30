package horizon

import (
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/horizon/db2"
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/swarmfund/horizon/render/hal"
	"gitlab.com/swarmfund/horizon/render/problem"
	"gitlab.com/swarmfund/horizon/resource/reviewablerequest"
	"strconv"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type reviewableRequestFilter struct {
	Fn fnReviewableRequestFilter
	Key string
}

type fnReviewableRequestFilter func(q history.ReviewableRequestQI, key, value string) (history.ReviewableRequestQI, error)

func reviewableRequestByEq(q history.ReviewableRequestQI, key, value string) (history.ReviewableRequestQI, error) {
	return q.ByDetailsEq(key, value), nil
}

func reviewableRequestMaskSet(q history.ReviewableRequestQI, key, rawValue string) (history.ReviewableRequestQI, error) {
	value, err := strconv.ParseInt(rawValue, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse int")
	}
	return q.ByDetailsMaskSet(key, value), nil
}

func reviewableRequestMaskNotSet(q history.ReviewableRequestQI, key, rawValue string) (history.ReviewableRequestQI, error) {
	value, err := strconv.ParseInt(rawValue, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse int")
	}
	return q.ByDetailsMaskNotSet(key, value), nil
}

// ReviewableRequestIndexAction renders slice of reviewable requests
type ReviewableRequestIndexAction struct {
	Action
	q history.ReviewableRequestQI
	Reviewer     string
	Requestor    string
	State        *int64
	UpdatedAfter *int64
	Records      []history.ReviewableRequest

	RequestTypes []xdr.ReviewableRequestType
	RequestSpecificFilters map[string]reviewableRequestFilter

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
	action.q = action.HistoryQ().ReviewableRequests()
	for requestParam, filter := range action.RequestSpecificFilters {
		value := action.GetString(requestParam)
		if value == "" {
			continue
		}

		var err error
		action.q, err = filter.Fn(action.q, filter.Key, value)
		if err != nil {
			action.SetInvalidField(requestParam, err)
			return
		}

		action.Page.Filters[requestParam] = value
	}

	action.Page.Filters["reviewer"] = action.Reviewer
	action.Page.Filters["requestor"] = action.Requestor
	action.Page.Filters["state"] = action.GetString("state")
}

func (action *ReviewableRequestIndexAction) checkAllowed() {
	action.IsAllowed(action.Requestor, action.Reviewer)
}

func (action *ReviewableRequestIndexAction) loadRecord() {

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

	action.q = action.q.ForTypes(action.RequestTypes).Page(action.PagingParams)
	var err error
	action.Records, err = action.q.Select()
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