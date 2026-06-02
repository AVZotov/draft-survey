package storage

import (
	"database/sql"
	"embed"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func NewDB(path string, fs embed.FS) (*sql.DB, error) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	if err = runMigrations(db); err != nil {
		return nil, err
	}

	if err = seedDictionaries(db, fs); err != nil {
		return nil, err
	}

	return db, nil
}
