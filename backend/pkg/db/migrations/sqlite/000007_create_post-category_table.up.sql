---category post 
CREATE TABLE IF NOT EXISTS POST_CATEGORY (
    post_id INTEGER NOT NULL,
    category_id INTEGER NOT NULL,

    PRIMARY KEY (post_id, category_id),

    FOREIGN KEY (post_id)
        REFERENCES POSTS(id)
        ON DELETE CASCADE,
    FOREIGN KEY (category_id)
        REFERENCES CATEGORY(id)
        ON DELETE CASCADE
);