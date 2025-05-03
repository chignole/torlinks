package cmd

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "torlinks",
	Short: "A program that creates symlinks to torrent files content to facilitate continued sharing.",
	Long: `
 ______  ____    ___    __    ____   _  __   __ __   ____
/_  __/ / __ \  / _ \  / /   /  _/  / |/ /  / //_/  / __/
 / /   / /_/ / / , _/ / /__ _/ /   /    /  / ,<    _\ \  
/_/    \____/ /_/|_| /____//___/  /_/|_/  /_/|_|  /___/  
                                                         

This program generates virtual links to the content of torrent files, enabling users to continue
sharing these files seamlessly. By creating these links, the program allows for the efficient 
distribution and access of shared files without the need for re-downloading or manually managing 
the original torrent files. This helps maintain the availability and integrity of shared content 
across multiple users and platforms.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.torlinks.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.

	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringP("config", "c", "", "config file (default is $HOME/.config/torlinks/config.yaml")
}

func initConfig() {
	log.Println("[INFO] Initizalizing configuration file ...")
	cfgFile := rootCmd.PersistentFlags().Lookup("config").Value.String()
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Printf("[ERROR] Error while getting home directory :%v\n", err)
		}
		configPath := filepath.Join(home, ".config", "torlinks")
		viper.AddConfigPath(configPath)
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("[ERROR] Error while reading config file : %v\n", err)
	}
}
