/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"chignole/torlinks/internal/files"
	"chignole/torlinks/internal/utils"
	"database/sql"
	"log"
	"math"
	"os"
	"path"

	"github.com/jedib0t/go-pretty/v6/table"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type torrentStats struct {
	fileName     string
	totalSize    int64
	presentSize  int64
	percentSize  float64
	totalFiles   int64
	presentFiles int64
	biggestFile  string
}

var (
	torrentsStats        []torrentStats
	completeSizeThresold float64
)

// dryRunCmd represents the dryRun command
var dryRunCmd = &cobra.Command{
	Use:   "dryRun",
	Short: "Similar to the Run command, but dry.",
	Run: func(cmd *cobra.Command, args []string) {
		source := viper.GetString("general.torrentsInbox")
		dbFile := viper.GetString("database.file")

		// Open database
		db, err := sql.Open("sqlite3", dbFile)
		if err != nil {
			log.Printf("[ERROR] Error while opening database file : %v", err)
		}
		defer db.Close()

		query := `SELECT COUNT(*) FROM files WHERE size = ?`

		// Search torrent files in source folder
		torrentsList := files.Find(source, ".torrent")

		// Farming stats
		for t := range torrentsList {
			var currentTorrent torrentStats
			filesInTorrent := files.Parse(torrentsList[t].File)

			currentTorrent.fileName = path.Base(torrentsList[t].File)
			currentTorrent.fileName = utils.TruncateString(currentTorrent.fileName, 50)

			for f := range filesInTorrent {
				var exists int
				biggestFileSize := 0

				currentTorrent.totalFiles = currentTorrent.totalFiles + 1
				currentTorrent.totalSize = currentTorrent.totalSize + filesInTorrent[f].Size
				if filesInTorrent[f].Size > int64(biggestFileSize) {
					currentTorrent.biggestFile = filesInTorrent[f].Path
					currentTorrent.biggestFile = utils.TruncateString(currentTorrent.biggestFile, 50)
				}

				err = db.QueryRow(query, filesInTorrent[f].Size).Scan(&exists)
				if exists > 0 {
					currentTorrent.presentFiles = currentTorrent.presentSize + 1
					currentTorrent.presentSize = currentTorrent.presentSize + filesInTorrent[f].Size
				}
			}
			currentTorrent.percentSize = (float64(currentTorrent.presentSize) / float64(currentTorrent.totalSize)) * 100
			currentTorrent.percentSize = math.Round(currentTorrent.percentSize*10) / 10
			if currentTorrent.percentSize >= completeSizeThresold {
				torrentsStats = append(torrentsStats, currentTorrent)
			}
		}

		// Display stats table
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"#", "Filename", "Exp. Size", "Mat. Size", "%", "Exp. Files", "Mat. Files", "Main File"})
		for i, stats := range torrentsStats {
			t.AppendRow([]interface{}{i + 1, stats.fileName, stats.totalSize, stats.presentSize, stats.percentSize, stats.totalFiles, stats.presentFiles, stats.biggestFile})
			t.AppendSeparator()
		}
		t.SetStyle(table.StyleColoredYellowWhiteOnBlack)
		t.Style().Options.SeparateRows = true
		t.Render()
	},
}

func init() {
	rootCmd.AddCommand(dryRunCmd)
	dryRunCmd.Flags().Float64VarP(&completeSizeThresold, "thresold", "t", 0, "Completed filesize thresold to display torrent")
}
