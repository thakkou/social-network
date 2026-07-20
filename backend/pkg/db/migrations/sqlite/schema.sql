-- -- USERS
-- CREATE TABLE IF NOT EXISTS USERS (
--     id INTEGER PRIMARY KEY AUTOINCREMENT,
--     created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

--     firstname TEXT NOT NULL,
--     lastname TEXT NOT NULL,
--     email TEXT NOT NULL UNIQUE,
--     password TEXT NOT NULL, -- for oauth
--     birthdate TEXT NOT NULL,
--     nickname TEXT UNIQUE, -- optional
--     aboutme TEXT, -- optional (may need larger text !)
--     avatar TEXT, -- nullable, stores "/uploads/avatars/xxx.png"
--     is_private INTEGER NOT NULL DEFAULT 0, -- 0->public/1-private
--     -- not used
--     last_seen DATETIME
-- );


-- CREATE INDEX IF NOT EXISTS idx_username  ON users(nickname COLLATE NOCASE);
-- -- SESSIONS
-- CREATE TABLE IF NOT EXISTS SESSIONS (
--     id TEXT PRIMARY KEY UNIQUE, -- uuid
--     expires_at DATETIME NOT NULL,
--     user_id INTEGER NOT NULL,
--     FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
-- );

-- -- POSTS
-- CREATE TABLE IF NOT EXISTS POSTS (
--     id         INTEGER  NOT NULL UNIQUE,
--     user_id    INTEGER  NOT NULL,
--     created_at DATETIME NOT NULL,
--     title      TEXT     NULL,
--     text       TEXT     NULL,
--     image      TEXT     NULL,
--     privacy    TEXT     NOT NULL DEFAULT 'public' CHECK (privacy IN ('public', 'almost_private', 'private')),
--     PRIMARY KEY (id AUTOINCREMENT),
--     FOREIGN KEY (user_id) REFERENCES USERS (id) ON DELETE CASCADE
-- );
-- -- For custom private posts (specific followers allowed to see it)
-- CREATE TABLE IF NOT EXISTS POST_ALLOWED_USERS (
--     post_id INTEGER NOT NULL,
--     user_id INTEGER NOT NULL,
--     PRIMARY KEY (post_id, user_id),
--     FOREIGN KEY (post_id) REFERENCES POSTS(id) ON DELETE CASCADE,
--     FOREIGN KEY (user_id) REFERENCES USERS(id) ON DELETE CASCADE
-- );

-- ---CATEGORY 
-- CREATE TABLE IF NOT EXISTS CATEGORY (
--     id INTEGER PRIMARY KEY AUTOINCREMENT,
--     name TEXT NOT NULL UNIQUE
-- );

-- INSERT OR IGNORE INTO CATEGORY (name) VALUES 
-- ('General'),
-- ('Lifestyle'),
-- ('Health & Fitness'),
-- ('Travel'),
-- ('Food & Cooking'),
-- ('Education'),
-- ('Business'),
-- ('Finance'),
-- ('Entertainment'),
-- ('Sports'),
-- ('Personal Dev'),
-- ('Culture'),
-- ('News');

-- ---category post 
-- CREATE TABLE IF NOT EXISTS POST_CATEGORY (
--     post_id INTEGER NOT NULL,
--     category_id INTEGER NOT NULL,

--     PRIMARY KEY (post_id, category_id),

--     FOREIGN KEY (post_id) REFERENCES POSTS(id) ON DELETE CASCADE,
--     FOREIGN KEY (category_id) REFERENCES CATEGORY(id) ON DELETE CASCADE
-- );

-- -- COMMENTS
-- CREATE TABLE IF NOT EXISTS COMMENTS (
--     id         INTEGER  NOT NULL UNIQUE,
--     user_id    INTEGER  NOT NULL,
--     post_id    INTEGER  NOT NULL,
--     created_at DATETIME NOT NULL,
--     text       TEXT     NULL    ,
--     PRIMARY KEY (id AUTOINCREMENT),
--     FOREIGN KEY (user_id) REFERENCES USERS (id) ON DELETE CASCADE,
--     FOREIGN KEY (post_id) REFERENCES POSTS (id) ON DELETE CASCADE 
-- );

-- -- POST REACTIONS
-- -- Fix POST_REACTIONS
-- CREATE TABLE IF NOT EXISTS POST_REACTIONS (
--   user_id INTEGER NOT NULL,
--   post_id INTEGER NOT NULL,
--   is_like INTEGER NOT NULL CHECK (is_like IN (-1, 1)),
--   PRIMARY KEY (user_id, post_id), -- <--- Add this
--   FOREIGN KEY (user_id) REFERENCES USERS(id) ON DELETE CASCADE,
--   FOREIGN KEY (post_id) REFERENCES POSTS(id) ON DELETE CASCADE 
-- );
-- -- reactions are unique by combination of both user_id and post_id

-- -- COMMENT REACTIONS
-- CREATE TABLE IF NOT EXISTS COMMENT_REACTIONS (
--   user_id INTEGER NOT NULL,
--   comment_id INTEGER NOT NULL,
--   is_like INTEGER NOT NULL DEFAULT 1 CHECK (is_like IN (-1, 1)), -- 1 for like / -1 for dislike
--   FOREIGN KEY (user_id) REFERENCES USERS (id) ON DELETE CASCADE ,
--   FOREIGN KEY (comment_id) REFERENCES COMMENTS (id) ON DELETE CASCADE 
-- );

-- --Rate Limits
-- CREATE TABLE IF NOT EXISTS rate_limits (
--     ip TEXT NOT NULL,
--     route TEXT NOT NULL,
--     last_request DATETIME NOT NULL,
--     PRIMARY KEY (ip, route)
-- );
-- --this for private messages
-- CREATE TABLE IF NOT EXISTS MESSAGES (
--     id INTEGER PRIMARY KEY AUTOINCREMENT,

--     conversation_id INTEGER NOT NULL,

--     sender_id INTEGER NOT NULL,

--     text TEXT NOT NULL,

-- created_at DATETIME NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%f', 'now')),

--     is_read INTEGER NOT NULL DEFAULT 0,

--     FOREIGN KEY (conversation_id) REFERENCES CONVERSATIONS(id) ON DELETE CASCADE,
    
--     FOREIGN KEY (sender_id) REFERENCES USERS(id) ON DELETE CASCADE
-- );


-- CREATE TABLE IF NOT EXISTS CONVERSATIONS (
--     id INTEGER PRIMARY KEY AUTOINCREMENT,

--     user1_id INTEGER NOT NULL,
--     user2_id INTEGER NOT NULL,

--     last_message TEXT,
--     last_message_at DATETIME,

--     user1_last_read_message_id INTEGER,
--     user2_last_read_message_id INTEGER,

--     created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

--     UNIQUE(user1_id, user2_id),

--     FOREIGN KEY (user1_id) REFERENCES USERS(id) ON DELETE CASCADE,
--     FOREIGN KEY (user2_id) REFERENCES USERS(id) ON DELETE CASCADE
-- );

-- --folowing table s


-- CREATE TABLE IF NOT EXISTS FOLLOWS (
--     follower_id INTEGER NOT NULL,
--     following_id INTEGER NOT NULL,

--     status TEXT NOT NULL DEFAULT 'pending',
--     -- pending | accepted
--     created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

--     PRIMARY KEY (follower_id, following_id),

--     FOREIGN KEY (follower_id)
--         REFERENCES USERS(id)
--         ON DELETE CASCADE,

--     FOREIGN KEY (following_id)
--         REFERENCES USERS(id)
--         ON DELETE CASCADE,

--     CHECK (follower_id != following_id)
-- );
-- --group table 
-- CREATE TABLE IF NOT EXISTS GROUPS (
--     id INTEGER PRIMARY KEY AUTOINCREMENT,

--     creator_id INTEGER NOT NULL,

--     title TEXT NOT NULL,
--     description TEXT,

--     created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

--     FOREIGN KEY (creator_id)
--         REFERENCES USERS(id)
--         ON DELETE CASCADE
-- );
-- CREATE TABLE IF NOT EXISTS GROUP_MEMBERS (
--     group_id INTEGER NOT NULL,
--     user_id INTEGER NOT NULL,

--     role TEXT DEFAULT 'member',
--     -- member | admin

--     joined_at DATETIME DEFAULT CURRENT_TIMESTAMP,


--     PRIMARY KEY(group_id,user_id),


--     FOREIGN KEY(group_id)
--         REFERENCES GROUPS(id)
--         ON DELETE CASCADE,


--     FOREIGN KEY(user_id)
--         REFERENCES USERS(id)
--         ON DELETE CASCADE
-- );
-- CREATE TABLE IF NOT EXISTS GROUP_INVITES (

--     id INTEGER PRIMARY KEY AUTOINCREMENT,

--     group_id INTEGER NOT NULL,

--     inviter_id INTEGER NOT NULL,

--     invited_user_id INTEGER NOT NULL,


--     status TEXT DEFAULT 'pending',
--     -- pending accepted rejected


--     created_at DATETIME DEFAULT CURRENT_TIMESTAMP,


--     FOREIGN KEY(group_id)
--         REFERENCES GROUPS(id)
--         ON DELETE CASCADE,


--     FOREIGN KEY(inviter_id)
--         REFERENCES USERS(id)
--         ON DELETE CASCADE,


--     FOREIGN KEY(invited_user_id)
--         REFERENCES USERS(id)
--         ON DELETE CASCADE
-- );

-- --group request
-- CREATE TABLE IF NOT EXISTS GROUP_REQUESTS (

--     id INTEGER PRIMARY KEY AUTOINCREMENT,

--     group_id INTEGER NOT NULL,

--     user_id INTEGER NOT NULL,


--     status TEXT DEFAULT 'pending',


--     created_at DATETIME DEFAULT CURRENT_TIMESTAMP,


--     UNIQUE(group_id,user_id),


--     FOREIGN KEY(group_id)
--         REFERENCES GROUPS(id)
--         ON DELETE CASCADE,


--     FOREIGN KEY(user_id)
--         REFERENCES USERS(id)
--         ON DELETE CASCADE
-- );

-- CREATE TABLE IF NOT EXISTS GROUP_MESSAGES (

--     id INTEGER PRIMARY KEY AUTOINCREMENT,

--     group_id INTEGER NOT NULL,

--     sender_id INTEGER NOT NULL,

--     text TEXT NOT NULL,

--     created_at DATETIME
--         DEFAULT CURRENT_TIMESTAMP,

--     FOREIGN KEY(group_id)
--         REFERENCES GROUPS(id)
--         ON DELETE CASCADE,

--     FOREIGN KEY(sender_id)
--         REFERENCES USERS(id)
--         ON DELETE CASCADE
-- );
-- CREATE TABLE IF NOT EXISTS GROUP_MESSAGE_READS (
--     group_id INTEGER NOT NULL,
--     user_id INTEGER NOT NULL,
--     last_read_message_id INTEGER,
--     PRIMARY KEY (group_id, user_id),
--     FOREIGN KEY (group_id) REFERENCES GROUPS(id) ON DELETE CASCADE,
--     FOREIGN KEY (user_id) REFERENCES USERS(id) ON DELETE CASCADE,
--     FOREIGN KEY (last_read_message_id) REFERENCES GROUP_MESSAGES(id) ON DELETE SET NULL
-- );
-- CREATE TABLE IF NOT EXISTS NOTIFICATIONS (
--     id INTEGER PRIMARY KEY AUTOINCREMENT,

--     user_id INTEGER NOT NULL,
--     actor_id INTEGER NOT NULL,

--     type TEXT NOT NULL,

--     object_type TEXT NOT NULL,
--     object_id INTEGER NOT NULL,

--     is_read INTEGER NOT NULL DEFAULT 0,

--     created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

--     FOREIGN KEY(user_id)
--         REFERENCES USERS(id)
--         ON DELETE CASCADE,

--     FOREIGN KEY(actor_id)
--         REFERENCES USERS(id)
--         ON DELETE CASCADE
-- );
  
-- --events in groups
-- -- GROUP EVENTS
-- CREATE TABLE IF NOT EXISTS GROUP_EVENTS (
--     id INTEGER PRIMARY KEY AUTOINCREMENT,
--     group_id INTEGER NOT NULL,
--     creator_id INTEGER NOT NULL,
--     title TEXT NOT NULL,
--     description TEXT,
--     event_time DATETIME NOT NULL,
--     created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
--     FOREIGN KEY(group_id) REFERENCES GROUPS(id) ON DELETE CASCADE,
--     FOREIGN KEY(creator_id) REFERENCES USERS(id) ON DELETE CASCADE
-- );

-- -- EVENT RESPONSES (Going / Not Going)
-- CREATE TABLE IF NOT EXISTS EVENT_RESPONSES (
--     event_id INTEGER NOT NULL,
--     user_id INTEGER NOT NULL,
--     status TEXT NOT NULL CHECK (status IN ('going', 'not_going')),
--     PRIMARY KEY (event_id, user_id),
--     FOREIGN KEY(event_id) REFERENCES GROUP_EVENTS(id) ON DELETE CASCADE,
--     FOREIGN KEY(user_id) REFERENCES USERS(id) ON DELETE CASCADE
-- );
-- CREATE TABLE IF NOT EXISTS GROUP_POSTS (
--     id INTEGER PRIMARY KEY AUTOINCREMENT,
--     group_id INTEGER NOT NULL,
--     user_id INTEGER NOT NULL,
--     created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
--     title TEXT,
--     text TEXT,
--     image TEXT,
--     FOREIGN KEY (group_id) REFERENCES GROUPS(id) ON DELETE CASCADE,
--     FOREIGN KEY (user_id) REFERENCES USERS(id) ON DELETE CASCADE
-- );

-- CREATE TABLE IF NOT EXISTS GROUP_POST_COMMENTS (
--     id INTEGER PRIMARY KEY AUTOINCREMENT,
--     group_post_id INTEGER NOT NULL,
--     user_id INTEGER NOT NULL,
--     created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
--     text TEXT,
--     FOREIGN KEY (group_post_id) REFERENCES GROUP_POSTS(id) ON DELETE CASCADE,
--     FOREIGN KEY (user_id) REFERENCES USERS(id) ON DELETE CASCADE
-- );