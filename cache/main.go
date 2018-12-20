package cache

import "time"

type Provider struct {
	assetCache *CoreAssetCache
}

func NewProvider() *Provider {
	return &Provider{
		assetCache: newcoreAssetCache(time.Duration(1)*time.Hour, time.Duration(1)*time.Hour),
	}
}
