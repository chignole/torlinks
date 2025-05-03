package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var buildDate string

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Displays build version.",
	Long:  `Displays build version.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("[INFO] Build : %s\n", buildDate)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
