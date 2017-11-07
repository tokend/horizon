package horizon

import (
	"bullioncoin.githost.io/development/go/xdr"
	"bullioncoin.githost.io/development/horizon/db2/core"
	"bullioncoin.githost.io/development/horizon/render/hal"
	"bullioncoin.githost.io/development/horizon/render/problem"
	"bullioncoin.githost.io/development/horizon/resource"
	"bullioncoin.githost.io/development/horizon/utils"
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

	action.Fees, err = action.CoreQ().FeesByTypeAssetAccount(int(xdr.FeeTypeEmissionFee), asset.Token, int64(xdr.EmissionFeeTypeSecondaryMarket), action.Account)
	if err != nil {
		action.Log.WithError(err).Error("Failed to load fee by asset and type")
		action.Err = &problem.ServerError
		return
	}

	action.Fees = utils.FillFeeGaps(action.Fees, core.FeeEntry{
		FeeType:     int(xdr.FeeTypeEmissionFee),
		Asset:       asset.Token,
		Fixed:       0,
		Percent:     0,
		Subtype:     int64(xdr.EmissionFeeTypeSecondaryMarket),
		AccountType: -1,
	})

}
