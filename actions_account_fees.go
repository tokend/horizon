package horizon

import (
	"github.com/go-errors/errors"
	"gitlab.com/swarmfund/horizon/db2/core"
	"gitlab.com/swarmfund/horizon/render/hal"
	"gitlab.com/swarmfund/horizon/render/problem"
	"gitlab.com/swarmfund/horizon/resource"
	. "gitlab.com/swarmfund/horizon/resource/smartfeetable"
	"gitlab.com/tokend/go/xdr"
)

// This file contains the actions:
//
// FeesForAccount: renders all registration requests
//FeesForAccount show all fees

type AccountFeesAction struct {
	Action
	AccountType *int32
	Account     string

	Assets   []core.Asset
	Response resource.FeesResponse
}

func (action *AccountFeesAction) JSON() {
	action.Do(
		action.loadParams,
		action.loadAccountType,
		action.loadAssets,
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
}

func (action *AccountFeesAction) loadAssets() {
	assets, err := action.CoreQ().Assets().Select()
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

func (action *AccountFeesAction) getFees(q core.FeeEntryQI) []core.FeeEntry {
	var result []core.FeeEntry
	err := q.Select(&result)
	if err != nil {
		action.Err = &problem.ServerError
		action.Log.WithError(err).Error("Could not get fee from the database")
		return nil
	}
	return result
}

func (action *AccountFeesAction) loadFees() SmartFeeTable {
	q := action.CoreQ().FeeEntries()
	account := action.getFees(q.ForAccount(action.Account))

	q = action.CoreQ().FeeEntries()
	accountType := action.getFees(q.ForAccountType(action.AccountType))

	//get general fees set for all, not to be confused with fees for General Account Type
	q = action.CoreQ().FeeEntries()
	general := action.getFees(q.ForAccountType(nil))

	sft := NewSmartFeeTable(account)
	sft.Update(accountType)
	sft.Update(general)
	sft.AddZeroFees(action.Assets)
	return sft
}

func (action *AccountFeesAction) loadResponse() {
	records := action.loadFees()
	if action.Err != nil {
		return
	}

	byAssets := records.GetValuesByAsset()
	action.Response.Fees = make(map[xdr.AssetCode][]resource.FeeEntry)
	var fee resource.FeeEntry
	for _, feesForAsset := range byAssets {
		for _, coreFee := range feesForAsset {
			fee.Populate(coreFee)
			ac := xdr.AssetCode(coreFee.Asset)
			action.Response.Fees[ac] = append(action.Response.Fees[ac], fee)
		}
	}

	for _, asset := range action.Assets {
		ac := xdr.AssetCode(asset.Code)
		action.Response.Fees[ac] = action.addDefaultEntriesForAsset(asset, action.Response.Fees[ac])
	}
}

func (action *AccountFeesAction) addDefaultEntriesForAsset(asset core.Asset, entries []resource.FeeEntry) []resource.FeeEntry {
	if entries == nil {
		entries = make([]resource.FeeEntry, 0)
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

func (action *AccountFeesAction) getDefaultFee(asset string, feeType int, subType int64) resource.FeeEntry {
	accountType := int32(-1)

	return resource.FeeEntry{
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
