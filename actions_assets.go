package horizon

import (
	"gitlab.com/swarmfund/horizon/render/hal"
	"gitlab.com/swarmfund/horizon/render/problem"
	"gitlab.com/swarmfund/horizon/resource"
)

type AssetsIndexAction struct {
	Action
	Owner string
	Assets []resource.Asset
}

func (action *AssetsIndexAction) JSON() {
	action.Do(
		action.loadParams,
		action.loadData,
		func() {
			hal.Render(action.W, action.Assets)
		},
	)
}

func (action *AssetsIndexAction) loadParams() {
	action.Owner = action.GetString("owner")
}

func (action *AssetsIndexAction) loadData() {
	assetsQ := action.CoreQ().Assets()
	if action.Owner != "" {
		assetsQ = assetsQ.ForOwner(action.Owner)
	}

	assets, err := assetsQ.Select()
	if err != nil {
		action.Err = &problem.ServerError
		action.Log.WithStack(err).WithError(err).Error("Could not get asset from the database")
		return
	}

	action.Assets = make([]resource.Asset, len(assets))
	for i := range assets {
		action.Assets[i].Populate(&assets[i])
	}
}

type AssetPairsAction struct {
	Action
	Assets []resource.AssetPair
}

func (action *AssetPairsAction) JSON() {
	action.Do(
		action.loadData,
		func() {
			hal.Render(action.W, action.Assets)
		},
	)
}

func (action *AssetPairsAction) loadData() {
	assets, err := action.CoreQ().AssetPairs()
	if err != nil {
		action.Err = &problem.ServerError
		action.Log.WithStack(err).WithError(err).Error("Could not get asset from the database")
		return
	}

	action.Assets = make([]resource.AssetPair, len(assets))
	for i := range assets {
		action.Assets[i].Populate(&assets[i])
	}
}
