-- POST REACTIONS
-- Fix POST_REACTIONS
CREATE TABLE IF NOT EXISTS POST_REACTIONS (
  user_id INTEGER NOT NULL,
  post_id INTEGER NOT NULL,
  is_like INTEGER NOT NULL CHECK (is_like IN (-1, 1)),

  PRIMARY KEY (user_id, post_id), -- <--- Add this

  FOREIGN KEY (user_id)
    REFERENCES USERS(id)
    ON DELETE CASCADE,
  FOREIGN KEY (post_id)
    REFERENCES POSTS(id)
    ON DELETE CASCADE 
);
-- reactions are unique by combination of both user_id and post_id