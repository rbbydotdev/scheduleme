CREATE TABLE auths (
	id            INTEGER PRIMARY KEY AUTOINCREMENT
	, user_id       INTEGER NOT NULL REFERENCES users (id) ON DELETE CASCADE
	, source        TEXT NOT NULL
	, source_id     TEXT NOT NULL
	, access_token  TEXT NOT NULL
	, refresh_token TEXT NOT NULL
	, expiry        DATETIME NOT NULL
	, created_at    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	, updated_at    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	, UNIQUE(user_id, source)
	, UNIQUE(source, source_id) 
);