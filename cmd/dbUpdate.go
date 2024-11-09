/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"chignole/torlinks/internal/filescanner"
	"chignole/torlinks/internal/utils"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// updateDbCmd represents the updateDb command
var dbUpdate = &cobra.Command{
	Use:   "dbUpdate",
	Short: "Updates files database.",
	Long:  `Creates an empty database if it doesn't already exists. Scan source directories and add new files to the database.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("[INFO] Database update started")
		dataDirectories := viper.GetStringSlice("general.data")
		database := viper.GetString("database.file")
		dbUpdatePing := viper.GetString("database.ping")
		log.Println(database)
		filescanner.ScanDirectories(dataDirectories[:], database)

		// Sending ping to healthcheck URL
		if dbUpdatePing != "" {
			utils.PingHealthCheck(dbUpdatePing)
		}
	},
}

func init() {
	rootCmd.AddCommand(dbUpdate)
}
