package horizon

import (
	"database/sql"

	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/horizon/db2/core"
	"gitlab.com/tokend/horizon/ledger"
	"gitlab.com/tokend/horizon/render/hal"
	"gitlab.com/tokend/horizon/render/problem"
	"gitlab.com/tokend/horizon/resource"
	"github.com/go-errors/errors"
)

// This file contains the actions:
//
// FeesAllAction: renders all registration requests

//FeesAllAction show all fees

type FeesAllAction struct {
	Action
	AccountType *int32
	Account     string

	IsOverview bool

	Response resource.FeesResponse
}

func (action *FeesAllAction) JSON() {
	action.Do(
		action.loadParams,
		action.loadData,
		func() {
			hal.Render(action.W, action.Response)
		},
	)
}

func (action *FeesAllAction) loadParams() {
	action.Account = action.GetString("account_id")
	action.AccountType = action.getAccountType("account_type")
	if action.Account != "" && action.AccountType != nil {
		action.SetInvalidField("account_type", errors.New("It's not allowed to set both filters"))
	}
}

func (action *FeesAllAction) getAccountType(name string) *int32 {
	rawAccountType := action.GetInt32(name)
	if action.Err != nil {
		return nil
	}

	if rawAccountType == 0 {
		return nil
	}

	for _, accountType := range xdr.AccountTypeAll {
		if int32(accountType) == rawAccountType {
			return &rawAccountType
		}
	}

	action.SetInvalidField(name, errors.New("Invalid"))
	return nil
}

func (action *FeesAllAction) loadData() {
	var ledgerHeader core.LedgerHeader
	err := action.CoreQ().LedgerHeaderBySequence(&ledgerHeader, ledger.CurrentState().CoreLatest)
	if err != nil {
		action.Log.WithError(err).Error("Failed to get latest ledger")
		action.Err = &problem.ServerError
		return
	}

	q := action.CoreQ().FeeEntries()
	// for the overview we need to return all the fee rules we have, so we just ignore filters
	if !action.IsOverview {
		// pass all the filters. Q will resolve them correctly
		q = q.ForAccount(action.Account).ForAccountType(action.AccountType)
	}

	actualFees := []core.FeeEntry{}
	err = q.Select(&actualFees)
	if err != nil {
		if err != sql.ErrNoRows {
			action.Err = &problem.ServerError
			action.Log.WithStack(err).WithError(err).Error("Could not get fee from the database")
			return
		}

		err = nil
	}

	// convert to map of resources
	action.Response.Fees = map[string][]resource.FeeEntry{}
	var fee resource.FeeEntry
	for _, coreFee := range actualFees {
		fee.Populate(coreFee)
		action.Response.Fees[coreFee.Asset] = append(action.Response.Fees[coreFee.Asset], fee)
	}

	// for overview we do not need to populate default fees
	if action.IsOverview {
		return
	}

	assets, err := action.CoreQ().Assets()
	if err != nil {
		action.Log.WithError(err).Error("Failed to load assets")
		action.Err = &problem.ServerError
		return
	}

	for _, asset := range assets {
		action.Response.Fees[asset.Code] = action.addDefaultEntriesForAsset(asset, action.Response.Fees[asset.Code])
	}
}

func feesContainsType(feeType int, entries []resource.FeeEntry) bool {
	for _, entry := range entries {
		if entry.FeeType == feeType {
			return true
		}
	}

	return false
}

func (action *FeesAllAction) addDefaultEntriesForAsset(asset core.Asset, entries []resource.FeeEntry) []resource.FeeEntry {
	for _, feeType := range xdr.FeeTypeAll {

		isAmountRangeSupported := feeType != xdr.FeeTypeReferralFee
		if !isAmountRangeSupported && feesContainsType(int(feeType), entries) {
			continue
		}

		switch feeType {
		case xdr.FeeTypeEmissionFee:
			for _, subType := range xdr.EmissionFeeTypeAll {
				entries = append(entries, action.getDefaultFee(asset.Code, int(feeType), int64(subType)))
			}
		default:
			entries = append(entries, action.getDefaultFee(asset.Code, int(feeType), int64(0)))
		}
	}

	return entries
}

func (action *FeesAllAction) getDefaultFee(asset string, feeType int, subType int64) resource.FeeEntry {
	accountType := int32(-1)
	if action.AccountType != nil {
		accountType = *action.AccountType
	}

	return resource.FeeEntry{
		Asset:       asset,
		FeeType:     feeType,
		Subtype:     subType,
		Percent:     "0",
		Fixed:       "0",
		LowerBound:  "0",
		UpperBound:  "0",
		AccountType: accountType,
		AccountID:   action.Account,
	}
}
