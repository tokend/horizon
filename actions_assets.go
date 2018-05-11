package horizon

import (
	"gitlab.com/swarmfund/horizon/render/hal"
	"gitlab.com/swarmfund/horizon/render/problem"
	"gitlab.com/swarmfund/horizon/resource"
	"gitlab.com/swarmfund/horizon/db2/core"
	"gitlab.com/tokend/go/amount"
	"github.com/go-errors/errors"
)

type AssetsIndexAction struct {
	Action
	Owner  string
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
	assets, err := action.CoreQ().AssetPairs().Select()
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

type AssetPairsConverterAction struct {
	Action
	AssetPair   core.AssetPair
	SourceAsset string
	DestAsset   string
	Amount      int64
}

func (action *AssetPairsConverterAction) loadParams() {
	action.SourceAsset = action.GetNonEmptyString("source_asset")
	action.DestAsset = action.GetNonEmptyString("dest_asset")
	action.Amount = action.GetAmount("amount")
	if action.Amount < 0 {
		action.SetInvalidField("amount", errors.New("negative is not allowed"))
		return
	}
}

func (action *AssetPairsConverterAction) JSON() {
	action.Do(
		action.loadParams,
		action.loadData,
	)
}

func (action *AssetPairsConverterAction) loadAssetPair(base, quote string) *core.AssetPair {
	assetPair, err := action.CoreQ().AssetPairs().ByCode(base, quote)
	if err != nil {
		action.Log.WithError(err).Error("Failed to get asset pair by code")
		action.Err = &problem.ServerError
		return nil
	}

	return assetPair
}

func (action *AssetPairsConverterAction) tryLoadMatchingAssetPair() core.AssetPair {
	assetPair := action.loadAssetPair(action.SourceAsset, action.DestAsset)
	if action.Err != nil {
		return core.AssetPair{}
	}

	if assetPair != nil {
		return *assetPair
	}

	assetPair = action.loadAssetPair(action.DestAsset, action.SourceAsset)
	if action.Err != nil {
		return core.AssetPair{}
	}

	if assetPair == nil {
		action.Err = &problem.NotFound
		return core.AssetPair{}
	}

	return *assetPair
}

func (action *AssetPairsConverterAction) loadData() {
	assetPair := action.tryLoadMatchingAssetPair()
	if action.Err != nil {
		return
	}

	result, isConverted, err := assetPair.ConvertToDestAsset(action.DestAsset, action.Amount)
	if err != nil {
		action.Log.WithError(err).Error("Failed to convert amount to dest asset")
		action.Err = &problem.ServerError
		return
	}

	if !isConverted {
		action.SetInvalidField("amount", errors.New("failed to convert due to overflow"))
		return
	}

	hal.Render(action.W, map[string]interface{}{
		"amount": amount.String(result),
	})
}
