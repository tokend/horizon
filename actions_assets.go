package horizon

import (
	"gitlab.com/distributed_lab/tokend/horizon/render/hal"
	"gitlab.com/distributed_lab/tokend/horizon/render/problem"
	"gitlab.com/distributed_lab/tokend/horizon/resource"
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
