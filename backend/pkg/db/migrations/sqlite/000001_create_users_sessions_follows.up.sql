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

CREATE INDEX IF NOT EXISTS idx_username ON users(nickname COLLATE NOCASE);

-- SESSIONS
CREATE TABLE IF NOT EXISTS SESSIONS (
    id TEXT PRIMARY KEY UNIQUE, -- uuid
    expires_at DATETIME NOT NULL,
    user_id INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- FOLLOWS
CREATE TABLE IF NOT EXISTS FOLLOWS (
    follower_id INTEGER NOT NULL,
    following_id INTEGER NOT NULL,

    status TEXT NOT NULL DEFAULT 'pending',
    -- pending | accepted
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (follower_id, following_id),

    FOREIGN KEY (follower_id)
        REFERENCES USERS(id)
        ON DELETE CASCADE,

    FOREIGN KEY (following_id)
        REFERENCES USERS(id)
        ON DELETE CASCADE,

    CHECK (follower_id != following_id)
);