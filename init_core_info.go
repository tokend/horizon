package horizon

func initStellarCoreInfo(app *App) {
	app.UpdateStellarCoreInfo()
}

func init() {
	appInit.Add("stellarCoreInfo", initStellarCoreInfo, "core_connector", "app-context", "log")
}
