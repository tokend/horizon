package horizon

import (
	"time"

	"bullioncoin.githost.io/development/go/amount"
	"bullioncoin.githost.io/development/go/xdr"
	"gitlab.com/distributed_lab/tokend/horizon/db2/core"
	"gitlab.com/distributed_lab/tokend/horizon/render/hal"
	"gitlab.com/distributed_lab/tokend/horizon/render/problem"
	"gitlab.com/distributed_lab/tokend/horizon/resource"
)

// This file contains the actions:
//
// AccountShowAction: details for single account (including stellar-core state)

// AccountShowAction renders a account summary found by its address.
type AccountShowAction struct {
	Action
	Address     string
	CoreRecord  *core.Account
	CoreSigners []core.Signer
	CoreLimits  core.Limits
	Balances    []core.Balance
	Policies    []core.ExchangePolicies
	Resource    resource.Account
}

// JSON is a method for actions.JSON
func (action *AccountShowAction) JSON() {
	action.Do(
		action.loadParams,
		action.checkAllowed,
		action.loadRecord,
		action.loadBalances,
		action.loadResource,
		func() {
			action.Resource.IncentivePerCoinExpiresAt = time.Date(2017, 10, 1, 0, 0, 0, 0, time.UTC).Unix()
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

func getTokenBalance(balances []core.Balance, token string) int64 {
	for i := range balances {
		if balances[i].Asset == token {
			return int64(balances[i].Amount)
		}
	}

	return 0
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

	for i := range balances {
		assetCode := balances[i].Asset
		asset, err := action.CachedQ().MustAssetByCode(assetCode)
		if err != nil {
			action.Log.WithError(err).Error("Failed to load asset")
			action.Err = &problem.ServerError
			return
		}

		if !action.IsAdmin && !asset.IsVisibleForUser(action.CoreRecord) {
			continue
		}

		if asset.Token == "" {
			action.Balances = append(action.Balances, balances[i])
			continue
		}

		coinsInCirculation, err := action.CachedQ().MustCoinsInCirculationForAsset(action.App.CoreInfo.MasterAccountID, asset.Token)
		if err != nil {
			action.Err = &problem.ServerError
			action.Log.WithError(err).Error("Failed to load coins in circulation")
			return
		}

		holdingTokens := getTokenBalance(balances, asset.Token)
		if coinsInCirculation.Amount == 0 || holdingTokens == 0 {
			action.Balances = append(action.Balances, balances[i])
			continue
		}

		incentive, isOverflow := amount.BigDivide(holdingTokens, 100*amount.One, coinsInCirculation.Amount, amount.ROUND_DOWN)
		if isOverflow {
			incentive = 0
		}

		balances[i].IncentivePerCoin = incentive
		action.Balances = append(action.Balances, balances[i])
	}
}

func (action *AccountShowAction) loadRecord() {
	var err error
	action.CoreRecord, err = action.CoreQ().
		Accounts().
		ForAddresses(action.Address).
		WithStatistics().
		First()

	if err != nil {
		action.Log.WithError(err).Error("Failed to get account from core DB")
		action.Err = &problem.ServerError
		return
	}

	if action.CoreRecord == nil {
		action.Err = &problem.NotFound
		return
	}

	action.CoreRecord.Statistics.ClearObsolete(time.Now().UTC())

	action.CoreSigners, err = action.GetSigners(action.CoreRecord)
	if err != nil {
		action.Log.WithError(err).Error("Failed to get signers for account")
		action.Err = &problem.ServerError
		return
	}

	action.CoreLimits, err = action.CoreQ().LimitsForAccount(action.CoreRecord.AccountID, action.CoreRecord.AccountType)

	if xdr.AccountType(action.CoreRecord.AccountType) == xdr.AccountTypeExchange {
		err = action.CoreQ().
			PoliciesByExchangeID(&action.Policies, action.Address)
		if err != nil {
			action.Log.WithError(err).Error("Failed to get policies for exchange")
			action.Err = &problem.ServerError
			return
		}
	}
}

func (action *AccountShowAction) loadResource() {
	err := action.Resource.Populate(
		action.Ctx,
		*action.CoreRecord,
		action.CoreSigners,
		action.Balances,
		&action.CoreLimits,
		action.Policies,
		action.App.CoreInfo.DemurragePeriod,
	)
	if err != nil {
		action.Log.WithError(err).Error("Failed to populate account response")
		action.Err = &problem.ServerError
		return
	}
}
