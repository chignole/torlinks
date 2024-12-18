package files

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"time"

	gotorrentparser "github.com/j-muller/go-torrent-parser"
	"github.com/spf13/viper"
)

type torrent struct {
	File   string
	Status string
}

type file struct {
	Path string
	Size int64
}

type torrentDetails struct {
	File      string
	Announce  []string
	Comment   string
	CreatedAt time.Time
	CreatedBy string
	Hash      string
}

// Search for specific file extension in a directory - Returning an array of torrent constructs
func Find(root, ext string) []torrent {
	var a []torrent
	filepath.WalkDir(root, func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			log.Fatalf("[ERROR] Error while looking for torrent files : %v", e)
		}
		if d.IsDir() && s != root {
			return filepath.SkipDir
		}
		if filepath.Ext(d.Name()) == ext {
			t := torrent{s, "unset"}
			a = append(a, t)
		}
		return nil
	})
	return a
}

// Parse a torrent - Returns an array of file structs
func Parse(torrent string) []file {
	var parsedFiles []file
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

// Parse a torrent and return some useful informations
func ParseDetails(torrent string) torrentDetails {
	var parsedDetails torrentDetails
	a, err := gotorrentparser.ParseFromFile(torrent)
	if err != nil {
		log.Fatalf("[ERROR] Error while parsing torrent : %v", err)
	}
	file := torrent
	announce := a.Announce
	comment := a.Comment
	createdAt := a.CreatedAt
	createdBy := a.CreatedBy
	hash := a.InfoHash

	parsedDetails = torrentDetails{file, announce, comment, createdAt, createdBy, hash}
	return parsedDetails
}

func Clean(t torrent) {
	torrentsWatchDir := viper.GetString("general.torrentsWatchDir")
	if t.Status == "linked" {
		fileName := filepath.Base(t.File)
		n := torrentsWatchDir + fileName
		log.Println("[INFO] Moving", fileName, "to", torrentsWatchDir)
		defer os.Rename(t.File, n)
	} else if t.Status == "multi" {
		n := t.File + ".multi"
		defer os.Rename(t.File, n)
	} else {
		n := t.File + ".delete"
		defer os.Rename(t.File, n)
	}
}
