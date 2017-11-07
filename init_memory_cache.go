package horizon

import (
	"bullioncoin.githost.io/development/horizon/cache"
)

func initMemoryCache(app *App) {
	app.cacheProvider = cache.NewProvider()
}

func init() {
	appInit.Add("memory_cache", initMemoryCache)
}
