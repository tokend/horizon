package horizon

import (
	"gitlab.com/tokend/horizon/db2/core"
	"gitlab.com/tokend/horizon/render/hal"
	"gitlab.com/tokend/horizon/render/problem"
	"gitlab.com/tokend/horizon/resource"
)

// This file contains the actions:
//
// EmissionRulesShowAction: renders fees for operationType
type EmissionRulesShowAction struct {
	Action
	AccountID string
	Account   *core.Account
	Asset     string

	Fees []core.FeeEntry
}

// JSON is a method for actions.JSON
func (action *EmissionRulesShowAction) JSON() {
	action.Do(
		action.loadParams,
		action.loadData,
		func() {
			response := resource.FeesResponse{
				Fees: map[string][]resource.FeeEntry{},
			}

			response.Fees[action.Asset] = make([]resource.FeeEntry, len(action.Fees))
			for i := range action.Fees {
				var fee resource.FeeEntry
				fee.Populate(action.Fees[i])
				response.Fees[action.Asset][i] = fee
			}
			hal.Render(action.W, response)
		},
	)
}
func (action *EmissionRulesShowAction) loadParams() {
	action.AccountID = action.GetNonEmptyString("account_id")
	action.Asset = action.GetNonEmptyString("asset")
}

func (action *EmissionRulesShowAction) loadData() {
	var err error
	action.Account, err = action.CoreQ().Accounts().ByAddress(action.AccountID)
	if err != nil {
		action.Log.WithError(err).Error("Failed to load account to get emission rules")
		action.Err = &problem.ServerError
		return
	}

	if action.Account == nil {
		action.Err = &problem.NotFound
		return
	}

	asset, err := action.CoreQ().AssetByCode(action.Asset)
	if err != nil {
		action.Log.WithError(err).Error("Failed to load asset to get emission rules")
		action.Err = &problem.ServerError
		return
	}

	if asset == nil {
		action.Err = &problem.NotFound
		return
	}
}
