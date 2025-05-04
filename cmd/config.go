package cmd

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Creates a configuration file.",
	Long:  `Creates a configuration file. Default path is $HOME/.config/torlinks/config.yaml`,
	Run: func(cmd *cobra.Command, args []string) {
		configExists := checkConfig()
		if configExists == true {
			log.Println("[INFO] Configuration file already exists.")
		} else {
			createConfig()
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}

// Check if config file  already exists
func checkConfig() bool {
	homeDir := os.Getenv("HOME")
	configPath := filepath.Join(homeDir, ".config", "torlinks", "config.yaml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func createConfig() {
	homeDir := os.Getenv("HOME")
	configPath := filepath.Join(homeDir, ".config", "torlinks", "config.yaml")
	err := os.MkdirAll(filepath.Dir(configPath), 0755)
	if err != nil {
		log.Printf("[ERROR] Error while creating directories : %v\n", err)
		return
	}
	err = os.WriteFile(configPath, []byte(strings.TrimSpace(defaultConfig)), 0644)
	if err != nil {
		log.Printf("[ERROR] Error while creating config file : %v\n", err)
		return
	}
	log.Printf("[INFO] Default configuration file created : %v\n", configPath)
}

// Default configuration file
var defaultConfig string = `
general:
  torrentsInbox: ""
  torrentsWatchDir: ""
  symlinkDir : ""
  data:
    - 
    -
    -

options:
  # Minimal size of kb size to process 
  minimalSize: 2000000

  # Ratio (%) of successful matches needed to consider a torrent seedable
  minimalMatch: 95

  # Default value for the processMultipleMatches flag
  # If set to true, will prompt user for match disambiguation
  processMultipleMatches: false

database:
  file: ""
  ping: ""

metrics: 
  file: ""
  ping: ""

transmission:
  server: ""
  user: ""
  password: ""
`
