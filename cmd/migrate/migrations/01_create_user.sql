CREATE TABLE IF NOT EXISTS users (
	id TEXT PRIMARY KEY,
	name TEXT,
	email TEXT UNIQUE,
	password BLOB,
	created_at TEXT
);