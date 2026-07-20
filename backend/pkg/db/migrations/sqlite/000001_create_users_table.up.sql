-- USERS
CREATE TABLE IF NOT EXISTS USERS (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    firstname TEXT NOT NULL,
    lastname TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL, -- for oauth
    birthdate TEXT NOT NULL,
    nickname TEXT UNIQUE, -- optional
    aboutme TEXT, -- optional (may need larger text !)
    avatar TEXT, -- nullable, stores "/uploads/avatars/xxx.png"
    is_private INTEGER NOT NULL DEFAULT 0, -- 0->public/1-private
    
    -- not used
    last_seen DATETIME
);

-- + index
CREATE INDEX IF NOT EXISTS idx_username  ON users(nickname COLLATE NOCASE);