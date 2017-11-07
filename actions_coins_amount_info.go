package horizon

import (
	"bullioncoin.githost.io/development/horizon/render/hal"
	"bullioncoin.githost.io/development/horizon/render/problem"
	"bullioncoin.githost.io/development/horizon/resource"
)

type CoinsAmountInfoAction struct {
	Action
	Resource resource.CoinsAmountInfo
}

func (action *CoinsAmountInfoAction) JSON() {
	action.Do(
		action.ValidateBodyType,
		action.checkAllowed,
		action.performRequest,
		func() {
			hal.Render(action.W, action.Resource)
		},
	)
}

func (action *CoinsAmountInfoAction) checkAllowed() {
	action.IsAllowed("")
}

func (action *CoinsAmountInfoAction) performRequest() {
	availableEmissions, err := action.App.obtainAvailableEmissions()
	if err != nil {
		action.Log.WithError(err).Error("Failed to load available emissions")
		action.Err = &problem.ServerError
		return
	}

	coinsInCirculation, err := action.obtainCoinsInCirculation()
	if err != nil {
		action.Log.WithError(err).Error("Failed to load amount of the coins in circulation")
		action.Err = &problem.ServerError
		return
	}

	assetStats, err := action.CoreQ().AssetStats(action.App.CoreInfo.MasterAccountID)
	if err != nil {
		action.Log.WithError(err).Error("failed to get asset stats")
		action.Err = &problem.ServerError
		return
	}

	action.Resource = resource.NewCoinsAmountInfo(availableEmissions, coinsInCirculation, assetStats)
}
