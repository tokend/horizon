package exchange

import (
	"github.com/stretchr/testify/mock"
	"gitlab.com/swarmfund/horizon/db2/core"
)

type assetProviderQMock struct {
	mock.Mock
}

func (q *assetProviderQMock) GetAssetsForPolicy(policy uint32) ([]core.Asset, error) {
	args := q.Called(policy)
	return args.Get(0).([]core.Asset), args.Error(1)
}

func (q *assetProviderQMock) GetAssetPairsForCodes(baseAssets []string, quoteAssets []string) ([]core.AssetPair, error) {
	args := q.Called(baseAssets, quoteAssets)
	return args.Get(0).([]core.AssetPair), args.Error(1)
}
