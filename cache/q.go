package cache

import (
	"fmt"

	"gitlab.com/swarmfund/horizon/db2/core"
	"gitlab.com/swarmfund/horizon/db2/history"
	"github.com/patrickmn/go-cache"
)

// loads data from db or cache
type QInterface interface {
	// returns error if asset not found
	MustAssetByCode(code string) (core.Asset, error)
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
