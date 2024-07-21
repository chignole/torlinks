/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"log"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	// "github.com/jedib0t/go-pretty/v6/text"

	"chignole/torlinks/internal/files"
	"github.com/spf13/cobra"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check [file]",
	Short: "Display the content of specified torrent file.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filename := args[0]

		details := files.ParseDetails(filename)
		files := files.Parse(filename)

		log.Println(details)
		// log.Println(files)

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"#", "Filename", "Size"})

		for i, file := range files {
			t.AppendRow([]interface{}{i + 1, file.Path, file.Size})
			t.AppendSeparator()
		}

		t.SetStyle(table.StyleColoredRedWhiteOnBlack)
		t.Style().Options.SeparateRows = true

		t.Render()

	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
}
