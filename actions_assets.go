package horizon

import (
	"bullioncoin.githost.io/development/horizon/render/hal"
	"bullioncoin.githost.io/development/horizon/render/problem"
	"bullioncoin.githost.io/development/horizon/resource"
)

type AssetsAllAction struct {
	Action
	Assets []resource.Asset
}

func (action *AssetsAllAction) JSON() {
	action.Do(
		action.loadData,
		func() {
			hal.Render(action.W, action.Assets)
		},
	)
}

func (action *AssetsAllAction) loadData() {
	assets, err := action.CoreQ().Assets()
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
