package handlers

import (
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon/db2/core2"
	"gitlab.com/tokend/horizon/exchange"
	"gitlab.com/tokend/horizon/web_v2/ctx"
	"gitlab.com/tokend/horizon/web_v2/resources"
	"gitlab.com/tokend/regources/generated"
	"net/http"
)

// balanceStateConverter - helper struct to convert balance states to different assets
type balanceStateConverter struct {
	converter *exchange.Converter
}

// newBalanceStateConverterForHandler - creates new instance of balanceStateConverter.
// returns nil and renders server error if failed to create
func newBalanceStateConverterForHandler(w http.ResponseWriter, r *http.Request) *balanceStateConverter {
	repo := ctx.CoreRepo(r)
	assetsProvider := struct {
		core2.AssetsQ
		core2.AssetPairsQ
	}{
		AssetsQ:     core2.NewAssetsQ(repo),
		AssetPairsQ: core2.NewAssetPairsQ(repo),
	}

	converter, err := exchange.NewConverter(assetsProvider)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to create new converter")
		ape.Render(w, problems.InternalError())
		return nil
	}

	return newBalanceStateConverted(converter)
}

// newBalanceStateConverted - creates new instance of balanceStateConverter
func newBalanceStateConverted(converter *exchange.Converter) *balanceStateConverter {
	return &balanceStateConverter{
		converter: converter,
	}
}

// Convert - returns converted balance state from existing one
func (c *balanceStateConverter) Convert(balance core2.Balance, toAsset string) (*regources.ConvertedBalanceState, error) {
	convertedAvailable, err := c.converter.TryToConvertWithOneHop(
		int64(balance.Amount),
		balance.AssetCode,
		toAsset,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert available amount")
	}

	convertedLocked, err := c.converter.TryToConvertWithOneHop(
		int64(balance.Locked),
		balance.AssetCode,
		toAsset,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert locked amount")
	}

	if convertedAvailable == nil || convertedLocked == nil {
		convertedState := resources.NewConvertedBalanceState(
			balance,
			regources.Amount(0),
			regources.Amount(0),
			false,
		)

		return &convertedState, nil
	}

	convertedState := resources.NewConvertedBalanceState(
		balance,
		regources.Amount(*convertedAvailable),
		regources.Amount(*convertedLocked),
		true,
	)

	return &convertedState, nil
}
