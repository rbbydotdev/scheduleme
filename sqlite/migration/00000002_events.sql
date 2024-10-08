CREATE TABLE events (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT DEFAULT '',
  duration INTEGER DEFAULT 60,
  avail_masks BLOB,
  user_id INTEGER NOT NULL,
  visible BOOLEAN DEFAULT TRUE,
  FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);