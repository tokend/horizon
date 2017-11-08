package cache

import "time"

type Provider struct {
	assetCache       *CoreAssetCache
	assetAmountCache *CoreAssetAmountCache
}

func NewProvider() *Provider {
	return &Provider{
		assetCache:       newcoreAssetCache(time.Duration(1)*time.Hour, time.Duration(1)*time.Hour),
		assetAmountCache: newcoreAssetAmountCache(time.Duration(5)*time.Minute, time.Duration(1)*time.Minute),
	}
}
