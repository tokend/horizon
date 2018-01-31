package exchange

import (
	"testing"

	"fmt"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/horizon/db2/core"
)

//getAssetsHelper return slice of assets
func getAssetsHelper() []core.Asset {
	assets := make([]core.Asset, 3)
	for i := range assets {
		assets[i] = core.Asset{
			Code:                 fmt.Sprintf("%s%d", "BTC", i),
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

func TestNewConverter(t *testing.T) {
	var (
		q          core.CoreQMock
		assetQMock core.AssetQMock
	)

	defer func() {
		q.AssertExpectations(t)
		assetQMock.AssertExpectations(t)
	}()

	q.On("Assets").Return(&assetQMock).Twice()

	Convey("test NewConverter", t, func() {
		Convey("Failed to load base assets", func() {
			expectedErr := errors.New("Invalid memory address or nil pointer exception")
			assetQMock.On("ForPolicy", uint32(xdr.AssetPolicyBaseAsset)).Return(&assetQMock).Once()

			var emptyAssets []core.Asset
			assetQMock.On("Select").Return(emptyAssets, expectedErr).Once()

			_, err := NewConverter(&q)

			assert.NotNil(t, err)
			assert.Equal(t, expectedErr, errors.Cause(err))
		})
		Convey("SUCCESS", func() {
			assetQMock.On("ForPolicy", uint32(xdr.AssetPolicyBaseAsset)).Return(&assetQMock).Once()
			assetQMock.On("Select").Return(getAssetsHelper(), nil).Once()

			converter, err := NewConverter(&q)
			assert.NoError(t, err)
			assert.NotNil(t, converter)

		})
	})
}
