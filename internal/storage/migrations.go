package storage

import (
	"database/sql"
	"time"
)

func runMigrations(db *sql.DB) error {
	const createMigrations = `CREATE TABLE IF NOT EXISTS migrations (
    		version INTEGER PRIMARY KEY,
    		applied_at TEXT NOT NULL DEFAULT (datetime('now'))
		);`

	_, err := db.Exec(createMigrations)
	if err != nil {
		return err
	}

	var version *int
	const getLastMigration = `SELECT max(version) FROM migrations;`
	row := db.QueryRow(getLastMigration)

	err = row.Scan(&version)
	if err != nil {
		return err
	}

	if err = updateMigrations(version, db); err != nil {
		return err
	}

	return nil
}

func updateMigrations(version *int, db *sql.DB) error {
	for _, m := range migrations {
		if version != nil && m.version <= *version {
			continue
		}

		if err := func() error {
			tx, err := db.Begin()
			if err != nil {
				//TODO: Add log if error or fatal
				return err
			}
			defer tx.Rollback()

			_, err = tx.Exec(`INSERT INTO migrations (version, applied_at) VALUES (?, ?);`, m.version, time.Now().Format("2006-01-02 15:04:05"))
			if err != nil {
				return err
			}

			_, err = tx.Exec(m.sql)
			if err != nil {
				return err
			}

			return tx.Commit()
		}(); err != nil {
			return err
		}

	}

	return nil
}
