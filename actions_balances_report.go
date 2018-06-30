package horizon

import (
	"gitlab.com/swarmfund/horizon/exchange"
	"gitlab.com/swarmfund/horizon/render/hal"
)

type Records struct {
	TotalZeroBalanceAccountsCount     int `json:"total_zero_balance_accounts_count"`
	TotalPositiveBalanceAccountsCount int `json:"total_positive_balance_accounts_count"`
	TotalAboveThresholdAccountsCount  int `json:"total_above_threshold_accounts_count"`
	TotalBelowThresholdAccountsCount  int `json:"total_below_threshold_accounts_count"`
}

type BalancesReportAction struct {
	Action
	threshold int64
	assetType string
	Records   Records `json:"records"`
}

func (action *BalancesReportAction) JSON() {
	action.Do(
		action.loadParams,
		action.loadRecords,
		func() {
			hal.Render(action.W, &action.Records)
		},
	)
}

func (action *BalancesReportAction) loadParams() {
	action.threshold = action.GetInt64("threshold")
	action.assetType = action.GetString("asset_type")
}

func (action *BalancesReportAction) loadRecords() {
	// find positive balances
	balances, err := action.CoreQ().Balances().NonZero().Select()
	if err != nil {
		action.Log.WithError(err).Warn("failed to filter non-zero balances")
	}

	USDAmountFromAccount := make(map[string]int64)

	converter, err := exchange.NewConverter(action.CoreQ())
	if err != nil {
		action.Log.Warn("failed to create converter")
	}

	// group balances by account and assign converted amounts
	for _, balance := range balances {
		amount, err := converter.TryToConvertWithOneHop(balance.Amount+balance.Locked, balance.Asset, action.assetType)
		if err != nil {
			action.Log.Warn("failed to convert to asset type, skipping")
		}
		if amount == nil {
			action.Log.WithField("base", balance.Asset).WithField("quote", action.assetType).Warn("the asset is unconvertable, skipping")
		} else {
			USDAmountFromAccount[balance.AccountID] += *amount
		}
	}

	// filter and count by threshold
	for account := range USDAmountFromAccount {
		if USDAmountFromAccount[account] < action.threshold {
			action.Records.TotalBelowThresholdAccountsCount += 1
		} else {
			action.Records.TotalAboveThresholdAccountsCount += 1
		}
	}

	// find zero balances
	zeroBalances, err := action.CoreQ().Balances().Zero().Select()
	if err != nil {
		action.Log.WithError(err).Warn("failed to filter zero balances")
	}

	isZeroAmountAccount := make(map[string]bool)

	// group zero balances by account
	for _, balance := range zeroBalances {
		isZeroAmountAccount[balance.AccountID] = true
	}

	action.Records.TotalZeroBalanceAccountsCount = len(isZeroAmountAccount)
	action.Records.TotalPositiveBalanceAccountsCount = len(USDAmountFromAccount)
}
