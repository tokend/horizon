package horizon

import (
	"gitlab.com/swarmfund/horizon/render/hal"
	"gitlab.com/swarmfund/horizon/render/problem"
	"gitlab.com/swarmfund/horizon/resource"
)

type AssetsShowAction struct {
	Action
	Code string
	Asset resource.Asset
}

func (action *AssetsShowAction) JSON() {
	action.Do(
		action.loadParams,
		action.loadData,
		func() {
			hal.Render(action.W, action.Asset)
		},
	)
}

func (action *AssetsShowAction) loadParams() {
	action.Code = action.GetString("code")
}

func (action *AssetsShowAction) loadData() {
	asset, err := action.CoreQ().Assets().ByCode(action.Code)
	if err != nil {
		action.Err = &problem.ServerError
		action.Log.WithStack(err).WithError(err).Error("Failed to load asset by code")
		return
	}

	if asset == nil {
		action.Err = &problem.NotFound
		return
	}

	action.Asset.Populate(asset)
}