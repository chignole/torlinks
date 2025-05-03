package cmd

import (
	"chignole/torlinks/internal/files"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/hekmon/transmissionrpc"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/exp/slices"
)

// seededCmd represents the seeded command
var seededCmd = &cobra.Command{
	Use:   "seeded",
	Short: "Verify which torrent files are currently not seeded by the Transmission client",
	Long:  `Verify which torrent files are currently not seeded by the Transmission client`,
	Run: func(cmd *cobra.Command, args []string) {
		torrentsInbox := viper.GetString("general.torrentsInbox")
		torrentsWatchDir := viper.GetString("general.torrentsWatchDir")
		transmissionServer := viper.GetString("transmission.server")
		transmissionUser := viper.GetString("transmission.user")
		transmissionPassword := viper.GetString("transmission.password")

		// Connect to the transmission client
		client, err := transmissionrpc.New(transmissionServer, transmissionUser, transmissionPassword, nil)
		if err != nil {
			log.Fatalf("[ERROR] Failed to connect to transmlission : %v", err)
		}

		// Get data from all torrents and put all hashes in a slice nammed seededTorrentHashes
		fields := []string{"id", "name", "status", "hashString"}
		torrents, err := client.TorrentGet(fields, nil)
		if err != nil {
			log.Fatalf("[ERROR] Failed to get torrents : %v", err)
		}

		var seededTorrentsHashes []string

		for _, torrent := range torrents {
			seededTorrentsHashes = append(seededTorrentsHashes, *torrent.HashString)
		}

		// Browse torrent folder looking for every added torrent, check if present in hashes slices
		// if not, then it's not actually seeded
		torrentfiles := files.Find(torrentsWatchDir, ".added")
		for _, torrent := range torrentfiles {
			torrentDetails := files.ParseDetails(torrent.File)
			if !slices.Contains(seededTorrentsHashes, torrentDetails.Hash) {
				newFileName := strings.TrimSuffix(filepath.Base(torrentDetails.File), ".added")
				newFilePath := filepath.Join(torrentsInbox, newFileName)
				log.Printf("\033[1;32m[INFO] Moved : %v\033[0m\n", filepath.Base(newFilePath))
				err := os.Rename(torrentDetails.File, newFilePath)
				if err != nil {
					log.Printf("[ERROR] Failed to move torrent files : %v", err)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(seededCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// seededCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// seededCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
