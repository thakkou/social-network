-- POSTS
CREATE TABLE IF NOT EXISTS POSTS (
    id         INTEGER  NOT NULL UNIQUE,
    user_id    INTEGER  NOT NULL,
    created_at DATETIME NOT NULL,
    title      TEXT     NULL,
    text       TEXT     NULL,
    image      TEXT     NULL,
    privacy    TEXT     NOT NULL DEFAULT 'public' CHECK (privacy IN ('public', 'almost_private', 'private')),
    PRIMARY KEY (id AUTOINCREMENT),
    FOREIGN KEY (user_id)
        REFERENCES USERS (id)
        ON DELETE CASCADE
);