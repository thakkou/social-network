--group request
CREATE TABLE IF NOT EXISTS GROUP_REQUESTS (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    group_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    status TEXT DEFAULT 'pending',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    UNIQUE(group_id,user_id),

    FOREIGN KEY(group_id)
        REFERENCES GROUPS(id)
        ON DELETE CASCADE,
    FOREIGN KEY(user_id)
        REFERENCES USERS(id)
        ON DELETE CASCADE
);