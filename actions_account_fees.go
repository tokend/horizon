package horizon

import (
	. "gitlab.com/swarmfund/horizon/fees"
	"gitlab.com/swarmfund/horizon/render/hal"
	"gitlab.com/swarmfund/horizon/render/problem"
	"gitlab.com/swarmfund/horizon/resource"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/regources"
)

// This file contains the actions:
//
// FeesForAccount: renders all the fees for a specific account
//FeesForAccount show all fees for account

type AccountFeesAction struct {
	Action
	AccountType *int32
	AccountID   string

	Records  SmartFeeTable
	Assets   []string
	Response resource.FeesResponse
}

func (action *AccountFeesAction) JSON() {
	action.Do(
		action.loadParams,
		action.loadAssets,
		action.loadFees,
		action.loadResponse,
		func() {
			hal.Render(action.W, action.Response)
		},
	)
}

func (action *AccountFeesAction) loadParams() {
	action.AccountID = action.GetString("account_id")

	acc := action.GetCoreAccount(action.AccountID, action.CoreQ())
	action.AccountType = &acc.AccountType
}

func (action *AccountFeesAction) loadAssets() {
	accountBalances, err := action.CoreQ().Balances().ByAddress(action.AccountID).Select()
	if err != nil {
		action.Log.WithError(err).Error("failed to load balances")
		action.Err = &problem.ServerError
		return
	}

	var assets []string
	for _, balance := range accountBalances {
		assets = append(assets, balance.Asset)
	}
	action.Assets = assets
}

func (action *AccountFeesAction) loadFees() {
	forAccount, err := action.CoreQ().FeeEntries().ForAccount(action.AccountID).Select()
	if err != nil {
		action.Err = &problem.ServerError
		action.Log.WithError(err).Error("can't get account fees from the database")
		return
	}

	forAccountType, err := action.CoreQ().FeeEntries().ForAccountType(action.AccountType).Select()
	if err != nil {
		action.Err = &problem.ServerError
		action.Log.WithError(err).Error("can't get account type fees from the database")
		return

	}

	//get general fees set for all, not to be confused with fees for General Account Type
	generalFees, err := action.CoreQ().FeeEntries().ForAccountType(nil).Select()
	if err != nil {
		action.Err = &problem.ServerError
		action.Log.WithError(err).Error("can't get general fees from the database")
		return
	}

	sft := NewSmartFeeTable(forAccount)
	sft.Update(forAccountType)
	sft.Update(generalFees)
	sft.AddZeroFees(action.Assets)
	action.Records = sft
}

func (action *AccountFeesAction) loadResponse() {
	if action.Err != nil {
		return
	}

	byAssets := action.Records.GetValuesByAsset()
	action.Response.Fees = make(map[xdr.AssetCode][]regources.FeeEntry)
	for _, feesForAsset := range byAssets {
		for _, coreFee := range feesForAsset {
			fee := resource.SmartPopulate(coreFee, *action.AccountType)
			ac := xdr.AssetCode(coreFee.Asset)
			action.Response.Fees[ac] = append(action.Response.Fees[ac], fee)
		}
	}
}
