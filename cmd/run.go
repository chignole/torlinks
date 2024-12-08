/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"chignole/torlinks/internal/files"
	"chignole/torlinks/internal/utils"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type matchingFile struct {
	ID   int
	Name string
	Path string
	Size int64
}

var processMultipleMatches bool

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Manual processing of torrents matching multiple files",
	Long:  "Manual processing of torrents matching multiple files",
	Run: func(cmd *cobra.Command, args []string) {
		dbFile := viper.GetString("database.file")
		torrentsInbox := viper.GetString("general.torrentsInbox")
		torrentsWatchDir := viper.GetString("general.torrentsWatchDir")
		minSize := viper.GetInt64("options.minimalSize")
		minimalMatch := viper.GetFloat64("options.minimalMatch")
		processMultipleMatches := viper.GetBool("options.processMultipleMatches")

		torrentFiles := files.Find(torrentsInbox, ".torrent")

		// Open database
		db, err := sql.Open("sqlite3", dbFile)
		if err != nil {
			log.Printf("[ERROR] Error while opening database file : %v", err)
		}
		defer db.Close()

	torrentFilesLoop:
		for t := range torrentFiles {
			// log.Printf("[-] Processing : %v\n", torrentFiles[t].File)
			filesToProcess := files.Parse(torrentFiles[t].File)
			matchedSize := 0
			totalSize := 0

			for _, f := range filesToProcess {

				if f.Size < minSize {
					totalSize = totalSize + int(f.Size)
					continue
				}
				files, err := getFilesBySize(db, f.Size)
				if err != nil {
					log.Fatalf("[ERROR] Failed getting matching files %v\n", err)
				}

				log.Printf("\033[3;35m[-] Processing: %s \033[0m\n", f.Path)
				displayMatchingFiles(files)
				numberOfMatchingFiles := len(files)
				switch {
				case numberOfMatchingFiles == 0:
					totalSize = totalSize + int(f.Size)
				case numberOfMatchingFiles == 1:
					createSymlink(f.Path, files[0].Path)
					totalSize = totalSize + int(f.Size)
					matchedSize = matchedSize + int(f.Size)
				case numberOfMatchingFiles > 1:
					if !processMultipleMatches {
						log.Println("\033[33m[-] Multiple matches - Skipping torrent file\033[0m")
						continue torrentFilesLoop
					}
					choice, err := selectMatchingFile(len(files))
					if err != nil {
						log.Printf("[ERROR] Failed to select matching file %v\n", err)
					}
					createSymlink(f.Path, files[choice].Path)
					totalSize = totalSize + int(f.Size)
					matchedSize = matchedSize + int(f.Size)
				}
			}

			matchedPercentage := utils.CalculatePercentage(float64(matchedSize), float64(totalSize))
			switch {
			case matchedPercentage < minimalMatch:
				log.Printf("\033[1;31m[-] Failed with : %.2f%% matched files\033[0m\n", matchedPercentage)
				n := torrentFiles[t].File + ".delete"
				defer os.Rename(torrentFiles[t].File, n)
			case matchedPercentage >= minimalMatch:
				log.Printf("\033[1;32m[-] Passed with : %.2f%% matched files\033[0m\n", matchedPercentage)
				filename := filepath.Base(torrentFiles[t].File)
				n := torrentsWatchDir + filename
				defer os.Rename(torrentFiles[t].File, n)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.PersistentFlags().BoolVarP(&processMultipleMatches, "processMultipleMatches", "m", false, "Manually process multiple matches")
	viper.BindPFlag("options.processMultipleMatches", runCmd.PersistentFlags().Lookup("processMultipleMatches"))
}

// Creates symlink to matched data - Handles subdirectories creation if needed
// link = needed Symlink - target = actual data
func createSymlink(link string, target string) {
	symlinkDir := viper.GetString("general.symlinkDir")
	checkFolder := regexp.MustCompile(`.*\/`)
	folder := checkFolder.FindStringSubmatch(link)
	if folder != nil {
		folderToCreate := filepath.Join(symlinkDir, folder[0])
		os.Mkdir(folderToCreate, 0744)
	}
	link = filepath.Join(symlinkDir, link)
	os.Symlink(target, link)
}

// Retrieves a list of files from the database, matching the specified size
func getFilesBySize(db *sql.DB, size int64) ([]matchingFile, error) {
	rows, err := db.Query("SELECT id, path, name, size FROM files where size = ?", size)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matchingFiles []matchingFile
	for rows.Next() {
		var file matchingFile
		if err := rows.Scan(&file.ID, &file.Path, &file.Name, &file.Size); err != nil {
			return nil, err
		}
		matchingFiles = append(matchingFiles, file)
	}
	return matchingFiles, nil
}

// Display a list of files to the user letting him choose the right match
func selectMatchingFile(max int) (int, error) {
	for {
		fmt.Print("\033[1;33m>>> Select matching file: \033[0m")
		var input string
		fmt.Scanln(&input)

		choice, err := strconv.Atoi(input)
		if err == nil && choice >= 0 && choice < max {
			return choice, nil
		}
		fmt.Println("\033[1;31mInvalid input. Please enter a valid number.\033[0m")
	}
}

// Display a list of matching files
// TODO : github.com/lithammer/fuzzysearch/fuzzy ?
func displayMatchingFiles(files []matchingFile) {
	for i, file := range files {
		log.Printf("[%d] %s (ID: %d, Size: %d)\n", i, file.Name, file.ID, file.Size)
	}
}
