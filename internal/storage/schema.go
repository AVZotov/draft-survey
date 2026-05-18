package storage

type migration struct {
	version int
	sql     string
}

var migrations = []migration{
	{1, `CREATE TABLE surveys (
    	id TEXT PRIMARY KEY,
    	imo text,
    	created_at TEXT NOT NULL DEFAULT (datetime('now')),
    	data TEXT NOT NULL
	);`},
	{2, `CREATE TABLE users (
    id   INTEGER PRIMARY KEY CHECK (id = 1),
    data TEXT NOT NULL
	);`},
}
