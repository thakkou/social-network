CREATE TABLE IF NOT EXISTS GROUP_INVITES (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    group_id INTEGER NOT NULL,
    inviter_id INTEGER NOT NULL,
    invited_user_id INTEGER NOT NULL,
    status TEXT DEFAULT 'pending', -- pending accepted rejected
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY(group_id)
        REFERENCES GROUPS(id)
        ON DELETE CASCADE,
    FOREIGN KEY(inviter_id)
        REFERENCES USERS(id)
        ON DELETE CASCADE,
    FOREIGN KEY(invited_user_id)
        REFERENCES USERS(id)
        ON DELETE CASCADE
);