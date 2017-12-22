package horizon

import (
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/horizon/db2"
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/swarmfund/horizon/render/hal"
	"gitlab.com/swarmfund/horizon/render/problem"
	"gitlab.com/swarmfund/horizon/resource/reviewablerequest"
)

// WithdrawalIndexAction renders slice of reviewable requests
type WithdrawalIndexAction struct {
	Action
	DestAsset    string
	Requester    string
	State        *int64
	Records      []history.ReviewableRequest
	PagingParams db2.PageQuery
	Page         hal.Page
}

// JSON is a method for actions.JSON
func (action *WithdrawalIndexAction) JSON() {
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

func (action *WithdrawalIndexAction) loadParams() {
	action.PagingParams = action.GetPageQuery()
	action.DestAsset = action.GetString("to_asset")
	action.State = action.GetOptionalInt64("state")
	action.Requester = action.GetString("requester")
	action.Page.Filters = map[string]string{
		"to_asset":  action.DestAsset,
		"state":     action.GetString("state"),
		"requester": action.Requester,
	}
}

func (action *WithdrawalIndexAction) checkAllowed() {
	action.IsAllowed("")
}

func (action *WithdrawalIndexAction) loadRecord() {
	q := action.HistoryQ().
		ReviewableRequests().
		ForType(int64(xdr.ReviewableRequestTypeWithdraw))

	if action.DestAsset != "" {
		q = q.ForDestAsset(action.DestAsset)
	}

	if action.Requester != "" {
		q = q.ForRequestor(action.Requester)
	}

	if action.State != nil {
		q = q.ForState(*action.State)
	}

	q = q.Page(action.PagingParams)
	var err error
	action.Records, err = q.Select()
	if err != nil {
		action.Log.WithError(err).Error("failed to load reviewable requests")
		action.Err = &problem.ServerError
		return
	}
}

func (action *WithdrawalIndexAction) loadPage() {
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
