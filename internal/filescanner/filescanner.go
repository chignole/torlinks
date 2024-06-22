package filescanner

import (
	"chignole/torlinks/internal/database"
	// "database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func ScanDirectories(dirs []string, dbPath string) error {
	db, err := database.OpenDatabase(dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	filesMap := make(map[string]struct{})

	for _, dir := range dirs {
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && isValidExtension(info.Name()) {
				fmt.Printf("[SCAN] %v \n", path)
				filesMap[path] = struct{}{}
				err := database.UpsertFile(db, info.Name(), path, info.Size())
				if err != nil {
					log.Printf("Error inserting/updating file %v", err)
				}
			}
			return nil
		})
		if err != nil {
			log.Printf("Error walking the path %q, %v\n", dir, err)
		}
	}

	err = database.CleanupDatabase(db, filesMap)
	if err != nil {
		log.Printf("Error cleaning up database: %v", err)
	}
	return nil
}

func isValidExtension(filename string) bool {
	validExtensions := []string{".mkv", ".mp4", ".wmv", ".avi"}
	for _, ext := range validExtensions {
		if strings.HasSuffix(strings.ToLower(filename), ext) {
			return true
		}
	}
	return false
}
