/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"chignole/torlinks/internal/database"
	"log"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var searchTerm string

// dbSearchCmd represents the dbSearch command
var dbSearchCmd = &cobra.Command{
	Use:   "dbSearch",
	Short: "Search your database for specific files",
	Run: func(cmd *cobra.Command, args []string) {
		if searchTerm == "" {
			log.Println("[ERROR] You must specify a search pattern")
			cmd.Usage()
			os.Exit(1)
		}
		dataBaseFile := viper.GetString("general.database")
		db, err := database.OpenDatabase(dataBaseFile)
		if err != nil {
			log.Printf("[ERROR] Error opening database : %v", err)
		}
		searchByNameResults, err := database.SearchByName(db, searchTerm)
		if err != nil {
			log.Printf("[ERROR] Error searching by name")
		}

		// Display results table
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"#", "Filename"})
		for i := range searchByNameResults {
			t.AppendRow([]interface{}{i + 1, searchByNameResults[i]})
			t.AppendSeparator()
		}
		t.SetStyle(table.StyleColoredRedWhiteOnBlack)
		t.Style().Options.SeparateRows = true
		t.Render()
	},
}

func init() {
	dbSearchCmd.Flags().StringVarP(&searchTerm, "name", "n", "", "Search term")
	rootCmd.AddCommand(dbSearchCmd)
}
