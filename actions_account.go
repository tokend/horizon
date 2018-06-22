package horizon

import (
	"time"

	"gitlab.com/swarmfund/horizon/db2/core"
	"gitlab.com/swarmfund/horizon/render/hal"
	"gitlab.com/swarmfund/horizon/render/problem"
	"gitlab.com/swarmfund/horizon/resource"
)

// This file contains the actions:
//
// AccountShowAction: details for single account (including stellar-core state)

// AccountShowAction renders a account summary found by its address.
type AccountShowAction struct {
	Action
	Address  string
	Resource resource.Account
}

// JSON is a method for actions.JSON
func (action *AccountShowAction) JSON() {
	action.Do(
		action.loadParams,
		action.checkAllowed,
		action.loadRecord,
		action.loadLimits,
		action.loadStatistics,
		action.loadBalances,
		action.loadExternalSystemAccountIDs,
		action.loadReferrals,
		func() {
			hal.Render(action.W, action.Resource)
		},
	)
}

func (action *AccountShowAction) loadParams() {
	action.Address = action.GetString("id")
}

func (action *AccountShowAction) checkAllowed() {
	action.IsAllowed(action.Address)
}

func (action *AccountShowAction) loadRecord() {
	coreRecord, err := action.CoreQ().
		Accounts().
		ForAddresses(action.Address).
		WithAccountKYC().
		First()

	if err != nil {
		action.Log.WithError(err).Error("Failed to get account from core DB")
		action.Err = &problem.ServerError
		return
	}

	if coreRecord == nil {
		action.Err = &problem.NotFound
		return
	}

	action.Resource.Populate(action.Ctx, *coreRecord)

	signers, err := action.GetSigners(coreRecord)
	if err != nil {
		action.Log.WithError(err).Error("Failed to get signers")
		action.Err = &problem.ServerError
		return
	}

	action.Resource.Signers.Populate(signers)
}

func (action *AccountShowAction) loadLimits() {
	limits, err := action.CoreQ().LimitsV2().ForAccountByAccountType(action.Address, action.Resource.AccountTypeI)
	if err != nil {
		action.Log.WithError(err).Error("Failed to load limits for account")
		action.Err = &problem.ServerError
		return
	}

	action.populateLimitsV2(limits)
}

func (action *AccountShowAction) loadStatistics() {
	statisticsV2Records, err := action.CoreQ().StatisticsV2().ForAccount(action.Address)
	if err != nil {
		action.Log.WithError(err).Error("Failed to get statistics for account from core DB")
		action.Err = &problem.ServerError
		return
	}

	action.populateStatisticsV2(statisticsV2Records)
}

func (action *AccountShowAction) loadBalances() {
	var balances []core.Balance
	err := action.CoreQ().
		BalancesByAddress(&balances, action.Address)
	if err != nil {
		action.Log.WithError(err).Error("Failed to get balances for account")
		action.Err = &problem.ServerError
		return
	}

	action.Resource.SetBalances(balances)
}

func (action *AccountShowAction) loadExternalSystemAccountIDs() {
	exSysIDs, err := action.CoreQ().ExternalSystemAccountID().ForAccount(action.Address).Select()
	if err != nil {
		action.Log.WithError(err).Error("Failed to load external system account ids")
		action.Err = &problem.ServerError
		return
	}

	action.Resource.ExternalSystemAccounts = make([]resource.ExternalSystemAccountID, 0, len(exSysIDs))
	for i := range exSysIDs {
		if exSysIDs[i].ExpiresAt == nil || *exSysIDs[i].ExpiresAt >= time.Now().Unix() {
			var result resource.ExternalSystemAccountID
			result.Populate(exSysIDs[i])
			action.Resource.ExternalSystemAccounts = append(action.Resource.ExternalSystemAccounts, result)
		}
	}
}

func (action *AccountShowAction) loadReferrals() {
	var coreReferrals []core.Account
	err := action.CoreQ().Accounts().ForReferrer(action.Address).Select(&coreReferrals)
	if err != nil {
		action.Log.WithError(err).Error("Failed to load referrals")
		action.Err = &problem.ServerError
		return
	}

	action.Resource.Referrals = make([]resource.Referral, len(coreReferrals))
	for i := range coreReferrals {
		action.Resource.Referrals[i].Populate(coreReferrals[i])
	}
}

func (action *AccountShowAction) populateLimitsV2(limitsV2Records []core.LimitsV2Entry) {
	for i, limitsV2 := range limitsV2Records {
		action.Resource.LimitsV2 = append(action.Resource.LimitsV2, resource.LimitsV2{})
		action.Resource.LimitsV2[i].Populate(limitsV2)
	}
}

func (action *AccountShowAction) populateStatisticsV2(statisticsV2Records []core.StatisticsV2Entry) {
	for i, statisticsV2 := range statisticsV2Records {
		action.Resource.StatisticsV2 = append(action.Resource.StatisticsV2, resource.StatisticsV2{})
		action.Resource.StatisticsV2[i].Populate(statisticsV2)
	}
}