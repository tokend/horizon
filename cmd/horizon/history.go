package main

import (
	"log"

	"github.com/spf13/cobra"
	"gitlab.com/tokend/horizon/db2/history/schema"
)

var historyCmd = &cobra.Command{
	Use:   "history [command]",
	Short: "commands to manage horizon's history",
}

func init() {
	historyCmd.AddCommand(historyDBMigrateCmd)
}

var historyDBMigrateCmd = &cobra.Command{
	Use:   "migrate [up|down|redo] [COUNT]",
	Short: "migrate schema",
	Long:  "performs a schema migration command",
	Run: func(cmd *cobra.Command, args []string) {
		err := conf.Init()
		if err != nil {
			log.Fatal(err)
		}
		migrateDB(cmd, args, conf.DatabaseURL, schema.Migrate)
	},
}
