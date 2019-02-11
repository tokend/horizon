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

func (action *AssetPairsConverterAction) loadData() {
	converter, err := action.CreateConverter()
	if err != nil {
		action.Err = &problem.ServerError
		action.Log.WithError(err).Error("failed to create converter")
		return
	}

	result, err := converter.TryToConvertWithOneHop(action.Amount, action.SourceAsset, action.DestAsset)
	if err != nil {
		action.Log.WithError(err).Error("Failed to convert amount to dest asset")
		action.Err = &problem.ServerError
		return
	}

	if result == nil {
		action.SetInvalidField("amount", errors.New("failed to convert due to overflow or "+
			"there is no asset pair path"))
		return
	}

	hal.Render(action.W, map[string]interface{}{
		"amount": amount.String(*result),
	})
}
