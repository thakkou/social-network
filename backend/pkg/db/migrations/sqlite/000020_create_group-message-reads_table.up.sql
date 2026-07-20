CREATE TABLE IF NOT EXISTS GROUP_MESSAGE_READS (
    group_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    last_read_message_id INTEGER,

    PRIMARY KEY (group_id, user_id),

    FOREIGN KEY (group_id)
        REFERENCES GROUPS(id)
        ON DELETE CASCADE,
    FOREIGN KEY (user_id)
        REFERENCES USERS(id)
        ON DELETE CASCADE,
    FOREIGN KEY (last_read_message_id)
        REFERENCES GROUP_MESSAGES(id)
        ON DELETE SET NULL
);