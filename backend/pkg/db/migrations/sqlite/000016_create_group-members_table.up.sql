CREATE TABLE IF NOT EXISTS GROUP_MEMBERS (
    group_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    role TEXT DEFAULT 'member', -- member | admin
    joined_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY(group_id,user_id),

    FOREIGN KEY(group_id)
        REFERENCES GROUPS(id)
        ON DELETE CASCADE,
    FOREIGN KEY(user_id)
        REFERENCES USERS(id)
        ON DELETE CASCADE
);