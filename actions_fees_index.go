package horizon

import (
	"database/sql"

	"github.com/go-errors/errors"
	"gitlab.com/swarmfund/horizon/db2/core"
	"gitlab.com/swarmfund/horizon/ledger"
	"gitlab.com/swarmfund/horizon/render/hal"
	"gitlab.com/swarmfund/horizon/render/problem"
	"gitlab.com/swarmfund/horizon/resource"
	. "gitlab.com/swarmfund/horizon/resource/smartfeetable"
	"gitlab.com/tokend/go/xdr"
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
	Records    SmartFeeTable
	Assets     []core.Asset
	Response   resource.FeesResponse
}

func (action *FeesAllAction) JSON() {
	action.Do(
		action.loadParams,
		action.loadData,
		action.loadFees,
		action.loadResponse,
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

	assets, err := action.CoreQ().Assets().Select()
	if err != nil {
		action.Log.WithError(err).Error("Failed to load assets")
		action.Err = &problem.ServerError
		return
	}
	action.Assets = assets
}

func (action *FeesAllAction) getFees(q core.FeeEntryQI) []core.FeeEntry {
	var result []core.FeeEntry
	err := q.Select(&result)
	if err != nil {
		if err != sql.ErrNoRows {
			action.Err = &problem.ServerError
			action.Log.WithStack(err).Error("Could not get fee from the database")
			return nil
		}
	}
	return result
}

func (action *FeesAllAction) loadFees() {
	q := action.CoreQ().FeeEntries()
	result := make(map[string][]core.FeeEntry)

	// for overview we only need general fees
	if action.IsOverview {
		result["overview"] = action.getFees(q.ForAccountType(action.AccountType))
		action.Records = NewSmartFeeTable(result["overview"])
		return
	}

	result["account"] = action.getFees(q.ForAccount(action.Account))
	result["account_type"] = action.getFees(q.ForAccountType(action.AccountType))
	//get general fees set for all, not to be confused with fees for General Account Type
	result["general"] = action.getFees(q.ForAccountType(nil))

	sft := NewSmartFeeTable(result["account"])
	sft.Update(result["account_type"])
	sft.Update(result["general"])
	sft.AddZeroFees(action.Assets)
	action.Records = sft
}

func (action *FeesAllAction) loadResponse() {
	byAssets := action.Records.GetValuesByAsset()
	action.Response.Fees = map[string][]resource.FeeEntry{}
	var fee resource.FeeEntry
	for _, feesForAsset := range byAssets {
		for _, coreFee := range feesForAsset {
			fee.Populate(coreFee)
			action.Response.Fees[coreFee.Asset] = append(action.Response.Fees[coreFee.Asset], fee)
		}
	}
}
