package exchange

import (
	"testing"

	"fmt"

	"math"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/horizon/db2/core"
	"gitlab.com/tokend/go/xdr"
)

func getAssetsHelper(assetCode string) []core.Asset {
	assets := make([]core.Asset, 3)
	for i := range assets {
		assets[i] = core.Asset{
			Code:                 fmt.Sprintf("%s%d", assetCode, i),
			Policies:             int32(xdr.AssetPolicyBaseAsset),
			Owner:                "SAV76USXIJOBMEQXPANUOQM6F5LIOTLPDIDVRJBFFE2MDJXG24TAPUU7",
			AvailableForIssuance: 1,
			PreissuedAssetSigner: "GCFXHS4GXL6BVUCXBWXGTITROWLVYXQKQLF4YH5O5JT3YZXCYPAFBJZB",
			MaxIssuanceAmount:    1000,
			Issued:               1,
			LockedIssuance:       0,
			PendingIssuance:      0,
			Details:              nil,
		}
	}
	return assets
}

func getAssetPairsHelper(baseAsset, quoteAsset string, currentPrice int64) []core.AssetPair {
	assetPairs := make([]core.AssetPair, 3)

	for i := range assetPairs {
		assetPairs[i] = core.AssetPair{
			BaseAsset:               fmt.Sprintf("%s%d", baseAsset, i),
			QuoteAsset:              fmt.Sprintf("%s%d", quoteAsset, i),
			CurrentPrice:            currentPrice,
			PhysicalPrice:           1001,
			PhysicalPriceCorrection: 1002,
			MaxPriceStep:            3,
			Policies:                int32(xdr.AssetPolicyBaseAsset),
		}
	}

	return assetPairs
}

func TestNewConverterConverter(t *testing.T) {
	var q assetProviderQMock
	defer q.AssertExpectations(t)

	Convey("Test newConverter()", t, func() {
		Convey("Success to create new converter", func() {
			assets := getAssetsHelper("SUN")
			q.On("GetAssetsForPolicy", uint32(xdr.AssetPolicyBaseAsset)).Return(assets, nil).Once()

			converter, err := newConverter(&q)

			assert.NoError(t, err)
			assert.NotNil(t, converter)

			//check valid
			for i := range assets {
				assert.Equal(t, assets[i].Code, converter.baseAssets[i])
			}
		})
		Convey("Failed to create new converter", func() {
			expectedErr := errors.New("failed to load base assets")
			var emptyAssets []core.Asset
			q.On("GetAssetsForPolicy", uint32(xdr.AssetPolicyBaseAsset)).Return(emptyAssets, expectedErr).Once()

			converter, err := newConverter(&q)

			assert.Nil(t, converter)
			assert.Error(t, err)
			assert.Equal(t, expectedErr, errors.Cause(err))
		})
	})
}

func TestConverter_loadPairsWithBaseAssets(t *testing.T) {
	var q assetProviderQMock
	defer q.AssertExpectations(t)

	assets := getAssetsHelper("SUN")
	q.On("GetAssetsForPolicy", uint32(xdr.AssetPolicyBaseAsset)).Return(assets, nil).Once()

	converter, _ := newConverter(&q)

	Convey("General test for loadPairsWithBaseAssets", t, func() {
		Convey("Success to load with base assets", func() {
			directPairs := getAssetPairsHelper("SUN", "BTC", 1000)
			q.On("GetAssetPairsForCodes", []string{"SUN0"}, []string{"SUN0", "SUN1", "SUN2"}).Return(directPairs, nil).Once()

			reversePairs := getAssetPairsHelper("BTC", "SUN", 1000)
			q.On("GetAssetPairsForCodes", []string{"SUN0", "SUN1", "SUN2"}, []string{"SUN0"}).Return(reversePairs, nil).Once()

			assetPairs, err := converter.loadPairsWithBaseAssets("SUN0")

			assert.NoError(t, err)
			assert.NotNil(t, assetPairs)

		})
		Convey("Failed to load base asset pairs", func() {
			Convey("Failed load direct", func() {
				expectedErr := errors.New("failed to load direct asset pairs")
				var emptyAssetPairs []core.AssetPair
				q.On("GetAssetPairsForCodes", []string{"SUN0"}, []string{"SUN0", "SUN1", "SUN2"}).Return(emptyAssetPairs, expectedErr).Once()

				assetPairs, err := converter.loadPairsWithBaseAssets("SUN0")

				assert.Error(t, err)
				assert.Nil(t, assetPairs)
				assert.Equal(t, expectedErr, errors.Cause(err))
			})
			Convey("Failed load reverse", func() {
				assetPairs := getAssetPairsHelper("SUN", "BTC", 1000)
				q.On("GetAssetPairsForCodes", []string{"SUN0"}, []string{"SUN0", "SUN1", "SUN2"}).Return(assetPairs, nil).Once()

				var emptyAssetPairs []core.AssetPair
				expectedErr := errors.New("failed to load reverse asset pairs")
				q.On("GetAssetPairsForCodes", []string{"SUN0", "SUN1", "SUN2"}, []string{"SUN0"}).Return(emptyAssetPairs, expectedErr).Once()

				assetPairs, err := converter.loadPairsWithBaseAssets("SUN0")

				assert.Error(t, err)
				assert.Nil(t, assetPairs)
				assert.Equal(t, expectedErr, errors.Cause(err))
			})
			Convey("If no rows in qsl table", func() {
				var emptyDirectPairs, emptyReversePairs []core.AssetPair
				q.On("GetAssetPairsForCodes", []string{"SUN0"}, []string{"SUN0", "SUN1", "SUN2"}).Return(emptyDirectPairs, nil).Once()

				q.On("GetAssetPairsForCodes", []string{"SUN0", "SUN1", "SUN2"}, []string{"SUN0"}).Return(emptyReversePairs, nil).Once()

				assetPairs, err := converter.loadPairsWithBaseAssets("SUN0")

				assert.NoError(t, err)
				assert.Nil(t, assetPairs)
			})
		})
	})
}

func TestTryLoadDirect(t *testing.T) {

	var q assetProviderQMock
	defer q.AssertExpectations(t)

	assets := getAssetsHelper("SUN")
	q.On("GetAssetsForPolicy", uint32(xdr.AssetPolicyBaseAsset)).Return(assets, nil).Once()

	converter, _ := newConverter(&q)

	Convey("General test for tryLoadDirect", t, func() {
		Convey("Success to load direct", func() {
			assetPairs := getAssetPairsHelper("SUN", "BTC", 1000)
			q.On("GetAssetPairsForCodes", []string{"SUN0", "BTC0"}, []string{"SUN0", "BTC0"}).Return(assetPairs, nil).Once()

			res, err := converter.tryLoadDirect("SUN0", "BTC0")

			assert.NoError(t, err)
			assert.NotNil(t, res)
		})
		Convey("When no errors and pair == nil", func() {
			assetPairs := getAssetPairsHelper("SUN", "BTC", 1000)
			q.On("GetAssetPairsForCodes", []string{"SUN0", "SUN0"}, []string{"SUN0", "SUN0"}).Return(assetPairs, nil).Once()

			res, err := converter.tryLoadDirect("SUN0", "SUN0")

			assert.NoError(t, err)
			assert.Nil(t, res)

		})
		Convey("Failed to load direct", func() {
			expectedErr := errors.New("failed to load direct asset pairs")
			var emptyAssetPairs []core.AssetPair
			q.On("GetAssetPairsForCodes", []string{"SUN0", "BTC0"}, []string{"SUN0", "BTC0"}).Return(emptyAssetPairs, expectedErr).Once()

			assetPairs, err := converter.tryLoadDirect("SUN0", "BTC0")

			assert.Error(t, err)
			assert.Nil(t, assetPairs)
			assert.Equal(t, expectedErr, errors.Cause(err))
		})
	})
}

func TestConverter_convertWithMaxPath(t *testing.T) {
	var q assetProviderQMock
	defer q.AssertExpectations(t)

	assets := getAssetsHelper("SUN")
	q.On("GetAssetsForPolicy", uint32(xdr.AssetPolicyBaseAsset)).Return(assets, nil).Once()

	converter, _ := newConverter(&q)

	amount := int64(1000)
	fromAsset, toAsset := "SUN0", "BTC0"

	fromPairs, toPairs := getAssetPairsHelper("SUN", "BTC", 1000), getAssetPairsHelper("BTC", "SUN", 1000)

	Convey("Test convertWithMaxPath", t, func() {
		Convey("Success to convertWithMaxPath", func() {
			path, err := converter.convertWithMaxPath(amount, fromAsset, toAsset, fromPairs, toPairs)

			assert.NoError(t, err)
			assert.NotNil(t, path)

		})

		Convey("Failed to convertWithMaxPath", func() {
			Convey("Failed to convert to asset to hop asset", func() {
				fromPairs := getAssetPairsHelper("SUN", "BTC", 0)

				path, err := converter.convertWithMaxPath(amount, fromAsset, toAsset, fromPairs, toPairs)

				assert.Error(t, err)
				assert.Nil(t, path)

			})
			Convey("Failed to convert from asset to hop asset", func() {
				//set invalid price
				fromPairs := getAssetPairsHelper("SUN", "BTC", 0)

				path, err := converter.convertWithMaxPath(amount, fromAsset, toAsset+"hgfds", fromPairs, toPairs)

				assert.Error(t, err)
				assert.Nil(t, path)
			})
		})
	})
}

func TestConverter_TryToConvertWithOneHop(t *testing.T) {
	var q assetProviderQMock
	defer q.AssertExpectations(t)

	assets := getAssetsHelper("SUN")
	q.On("GetAssetsForPolicy", uint32(xdr.AssetPolicyBaseAsset)).Return(assets, nil).Once()

	converter, _ := newConverter(&q)

	amount := int64(1000)
	fromAsset, toAsset := "SUN0", "BTC0"

	Convey("General test for TryToConvertWithOneHop", t, func() {
		Convey("Success to convert in one hop ", func() {
			Convey("Convert to dest asset success", func() {
				assetPairs := getAssetPairsHelper("SUN", "BTC", 1000)
				q.On("GetAssetPairsForCodes", []string{"SUN0", "BTC0"}, []string{"SUN0", "BTC0"}).Return(assetPairs, nil).Once()

				res, err := converter.TryToConvertWithOneHop(amount, fromAsset, toAsset)

				assert.NoError(t, err)
				assert.NotNil(t, res)
			})

			Convey("Success to load pairs with base assets", func() {
				assetPairs := getAssetPairsHelper(fromAsset, toAsset, 1000)
				q.On("GetAssetPairsForCodes", []string{"SUN0", "BTC0"}, []string{"SUN0", "BTC0"}).Return(assetPairs, nil).Once()

				//success to load fromAsset
				q.On("GetAssetPairsForCodes", []string{"SUN0"}, []string{"SUN0", "SUN1", "SUN2"}).Return(assetPairs, nil).Once()
				q.On("GetAssetPairsForCodes", []string{"SUN0", "SUN1", "SUN2"}, []string{"SUN0"}).Return(assetPairs, nil).Once()

				//success to load toAsset
				q.On("GetAssetPairsForCodes", []string{"BTC0"}, []string{"SUN0", "SUN1", "SUN2"}).Return(assetPairs, nil).Once()
				q.On("GetAssetPairsForCodes", []string{"SUN0", "SUN1", "SUN2"}, []string{"BTC0"}).Return(assetPairs, nil).Once()

				res, err := converter.TryToConvertWithOneHop(amount, fromAsset, toAsset)

				assert.NoError(t, err)
				assert.NotNil(t, res)
			})
		})
		Convey("Failed to convert with one hop", func() {
			Convey("failed to convert to dest asset", func() {
				Convey("failed to load pairs for toAsset", func() {
					assetPairs := getAssetPairsHelper(fromAsset, toAsset, 1000)
					q.On("GetAssetPairsForCodes", []string{"SUN0", "BTC0"}, []string{"SUN0", "BTC0"}).Return(assetPairs, nil).Once()

					//success to load fromAsset
					q.On("GetAssetPairsForCodes", []string{"SUN0"}, []string{"SUN0", "SUN1", "SUN2"}).Return(assetPairs, nil).Once()
					q.On("GetAssetPairsForCodes", []string{"SUN0", "SUN1", "SUN2"}, []string{"SUN0"}).Return(assetPairs, nil).Once()

					var emptyAssetPairs []core.AssetPair
					expectedErr := errors.New("failed to load pairs with base asset for to asset")
					q.On("GetAssetPairsForCodes", []string{"BTC0"}, []string{"SUN0", "SUN1", "SUN2"}).Return(emptyAssetPairs, expectedErr).Once()

					res, err := converter.TryToConvertWithOneHop(amount, fromAsset, toAsset)

					assert.Nil(t, res)
					assert.Equal(t, expectedErr, errors.Cause(err))
				})
				Convey("failed to load pairs for fromAsset", func() {
					assetPairs := getAssetPairsHelper(fromAsset, toAsset, 1000)
					q.On("GetAssetPairsForCodes", []string{"SUN0", "BTC0"}, []string{"SUN0", "BTC0"}).Return(assetPairs, nil).Once()

					var emptyAssetPairs []core.AssetPair

					expectedErr := errors.New("failed to load pairs with base asset for from asset")
					q.On("GetAssetPairsForCodes", []string{"SUN0"}, []string{"SUN0", "SUN1", "SUN2"}).Return(emptyAssetPairs, expectedErr).Once()

					res, err := converter.TryToConvertWithOneHop(amount, fromAsset, toAsset)

					assert.Nil(t, res)
					assert.Error(t, err)
					assert.Equal(t, expectedErr, errors.Cause(err))

				})
				Convey("When !isConverted", func() {
					assetPairs := getAssetPairsHelper("SUN", "BTC", math.MaxInt64)
					q.On("GetAssetPairsForCodes", []string{"SUN0", "BTC0"}, []string{"SUN0", "BTC0"}).Return(assetPairs, nil).Once()

					amount = math.MaxInt64

					res, err := converter.TryToConvertWithOneHop(amount, fromAsset, toAsset)

					assert.Nil(t, res)
					assert.NoError(t, err)
				})
				Convey("failed to convert because of invalid price", func() {
					assetPairs := getAssetPairsHelper("SUN", "BTC", 0)
					q.On("GetAssetPairsForCodes", []string{"SUN0", "BTC0"}, []string{"SUN0", "BTC0"}).Return(assetPairs, nil).Once()

					res, err := converter.TryToConvertWithOneHop(amount, fromAsset, toAsset)

					assert.Nil(t, res)
					assert.Error(t, err)
				})

			})
			Convey("Failed to load direct", func() {
				expectedErr := errors.New("failed to load direct asset pairs")
				var emptyAssetPairs []core.AssetPair
				q.On("GetAssetPairsForCodes", []string{"SUN0", "BTC0"}, []string{"SUN0", "BTC0"}).Return(emptyAssetPairs, expectedErr).Once()

				res, err := converter.TryToConvertWithOneHop(amount, fromAsset, toAsset)

				assert.Nil(t, res)
				assert.Error(t, err)
				assert.Equal(t, expectedErr, errors.Cause(err))
			})
			Convey("When from asset == to asset", func() {
				fromAsset, toAsset = "SUN0", "SUN0"

				res, err := converter.TryToConvertWithOneHop(amount, fromAsset, toAsset)
				assert.Equal(t, *res, amount)
				assert.NoError(t, err)
			})
		})
	})
}
