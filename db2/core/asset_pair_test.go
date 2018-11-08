package core

import (
	"testing"

	"gitlab.com/distributed_lab/logan/v3/errors"

	"gitlab.com/tokend/go/amount"

	"github.com/stretchr/testify/assert"
)

func TestAssetPair_ConvertToDestAsset(t *testing.T) {
	var q mockAssetLoader
	defer q.AssertExpectations(t)

	t.Run("Invalid current price", func(t *testing.T) {
		assetPair := AssetPair{}
		res, isConverted, err := assetPair.ConvertToDestAsset("TEST", 123, &q)
		assert.Zero(t, res)
		assert.False(t, isConverted)
		assert.Error(t, err)
	})

	t.Run("Failed to select asset", func(t *testing.T) {
		destCode := "TEST"
		assetPair := AssetPair{CurrentPrice: 2 * amount.One}
		q.On("LoadAsset", destCode).Return(nil, errors.New("not nil"))
		res, isConverted, err := assetPair.ConvertToDestAsset(destCode, 123, &q)
		assert.Zero(t, res)
		assert.False(t, isConverted)
		assert.Error(t, err)
	})

	t.Run("Asset not found", func(t *testing.T) {
		destCode := "TEST"
		assetPair := AssetPair{CurrentPrice: 2 * amount.One}
		q.On("LoadAsset", destCode).Return(nil, nil)
		res, isConverted, err := assetPair.ConvertToDestAsset(destCode, 123, &q)
		assert.Zero(t, res)
		assert.False(t, isConverted)
		assert.Error(t, err)
	})

	t.Run("dest code invalid", func(t *testing.T) {
		assetPair := AssetPair{
			CurrentPrice: 2 * amount.One,
			BaseAsset:    "BASE",
		}
		fakeDestCode := "FAKEBASE"
		q.On("LoadAsset", fakeDestCode).Return(&Asset{TrailingDigits: 5}, nil)
		res, isConverted, err := assetPair.ConvertToDestAsset(fakeDestCode, 123, &q)
		assert.Zero(t, res)
		assert.False(t, isConverted)
		assert.Error(t, err)
	})

	t.Run("Success because quote asset match", func(t *testing.T) {
		assetPair := AssetPair{
			CurrentPrice: 2 * amount.One,
			QuoteAsset:   "QUOTE",
		}
		q.On("LoadAsset", assetPair.QuoteAsset).Return(&Asset{TrailingDigits: 5}, nil)
		res, isConverted, err := assetPair.ConvertToDestAsset(assetPair.QuoteAsset, 123, &q)
		assert.Equal(t, res, int64(250))
		assert.True(t, isConverted)
		assert.NoError(t, err)
	})

	t.Run("Success because base asset match", func(t *testing.T) {
		assetPair := AssetPair{
			CurrentPrice: 2 * amount.One,
			BaseAsset:    "BASE",
		}
		q.On("LoadAsset", assetPair.BaseAsset).Return(&Asset{TrailingDigits: 5}, nil)
		res, isConverted, err := assetPair.ConvertToDestAsset(assetPair.BaseAsset, 123, &q)
		assert.Equal(t, res, int64(70))
		assert.True(t, isConverted)
		assert.NoError(t, err)
	})
}
