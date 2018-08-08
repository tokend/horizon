package horizon

import (
	"github.com/go-errors/errors"
	"gitlab.com/swarmfund/horizon/db2/core"
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
	Account     string

	Records  SmartFeeTable
	Assets   []core.Asset
	Response resource.FeesResponse
}

func (action *AccountFeesAction) JSON() {
	action.Do(
		action.loadParams,
		action.loadAccountType,
		action.loadAssets,
		action.loadFees,
		action.loadResponse,
		func() {
			hal.Render(action.W, action.Response)
		},
	)
}

func (action *AccountFeesAction) loadParams() {
	action.Account = action.GetString("account_id")
	if action.Account == "" {
		action.SetInvalidField("account_id", errors.New("cannot be blank"))
	}

	acc, err := action.CoreQ().Accounts().ByAddress(action.Account)
	if acc == nil {
		action.Log.Error("account does not exist")
		action.Err = &problem.BadRequest
	}
	if err != nil {
		action.Log.WithError(err)
		action.Err = &problem.ServerError
	}
}

func (action *AccountFeesAction) loadAssets() {
	assets, err := action.CoreQ().Assets().ForOwner(action.Account).Select()
	if err != nil {
		action.Log.WithError(err).Error("Failed to load assets")
		action.Err = &problem.ServerError
		return
	}
	action.Assets = assets
}
func (action *AccountFeesAction) loadAccountType() {
	account, err := action.CoreQ().Accounts().ByAddress(action.Account)
	if err != nil {
		action.Log.WithError(err).Error("Failed to get account info")
		action.Err = &problem.ServerError
		return
	}
	action.AccountType = &account.AccountType
}

func (action *AccountFeesAction) loadFees() {
	account, err := action.CoreQ().FeeEntries().ForAccount(action.Account).Select()
	if err != nil {
		action.Err = &problem.ServerError
		action.Log.WithError(err).Error("Could not get account fees from the database")
	}

	accountType, err := action.CoreQ().FeeEntries().ForAccountType(action.AccountType).Select()
	if err != nil {
		action.Err = &problem.ServerError
		action.Log.WithError(err).Error("Could not get account type fees from the database")

	}

	//get general fees set for all, not to be confused with fees for General Account Type
	general, err := action.CoreQ().FeeEntries().ForAccountType(nil).Select()
	if err != nil {
		action.Err = &problem.ServerError
		action.Log.WithError(err).Error("Could not get general fees from the database")
	}

	sft := NewSmartFeeTable(account)
	sft.Update(accountType)
	sft.Update(general)
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

	for _, asset := range action.Assets {
		ac := xdr.AssetCode(asset.Code)
		action.Response.Fees[ac] = action.addDefaultEntriesForAsset(asset, action.Response.Fees[ac])
	}
}

func (action *AccountFeesAction) addDefaultEntriesForAsset(asset core.Asset, entries []regources.FeeEntry) []regources.FeeEntry {
	if entries == nil {
		entries = make([]regources.FeeEntry, 0)
	}
	for _, feeType := range xdr.FeeTypeAll {
		subtypes := []int64{0}
		if feeType == xdr.FeeTypePaymentFee {
			subtypes = []int64{int64(xdr.PaymentFeeTypeIncoming), int64(xdr.PaymentFeeTypeOutgoing)}
		}

		for _, subtype := range subtypes {
			entries = append(entries, action.getDefaultFee(asset.Code, int(feeType), subtype))
		}
	}

	return entries
}

func (action *AccountFeesAction) getDefaultFee(asset string, feeType int, subType int64) regources.FeeEntry {
	accountType := int32(-1)
	if action.AccountType != nil {
		accountType = *action.AccountType
	}
	return regources.FeeEntry{
		Asset:       asset,
		FeeType:     feeType,
		Subtype:     subType,
		Percent:     "0",
		Fixed:       "0",
		LowerBound:  "0",
		UpperBound:  "0",
		AccountType: accountType,
		FeeAsset:    asset,
	}
}
