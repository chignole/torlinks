package symlink

import (
	"chignole/torlinks/internal/files"
	"database/sql"
	"log"
	"os"
	"regexp"
	"strings"
)

func Create(file string, size int64, db *sql.DB) bool {
	var targetPath string
	var targetSize int64

	linked := false

	if size < 2000000 {
		log.Printf("[INFO] Skipping file %v. Size < 2MB\n", file)
	}

	stmt, err := db.Prepare("SELECT path, size FROM files where size =?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(size).Scan(&targetPath, &targetSize)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("[WARN] No file found in database for: %v\n", file)
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

func Retry(source string) {
	failedFiles := files.Find(source, ".delete")
	for f := range failedFiles {
		trimmedName := strings.TrimSuffix(failedFiles[f].File, ".delete")
		os.Rename(failedFiles[f].File, trimmedName)
	}
}
