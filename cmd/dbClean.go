package cmd

import (
	"chignole/torlinks/internal/filescanner"
	"chignole/torlinks/internal/utils"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// dbCleanCmd represents the dbClean command
var dbCleanCmd = &cobra.Command{
	Use:   "dbClean",
	Short: "Rebuild files database.",
	Long:  `Rebuild files database`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("[INFO] Database cleaning started")
		dataDirectories := viper.GetStringSlice("general.data")
		database := viper.GetString("database.file")
		dbUpdatePing := viper.GetString("database.ping")
		os.Remove(database)
		filescanner.ScanDirectories(dataDirectories[:], database)

		if dbUpdatePing != "" {
			utils.PingHealthCheck(dbUpdatePing)
		}
	},
}

func init() {
	rootCmd.AddCommand(dbCleanCmd)
}
