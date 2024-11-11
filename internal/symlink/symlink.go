package symlink

import (
	"chignole/torlinks/internal/files"
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/viper"
)

func Create(file string, size int64, db *sql.DB) string {
	var targetPath string
	var targetSize int64

	torrentStatus := "unset"
	minSize := viper.GetInt64("options.minimalSize")
	symlinkDir := viper.GetString("general.symlinkDir")

	// Check file size and skip it if needed. Useful to ignore .nfo, .sfv, etc...
	if size < minSize {
		log.Printf("[INFO] Skipping file %v", file)
		return torrentStatus
	}

	// Count size matches - If it's more than 1, mark the torrent and skip it
	var count int
	err := db.QueryRow("SELECT COUNT (*) FROM files where size=?", size).Scan(&count)
	if err != nil {
		log.Fatalf("[ERROR] Failed to count the number of matchs %v\n", err)
	}

	if count > 1 {
		torrentStatus = "multi"
		return torrentStatus
	}

	// Prepare query
	stmt, err := db.Prepare("SELECT path, size FROM files where size =?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	// Query database for similar filesize - Returns matching file path and size
	err = stmt.QueryRow(size).Scan(&targetPath, &targetSize)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("[WARN] No file found in database for: %v\n", file)
			return torrentStatus
		} else {
			log.Fatalf("[ERROR] Failed to execute query. %v\n", err)
		}
	}

	// rows, err := stmt.Query(size)
	// for rows.Next() {
	// 	err := rows.Scan(&targetPath, &targetSize)
	// 	if err != nil {
	// 		log.Fatalf("[ERROR] Failed to scan row. %v\n", err)
	// 	}
	// 	log.Println(targetPath)
	// }

	// Check for subfolder - Creates it if needed
	checkFolder := regexp.MustCompile(`.*\/`)
	folder := checkFolder.FindStringSubmatch(file)
	if folder != nil {
		folderToCreate := filepath.Join(symlinkDir, folder[0])
		os.Mkdir(folderToCreate, 0744)
	}

	// Creates symlink
	file = filepath.Join(symlinkDir, file)
	os.Symlink(targetPath, file)
	torrentStatus = "linked"
	return torrentStatus
}

func Retry(source string) {
	failedFiles := files.Find(source, ".delete")
	for f := range failedFiles {
		trimmedName := strings.TrimSuffix(failedFiles[f].File, ".delete")
		os.Rename(failedFiles[f].File, trimmedName)
	}
}
