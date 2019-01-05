package horizon

func initCoreInfo(app *App) {
	app.UpdateCoreInfo()
}

func init() {
	appInit.Add("core-info", initCoreInfo, "core_connector", "app-context", "log")
}
