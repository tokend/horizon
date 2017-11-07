package horizon

import (
	"gitlab.com/distributed_lab/tokend/horizon/db2/core"
	"gitlab.com/distributed_lab/tokend/horizon/render/hal"
	"gitlab.com/distributed_lab/tokend/horizon/render/problem"
	"gitlab.com/distributed_lab/tokend/horizon/resource"
)

type AccountReferralsAction struct {
	Action

	AccountID string

	Records  []core.Account
	Balances map[string][]core.Balance
	Page     hal.Page
}

func (action *AccountReferralsAction) JSON() {
	action.Do(
		action.loadParams,
		action.checkAllowed,
		action.loadRecords,
		action.loadBalances,
		action.loadResource,
		func() {
			hal.Render(action.W, action.Page)
		},
	)
}

func (action *AccountReferralsAction) loadParams() {
	action.AccountID = action.GetNonEmptyString("id")
}

func (action *AccountReferralsAction) checkAllowed() {
	action.IsAllowed("")
}

func (action *AccountReferralsAction) loadRecords() {
	err := action.CoreQ().Accounts().
		ForReferrer(action.AccountID).
		WithStatistics().
		Select(&action.Records)
	if err != nil {
		action.Log.WithError(err).Error("failed to load accounts")
		action.Err = &problem.ServerError
		return
	}

}

func (action *AccountReferralsAction) loadBalances() {
	action.Balances = map[string][]core.Balance{}
	for _, record := range action.Records {
		balances := []core.Balance{}
		err := action.CoreQ().
			BalancesByAddress(&balances, record.AccountID)
		if err != nil {
			action.Log.WithError(err).Error("Failed to get balances for account")
			action.Err = &problem.ServerError
			return
		}
		action.Balances[record.AccountID] = balances
	}
}

func (action *AccountReferralsAction) loadResource() {
	for _, record := range action.Records {
		var r resource.Account
		err := r.Populate(
			action.Ctx,
			record,
			[]core.Signer{},
			action.Balances[record.AccountID],
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
	action.Page.PopulateLinks()
	action.Page.BaseURL = action.BaseURL()
	action.Page.BasePath = action.Path()
}
