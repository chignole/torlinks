/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"chignole/torlinks/internal/files"
	"chignole/torlinks/internal/symlink"
	"database/sql"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// linkCmd represents the link command
var linkCmd = &cobra.Command{
	Use:   "link",
	Short: "Search for torrent files and create symlinks to their data.",
	Long:  `Search for torrent files and create symlinks to their data.`,
	Run: func(cmd *cobra.Command, args []string) {
		dbFile := viper.GetString("database.file")
		source := viper.GetString("general.torrentsInbox")

		// Open database
		db, err := sql.Open("sqlite3", dbFile)
		if err != nil {
			log.Printf("[ERROR] Error while opening database file : %v", err)
		}
		defer db.Close()

		// Look for torrent files
		torrents := files.Find(source, ".torrent")

		// Process every torrent files
		for t := range torrents {
			log.Printf("[INFO] Processing : %v\n", torrents[t].File)
			filesToProcess := files.Parse(torrents[t].File)
		processLoop:
			for _, f := range filesToProcess {
				torrentStatus := symlink.Create(f.Path, f.Size, db)
				switch torrentStatus {
				case "linked":
					torrents[t].Status = "linked"
				case "multi":
					torrents[t].Status = "multi"
					break processLoop
				default:
					torrents[t].Status = "unset"
				}
			}
		}

		// Clean inbox after processing
		for _, t := range torrents {
			files.Clean(t)
		}
	},
}

func init() {
	rootCmd.AddCommand(linkCmd)
}
