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
		totalFailedFiles := getFailedFiles(source)
		totalInboxFiles := getInboxFiles(source)

		log.Printf("[STATS] Total files size....: %dGb\n", totalSize)
		log.Printf("[STATS] Total failed files..: %d\n", totalFailedFiles)
		log.Printf("[STATS] Total inbox files...: %d\n", totalInboxFiles)
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

func getInboxFiles(source string) int64 {
	inboxFiles := files.Find(source, ".torrent")
	totalInboxFiles := len(inboxFiles)
	return int64(totalInboxFiles)
}

func getFailedFiles(source string) int64 {
	failedFiles := files.Find(source, ".delete")
	totalFailedFiles := len(failedFiles)
	return int64(totalFailedFiles)
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
