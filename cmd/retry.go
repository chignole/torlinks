package cmd

import (
	"chignole/torlinks/internal/symlink"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// retryCmd represents the retry command
var retryCmd = &cobra.Command{
	Use:   "retry",
	Short: "Allows to reprocess failed torrent files.",
	Long: `Allow to reprocess failed torrent files by removing the .delete extension
  of those files in inbox folder.`,
	Run: func(cmd *cobra.Command, args []string) {
		source := viper.GetString("general.torrentsInbox")
		symlink.Retry(source)
	},
}

func init() {
	rootCmd.AddCommand(retryCmd)
}
