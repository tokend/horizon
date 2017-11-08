package horizon

import (
	"gitlab.com/tokend/horizon/cache"
)

func initMemoryCache(app *App) {
	app.cacheProvider = cache.NewProvider()
}

func init() {
	appInit.Add("memory_cache", initMemoryCache)
}
