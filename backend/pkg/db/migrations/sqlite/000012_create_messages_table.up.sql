--this for private messages
CREATE TABLE IF NOT EXISTS MESSAGES (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    conversation_id INTEGER NOT NULL,
    sender_id INTEGER NOT NULL,
    text TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%f', 'now')),
    is_read INTEGER NOT NULL DEFAULT 0,

    FOREIGN KEY (conversation_id)
        REFERENCES CONVERSATIONS(id)
        ON DELETE CASCADE,
    FOREIGN KEY (sender_id)
        REFERENCES USERS(id)
        ON DELETE CASCADE
);