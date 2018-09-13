package horizon

import (
	"gitlab.com/swarmfund/horizon/db2"
	"gitlab.com/swarmfund/horizon/db2/history"
	"gitlab.com/swarmfund/horizon/render/hal"
	"gitlab.com/swarmfund/horizon/render/problem"
	"gitlab.com/swarmfund/horizon/resource/reviewablerequest"
	"gitlab.com/tokend/go/doorman"
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
	HistRecords        []history.ReviewableRequest
	Records            []regources.ReviewableRequest

	RequestTypes []xdr.ReviewableRequestType

	PagingParams  db2.PageQuery
	Page          hal.Page
	DisablePaging bool
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
			hal.Render(action.W, action.Records)
		},
	)
}

func (action *ReviewableRequestIndexAction) loadParams() {
	action.DisablePaging = false
	action.PagingParams = action.GetPageQuery()
	if action.PagingParams.Cursor == "" {
		action.DisablePaging = true
	}
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

	constrains := []doorman.SignerConstraint{}
	extentions := []doorman.SignerExtension{}
	for _, actionType := range action.RequestTypes {
		if actionType == xdr.ReviewableRequestTypeIssuanceCreate {
			extentions = append(extentions, doorman.SignerExternsionPendingIssuance, doorman.SignerExternsionIssuanceHistory)
		}
		if actionType == xdr.ReviewableRequestTypeUpdateKyc {
			extentions = append(extentions, doorman.SignerExternsionPendingKYC, doorman.SignerExternsionKYCHistory)

		}
		if actionType == xdr.ReviewableRequestTypeSale {
			extentions = append(extentions, doorman.SignerExternsionCrowdfundingCampaign)
		}
	}

	for _, ext := range extentions {
		if action.Requestor != "" {
			constrains = append(constrains, doorman.SignerOfWithPermission(action.Requestor, ext))
		}
		if action.Reviewer != "" {
			constrains = append(constrains, doorman.SignerOfWithPermission(action.Reviewer, ext))
		}
		constrains = append(constrains, doorman.SignerOfWithPermission(action.App.CoreInfo.MasterAccountID, ext))
	}

	action.Doorman().Check(action.R, constrains...)
}

func (action *ReviewableRequestIndexAction) loadRecord() {
	if !action.DisablePaging {
		return
	}
	action.q = action.HistoryQ().ReviewableRequests()

	if action.CustomFilter != nil {
		action.CustomFilter(action)
	}

	if action.Err != nil {
		return
	}

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

	action.q = action.q.ForTypes(action.RequestTypes)

	var err error
	action.HistRecords, err = action.q.Select()
	if err != nil {
		action.Log.WithError(err).Error("failed to load reviewable requests")
		action.Err = &problem.ServerError
		return
	}
}

func (action *ReviewableRequestIndexAction) loadPage() {
	for i := range action.HistRecords {
		res, err := reviewablerequest.PopulateReviewableRequest(&action.HistRecords[i])
		if err != nil {
			action.Log.WithError(err).Error("Failed to populate reviewable request")
			action.Err = &problem.ServerError
			return
		}
		action.Records = append(action.Records, *res)
	}
}
