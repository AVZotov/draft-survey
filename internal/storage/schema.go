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
	{3, `CREATE TABLE countries (
    code TEXT PRIMARY KEY,
    name TEXT NOT NULL
	);`},
	{4, `CREATE TABLE ports (
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    locode       TEXT,
    name         TEXT NOT NULL,
    country_code TEXT NOT NULL,
    coordinates  TEXT
	);
	CREATE INDEX idx_ports_name ON ports(name);
	CREATE INDEX idx_ports_country ON ports(country_code);`},
	{5, `ALTER TABLE users ADD COLUMN signature BLOB;`},
	{6, `CREATE TABLE cargo_types (
    id   INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE);`},
	{7, `CREATE TABLE packing (
    id   INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE);`},
}
