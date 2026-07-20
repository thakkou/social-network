-- COMMENTS
CREATE TABLE IF NOT EXISTS COMMENTS (
    id         INTEGER  NOT NULL UNIQUE,
    user_id    INTEGER  NOT NULL,
    post_id    INTEGER  NOT NULL,
    created_at DATETIME NOT NULL,
    text       TEXT     NULL    ,
    PRIMARY KEY (id AUTOINCREMENT),
    FOREIGN KEY (user_id)
        REFERENCES USERS (id)
        ON DELETE CASCADE,
    FOREIGN KEY (post_id)
        REFERENCES POSTS (id)
        ON DELETE CASCADE 
);