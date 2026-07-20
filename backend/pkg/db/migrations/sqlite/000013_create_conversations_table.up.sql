CREATE TABLE IF NOT EXISTS CONVERSATIONS (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    user1_id INTEGER NOT NULL,
    user2_id INTEGER NOT NULL,

    last_message TEXT,
    last_message_at DATETIME,

    user1_last_read_message_id INTEGER,
    user2_last_read_message_id INTEGER,

    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    UNIQUE(user1_id, user2_id),

    FOREIGN KEY (user1_id)
        REFERENCES USERS(id)
        ON DELETE CASCADE,
    FOREIGN KEY (user2_id)
        REFERENCES USERS(id)
        ON DELETE CASCADE
);