-- EVENT RESPONSES (Going / Not Going)
CREATE TABLE IF NOT EXISTS EVENT_RESPONSES (
    event_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('going', 'not_going')),

    PRIMARY KEY (event_id, user_id),

    FOREIGN KEY(event_id)
        REFERENCES GROUP_EVENTS(id)
        ON DELETE CASCADE,
    FOREIGN KEY(user_id)
        REFERENCES USERS(id)
        ON DELETE CASCADE
);