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