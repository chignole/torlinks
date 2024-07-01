/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"chignole/torlinks/internal/files"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// inboxCmd represents the inbox command
var inboxCmd = &cobra.Command{
	Use:   "inbox",
	Short: "Provides some useful stats about your inbox folder.",
	Long:  `Provides some useful stats about your inbox folder.`,
	Run: func(cmd *cobra.Command, args []string) {
		source := viper.GetString("general.source")

		totalSize := getTotalSize(source)
		log.Printf("[STATS] Total files size : %dGb\n", totalSize)

		getBiggestFiles(source)
	},
}

func init() {
	rootCmd.AddCommand(inboxCmd)
}

func getTotalSize(source string) int64 {
	var totalSize int64
	torrents := files.Find(source, ".torrent")
	for t := range torrents {
		files := files.Parse(torrents[t].File)
		for _, f := range files {
			totalSize = totalSize + f.Size
		}
	}
	totalSize = totalSize / 1073741824
	return totalSize
}

func getBiggestFiles(source string) []string {
	torrents := files.Find(source, ".torrent")
	for t := range torrents {
		files := files.Parse(torrents[t].File)
		for _, f := range files {
			log.Println(f)
		}
	}
	return nil
}
