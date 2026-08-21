package storage

import (
	"database/sql"
	"embed"
	"os"
	"path/filepath"
	"strings"

	_ "modernc.org/sqlite"
)

func NewDB(path string, fs embed.FS) (*sql.DB, error) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite", withPragmas(path))
	if err != nil {
		return nil, err
	}

	// SQLite allows only one writer at a time no matter the journal mode.
	// Pinning the pool to a single physical connection makes Go's pool queue
	// concurrent callers instead of racing separate connections into
	// SQLITE_BUSY (autosave bursts from the UI fire several writes at once).
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(0) // never recycle the single connection

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

// withPragmas appends the WAL + busy_timeout + synchronous pragmas as
// modernc.org/sqlite DSN params instead of running them as one-off PRAGMA
// statements after Open. The driver re-applies every "_pragma" param on each
// new physical connection it opens (not just the first), which matters
// because busy_timeout and synchronous are per-connection settings — a plain
// db.Exec("PRAGMA busy_timeout=...") right after Open would only cover
// whichever connection happened to be live at that moment. journal_mode=WAL
// is persisted in the database file itself, but setting it here too keeps
// all three pragmas in one place and guarantees it for a brand-new file.
// synchronous=NORMAL is the mode SQLite's own docs recommend pairing with
// WAL: fsyncs on checkpoint instead of every commit, safe against process
// crashes (still durable across app-level failures), only a lost OS-level
// power failure could roll back the last few commits — an acceptable
// tradeoff for an offline desktop tool.
func withPragmas(path string) string {
	sep := "?"
	if strings.Contains(path, "?") {
		sep = "&"
	}
	return path + sep + "_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)"
}
