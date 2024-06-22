package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func OpenDatabase(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	createTableSQL := `CREATE TABLE IF NOT EXISTS files (
    id INTEGER PRIMARY KEY,
    name TEXT,
    path TEXT UNIQUE,
    size INTERGER
  );`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func UpsertFile(db *sql.DB, name, path string, size int64) error {
	query := `INSERT INTO  files  (name, path, size) VALUES (?, ?, ?)
            ON CONFLICT (path) DO UPDATE SET name=excluded.name, size=excluded.size;`
	_, err := db.Exec(query, name, path, size)
	return err
}

func CleanupDatabase(db *sql.DB, filesMap map[string]struct{}) error {
	rows, err := db.Query("SELECT path FROM files")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var path string
		err := rows.Scan(&path)
		if err != nil {
			return err
		}
		if _, exists := filesMap[path]; !exists {
			_, err := db.Exec("DELETE FROM files WHERE path = ?", path)
			if err != nil {
				return err
			}
		}
	}
	return rows.Err()
}
