package horizon

import (
	"fmt"

	"github.com/go-errors/errors"
	"gitlab.com/swarmfund/horizon/db2/core"
	"gitlab.com/swarmfund/horizon/render/hal"
	"gitlab.com/swarmfund/horizon/render/problem"
	"gitlab.com/swarmfund/horizon/resource"
	. "gitlab.com/swarmfund/horizon/resource/smartfeetable"
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
	)
}

func (action *AccountFeesAction) loadParams() {
	action.Account = action.GetString("account_id")
	if action.Account == "" {
		action.SetInvalidField("account_id", errors.New(""))
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
	fmt.Println()
	if err != nil {
		action.Err = &problem.ServerError
		action.Log.WithStack(err).Error("Could not get fee from the database")
		return nil
	}
	return result
}

func (action *AccountFeesAction) loadFees() SmartFeeTable {

	result := make(map[string][]core.FeeEntry)

	q := action.CoreQ().FeeEntries()
	result["account"] = action.getFees(q.ForAccount(action.Account))

	q = action.CoreQ().FeeEntries()
	result["account_type"] = action.getFees(q.ForAccountType(action.AccountType))

	//get general fees set for all, not to be confused with fees for General Account Type
	q = action.CoreQ().FeeEntries()
	result["general"] = action.getFees(q.ForAccountType(nil))

	sft := NewSmartFeeTable(result["account"])
	sft.Update(result["account_type"])
	sft.Update(result["general"])
	sft.AddZeroFees(action.Assets)
	return sft
}

func (action *AccountFeesAction) loadResponse() {
	records := action.loadFees()
	if action.Err != nil {
		return
	}

	byAssets := records.GetValuesByAsset()
	response := resource.FeesResponse{}
	response.Fees = map[string][]resource.FeeEntry{}
	var fee resource.FeeEntry
	for _, feesForAsset := range byAssets {
		for _, coreFee := range feesForAsset {
			fee.Populate(coreFee)
			response.Fees[coreFee.Asset] = append(response.Fees[coreFee.Asset], fee)
		}
	}

	hal.Render(action.W, response)
}
