package horizon

import (
	"bullioncoin.githost.io/development/go/xdr"
	"bullioncoin.githost.io/development/horizon/db2/core"
	"bullioncoin.githost.io/development/horizon/render/hal"
	"bullioncoin.githost.io/development/horizon/render/problem"
	"bullioncoin.githost.io/development/horizon/resource"
)

type AccountIndexAction struct {
	Action
	Types []xdr.AccountType

	Records []core.Account
	Page    hal.Page
}

func (action *AccountIndexAction) JSON() {
	action.Do(
		action.EnsureHistoryFreshness,
		action.checkAllowed,
		action.ValidateCursorWithinHistory,
		action.loadRecords,
		action.loadPage,
		func() {
			hal.Render(action.W, action.Page)
		},
	)
}

func (action *AccountIndexAction) checkAllowed() {
	action.IsAllowed("")
}

func (action *AccountIndexAction) loadRecords() {
	// pagination is turned off intentionally, coz we can't have string cursors atm
	err := action.CoreQ().Accounts().
		ForTypes(action.Types).
		Select(&action.Records)
	if err != nil {
		action.Log.WithError(err).Error("failed to load accounts")
		action.Err = &problem.ServerError
		return
	}
}

func (action *AccountIndexAction) loadPage() {
	for _, record := range action.Records {
		var r resource.Account
		err := r.Populate(
			action.Ctx,
			record,
			[]core.Signer{},
			[]core.Balance{},
			nil,
			nil,
			action.App.CoreInfo.DemurragePeriod,
		)
		if err != nil {
			action.Log.WithError(err).WithField("account", record.AccountID).
				Error("failed to populate resources")
			action.Err = &problem.ServerError
			return
		}
		action.Page.Add(r)
	}
	action.Page.BaseURL = action.BaseURL()
	action.Page.BasePath = action.Path()
}
