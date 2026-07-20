-- For custom private posts (specific followers allowed to see it)
CREATE TABLE IF NOT EXISTS POST_ALLOWED_USERS (
    post_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    PRIMARY KEY (post_id, user_id),
    FOREIGN KEY (post_id)
        REFERENCES POSTS(id)
        ON DELETE CASCADE,
    FOREIGN KEY (user_id)
        REFERENCES USERS(id)
        ON DELETE CASCADE
);