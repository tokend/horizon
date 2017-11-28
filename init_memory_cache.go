package horizon

import (
	"gitlab.com/swarmfund/horizon/cache"
)

func initMemoryCache(app *App) {
	app.cacheProvider = cache.NewProvider()
}

func init() {
	appInit.Add("memory_cache", initMemoryCache)
}
