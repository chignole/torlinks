package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/j-muller/go-torrent-parser"
)

// Options
// const torrentFolder = "./"
// const showFolder = "/mnt/medias.1/Series/"

// Variables
var torrents []string
var files []string

type file struct {
	path        string
	size        int64
	showTitle   string
	showEpisode string
}

// Search for specific file extension in a directory - Returning an array of files
func find(root, ext string) []string {
	var a []string
	filepath.WalkDir(root, func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
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
	a, err := gotorrentparser.ParseFromFile(torrent)
	showPattern := regexp.MustCompile(`(.*)(?:\.S\d{2}.*)(S\d{2}E\d{2})`)
	if err != nil {
		fmt.Println(err)
	} else {
		for _, f := range a.Files {
			path := filepath.Join(f.Path...)
			size := f.Length
			re := showPattern.FindStringSubmatch(path)
			showTitle := re[1]
			showEpisode := re[2]
			parsedFile := file{path, size, showTitle, showEpisode}
			parsedFiles = append(parsedFiles, parsedFile)
		}
	}
	return parsedFiles
}

// Browse a folder, looking for a specific file, and  creates a symlink to this file
func createSymlink(f file, showFolder string) {
	checkFolder := regexp.MustCompile(`.*\/`)
	folder := checkFolder.FindStringSubmatch(f.path)
	if folder != nil {
		os.Mkdir(folder[0], 0744)
	}
	filepath.WalkDir(showFolder, func(a string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if strings.Contains(a, f.showTitle) {
			fileStats, err := os.Stat(a)
			if err != nil {
				return err
			}
			if fileStats.Size() == f.size {
				// fmt.Println(f.path, "->", a)
				os.Symlink(a, f.path)
			}
		}
		return nil
	})
}

func main() {
	torrentFolder := os.Args[1]
	showFolder := os.Args[2]

	torrents := find(torrentFolder, ".torrent")

	for _, t := range torrents {
		fmt.Println("Processing : ", t)
		filesToProcess := parse(t)

		for _, f := range filesToProcess {
			fmt.Println("Processing :", f)
			createSymlink(f, showFolder)
		}
	}
}
