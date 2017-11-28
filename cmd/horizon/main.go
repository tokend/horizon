package main

import (
	"log"
	"runtime"

	"gitlab.com/swarmfund/horizon"
	"gitlab.com/swarmfund/horizon/config"
	"github.com/spf13/cobra"
)

var app *horizon.App
var conf config.Config
var version string

var rootCmd *cobra.Command

func main() {
	if version != "" {
		horizon.SetVersion(version)
	}
	runtime.GOMAXPROCS(runtime.NumCPU())
	rootCmd.Execute()
}

func init() {
	rootCmd = &cobra.Command{
		Use:   "horizon",
		Short: "client-facing api server for the stellar network",
		Long:  "client-facing api server for the stellar network",
		Run: func(cmd *cobra.Command, args []string) {
			initApp(cmd, args)
			app.Serve()
		},
	}

	conf.DefineConfigStructure(rootCmd)

	rootCmd.AddCommand(dbCmd)
}

func initApp(cmd *cobra.Command, args []string) {
	err := conf.Init()
	if err != nil {
		log.Println("Failed to init config")
		log.Fatal(err.Error())
	}
	app, err = horizon.NewApp(conf)

	if err != nil {
		log.Fatal(err.Error())
	}
}
