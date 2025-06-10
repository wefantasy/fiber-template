CREATE TABLE "user" (
	id INTEGER PRIMARY KEY NOT NULL,
	username TEXT UNIQUE,
	password TEXT,
	created_at DATETIME DEFAULT (datetime(current_timestamp, 'localtime')),
	updated_at DATETIME,
	deleted_at DATETIME
);