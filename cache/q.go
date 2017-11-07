package cache

import (
	"fmt"

	"gitlab.com/distributed_lab/tokend/horizon/db2/core"
	"gitlab.com/distributed_lab/tokend/horizon/db2/history"
	"github.com/patrickmn/go-cache"
)

// loads data from db or cache
type QInterface interface {
	// returns error if asset not found
	MustAssetByCode(code string) (core.Asset, error)
	// tries not load number of coins in circulation, returns error if fails to load
	MustCoinsInCirculationForAsset(masterAccountID, asset string) (core.AssetAmount, error)
}

type Q struct {
	coreQ         core.QInterface
	historyQ      history.QInterface
	cacheProvider *Provider
}

func NewQ(coreQ core.QInterface, historyQ history.QInterface, provider *Provider) *Q {
	return &Q{
		coreQ:         coreQ,
		historyQ:      historyQ,
		cacheProvider: provider,
	}
}

func (q *Q) MustAssetByCode(code string) (core.Asset, error) {
	if asset, isFound := q.cacheProvider.assetCache.Get(code); isFound {
		return asset, nil
	}

	asset, err := q.coreQ.AssetByCode(code)
	if err != nil {
		return core.Asset{}, err
	}

	if asset == nil {
		err = fmt.Errorf("asset %s not found", code)
		return core.Asset{}, err
	}

	q.cacheProvider.assetCache.Set(code, *asset, cache.DefaultExpiration)

	return *asset, nil
}

func (q *Q) MustCoinsInCirculationForAsset(masterAccountID, asset string) (core.AssetAmount, error) {
	// we can ignore master account id here as it never changes
	if asset, isFound := q.cacheProvider.assetAmountCache.Get(asset); isFound {
		return asset, nil
	}

	result, err := q.coreQ.MustCoinsInCirculationForAsset(masterAccountID, asset)
	if err != nil {
		return core.AssetAmount{}, err
	}

	q.cacheProvider.assetAmountCache.Set(asset, result, cache.DefaultExpiration)

	return result, nil
}
