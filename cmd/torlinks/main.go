package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"

	// "chignole/torlinks/internal/database"
	"chignole/torlinks/internal/filescanner"

	"github.com/j-muller/go-torrent-parser"
)

// TODO Dissociate failure and success in finding files for symlinks
// TODO Use viper for confguration

// Variables
var torrentFolder string
var dataFolder []string
var updateDatabaseOnStart string
var torrents []string
var files []string

type file struct {
	path string
	size int64
}

type torrent struct {
	file   string
	linked bool
}

type configuration struct {
	General general `json:"general"`
}

type general struct {
	Source                string   `json:"source"`
	Destination           string   `json:"destination"`
	Data                  []string `json:"data"`
	UpdateDatabaseOnStart string   `json:"updatedatabaseonstart"`
}

var loadConfig configuration

// Processing configuration file
func initConfig() configuration {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Panic(err)
	}
	file := homeDir + "/.config/torlinks/config.json"
	configFile, err := os.Open(file)
	if err != nil {
		log.Fatalf("[ERROR] Error opening configuration file : %v", err)
	}
	defer configFile.Close()
	config, err := ioutil.ReadAll(configFile)
	if err != nil {
		log.Panicf("[ERROR] Error processing configuration file : %v", err)
	}
	err = json.Unmarshal(config, &loadConfig)
	if err != nil {
		log.Panicf("[ERROR] Error processing configuration file : %v", err)
	}
	return loadConfig
}

// Search for specific file extension in a directory - Returning an array of torrent constructs
func find(root, ext string) []torrent {
	var a []torrent
	filepath.WalkDir(root, func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			log.Fatalf("[ERROR] Error while looking for torrent files : %v", e)
		}
		if d.IsDir() && s != root {
			return filepath.SkipDir
		}
		if filepath.Ext(d.Name()) == ext {
			t := torrent{s, false}
			a = append(a, t)
		}
		return nil
	})
	return a
}

// Parse a torrent - Returns an array of file structs
func parseTorrent(torrent string) []file {
	var parsedFiles []file
	// TODO Modify regex so it can accept sXXeXX or SXXEXX
	a, err := gotorrentparser.ParseFromFile(torrent)
	if err != nil {
		log.Fatalf("[ERROR] Error while parsing torrent : %v", err)
	} else {
		for _, f := range a.Files {
			path := filepath.Join(f.Path...)
			size := f.Length
			parsedFile := file{path, size}
			parsedFiles = append(parsedFiles, parsedFile)
		}
	}
	return parsedFiles
}

func clean(t torrent) {
	if t.linked == true {
		n := loadConfig.General.Destination + t.file
		fmt.Println("[INFO] Moving", t.file, "to", loadConfig.General.Destination)
		log.Println("[INFO] Moving", t.file, "to", loadConfig.General.Destination)
		defer os.Rename(t.file, n)
	} else {
		n := t.file + ".delete"
		defer os.Rename(t.file, n)
	}
}

func createSymlink(file string, size int64, db *sql.DB) bool {
	var targetPath string
	var targetSize int64

	linked := false

	if size < 2000000 {
		fmt.Printf("[INFO] Skipping file %v. Size < 2MB\n", file)
	}

	stmt, err := db.Prepare("SELECT path, size FROM files where size =?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(size).Scan(&targetPath, &targetSize)

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("[WARN] No file found in database for: %v\n", file)
			return linked
		} else {
			log.Fatalf("[ERROR] Failed to execute query. %v\n", err)
		}
	}
	// Check for subfolder
	checkFolder := regexp.MustCompile(`.*\/`)
	folder := checkFolder.FindStringSubmatch(file)
	if folder != nil {
		os.Mkdir(folder[0], 0744)
	}

	os.Symlink(targetPath, file)
	linked = true
	return linked
}

// Main function
func main() {
	// Loading configuration file
	log.Println("[INFO] Loading configuration file...")
	loadConfig := initConfig()
	torrentFolder = loadConfig.General.Source
	dataFolder = loadConfig.General.Data
	updateDatabaseOnStart = loadConfig.General.UpdateDatabaseOnStart

	// Logging setup
	f, err := os.OpenFile("torlinks.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("[ERROR] Error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	// Opendatabase
	log.Println("[INFO] Opening database ./files.db")
	db, err := sql.Open("sqlite3", "./files.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Update files database
	if updateDatabaseOnStart == "true" {
		log.Println("[INFO] Starting database update")
		directoriesToScan := dataFolder
		filescanner.ScanDirectories(directoriesToScan[:], "./files.db")
	}

	// Search for torrents in specified folder
	torrents := find(torrentFolder, ".torrent")

	for t := range torrents {
		fmt.Printf("[INFO] Processing : %v\n", torrents[t].file)
		filesToProcess := parseTorrent(torrents[t].file)
		for _, f := range filesToProcess {
			linked := createSymlink(f.path, f.size, db)
			if linked == true {
				torrents[t].linked = true
			}
		}
	}

	// Clean torrent files after procesing
	for _, t := range torrents {
		clean(t)
	}
}
