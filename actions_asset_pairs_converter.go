package horizon

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/horizon/db2/core"
	"gitlab.com/tokend/horizon/render/hal"
	"gitlab.com/tokend/horizon/render/problem"
)

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
