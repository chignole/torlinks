package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/j-muller/go-torrent-parser"
)

// TODO Dissociate failure and success in finding files for symlinks
// TODO Use viper for confguration

// Variables
var torrentFolder string
var showFolder string
var movieFolder string
var torrents []string
var files []string

type file struct {
	path        string
	size        int64
	showTitle   string
	showEpisode string
}

type configuration struct {
	General general `json:"general"`
}

type general struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
	TvShows     string `json:"tvshows"`
	Movies      string `json:"movies"`
}

// Search for specific file extension in a directory - Returning an array of files
func find(root, ext string) []string {
	var a []string
	filepath.WalkDir(root, func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			log.Fatalf("[ERROR] Error while looking for torrent files : %v", e)
		}
		if filepath.Ext(d.Name()) == ext {
			a = append(a, s)
		}
		return nil
	})
	return a
}

// Parse a torrent - Returns an array of file structs
func parse(torrent string) []file {
	var parsedFiles []file
	showPattern := regexp.MustCompile(`(.*)(?:\.S\d{2}.*)(S\d{2}E\d{2})`)
	a, err := gotorrentparser.ParseFromFile(torrent)
	if err != nil {
		log.Fatalf("[ERROR] Error while parsing torrent : %v", err)
	} else {
		for _, f := range a.Files {
			var showTitle string = "none"
			var showEpisode string = "none"
			path := filepath.Join(f.Path...)
			size := f.Length
			re := showPattern.FindStringSubmatch(path)
			if len(re) > 1 {
				showTitle = re[1]
				showEpisode = re[2]
			}
			parsedFile := file{path, size, showTitle, showEpisode}
			parsedFiles = append(parsedFiles, parsedFile)
		}
	}
	return parsedFiles
}

// Browse a folder, looking for a specific file, and  creates a symlink to this file
func createSymlink(f file, showFolder string) bool {
	var mediaFolder string
	var symLink bool = false
	// Check for subfolder
	checkFolder := regexp.MustCompile(`.*\/`)
	folder := checkFolder.FindStringSubmatch(f.path)
	if f.showTitle != "none" {
		mediaFolder = showFolder
	} else {
		mediaFolder = movieFolder
	}
	filepath.WalkDir(mediaFolder, func(a string, d fs.DirEntry, e error) error {
		if e != nil {
			log.Fatalf("[ERROR] Error while browsing media directory : %v", e)
		}
		fileStats, err := os.Stat(a)
		if err != nil {
			log.Fatalf("[ERROR] Error while getting filestats : %v", err)
		}
		if fileStats.Size() == f.size {
			if folder != nil {
				os.Mkdir(folder[0], 0744)
			}
			os.Symlink(a, f.path)
			log.Println("[PASS]", a, f.path)
			symLink = true
		}
		return nil
	})
	return symLink
}

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
	loadConfig := configuration{}
	err = json.Unmarshal(config, &loadConfig)
	if err != nil {
		log.Panicf("[ERROR] Error processing configuration file : %v", err)
	}
	return loadConfig
}

// Main function
func main() {
	// Logging setup
	f, err := os.OpenFile("torlinks.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("[ERROR] Error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	// Loading configuration file
	loadConfig := initConfig()
	torrentFolder = loadConfig.General.Source
	movieFolder = loadConfig.General.Movies
	showFolder = loadConfig.General.TvShows

	// Search for torrent files in specified folder
	torrents := find(torrentFolder, ".torrent")

	// Parse every torrent file, then browse media folder and create symlinks
	for _, t := range torrents {
		fmt.Println("Processing : ", t)
		filesToProcess := parse(t)

		for _, f := range filesToProcess {
			fmt.Println("Processing :", f)
			linkCreated := createSymlink(f, showFolder)
			fmt.Println(linkCreated)
			if linkCreated == true {
				n := loadConfig.General.Destination + t
				defer os.Rename(t, n)
			}
		}
	}
}
