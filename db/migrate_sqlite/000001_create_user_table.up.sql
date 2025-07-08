CREATE TABLE IF NOT EXISTS "user"
(
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    username   TEXT NOT NULL UNIQUE,
    password   TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT (datetime(current_timestamp, 'localtime')),
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);