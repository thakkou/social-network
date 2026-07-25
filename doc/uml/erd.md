# Entity Relationship Diagram — Database Schema

```mermaid
%%{init: {'theme':'dark'}}%%
erDiagram
    USERS {
        int id PK
        datetime created_at
        string firstname
        string lastname
        string email UK
        string password
        string birthdate
        string nickname UK
        string aboutme
        string avatar
        bool is_private
        datetime last_seen
    }

    SESSIONS {
        string id PK
        datetime expires_at
        int user_id FK
    }

    POSTS {
        int id PK
        int user_id FK
        datetime created_at
        string title
        string text
        string image
        enum privacy "public | almost_private | private"
    }

    POST_ALLOWED_USERS {
        int post_id PK,FK
        int user_id PK,FK
    }

    POST_REACTIONS {
        int user_id PK,FK
        int post_id PK,FK
        int is_like "1 = like, -1 = dislike"
    }

    CATEGORY {
        int id PK
        string name UK
    }

    POST_CATEGORY {
        int post_id PK,FK
        int category_id PK,FK
    }

    COMMENTS {
        int id PK
        int user_id FK
        int post_id FK
        datetime created_at
        string text
        string image
    }

    COMMENT_REACTIONS {
        int user_id PK,FK
        int comment_id PK,FK
        int is_like "1 = like, -1 = dislike"
    }

    FOLLOWS {
        int follower_id PK,FK
        int following_id PK,FK
        enum status "pending | accepted"
        datetime created_at
    }

    CONVERSATIONS {
        int id PK
        int user1_id FK
        int user2_id FK
        string last_message
        datetime last_message_at
        int user1_last_read_message_id
        int user2_last_read_message_id
        datetime created_at
    }

    MESSAGES {
        int id PK
        int conversation_id FK
        int sender_id FK
        string text
        datetime created_at
        bool is_read
    }

    GROUPS {
        int id PK
        string logo
        string background
        int creator_id FK
        string title
        string description
        datetime created_at
    }

    GROUP_MEMBERS {
        int group_id PK,FK
        int user_id PK,FK
        enum role "member | admin"
        datetime joined_at
    }

    GROUP_INVITES {
        int id PK
        int group_id FK
        int inviter_id FK
        int invited_user_id FK
        enum status "pending | accepted | rejected"
        datetime created_at
    }

    GROUP_REQUESTS {
        int id PK
        int group_id FK
        int user_id FK
        enum status "pending"
        datetime created_at
    }

    GROUP_MESSAGES {
        int id PK
        int group_id FK
        int sender_id FK
        string text
        datetime created_at
    }

    GROUP_MESSAGE_READS {
        int group_id PK,FK
        int user_id PK,FK
        int last_read_message_id FK
    }

    GROUP_POSTS {
        int id PK
        int group_id FK
        int user_id FK
        datetime created_at
        string title
        string text
        string image
    }

    GROUP_POST_REACTIONS {
        int group_post_id PK,FK
        int user_id PK,FK
        int is_like "1 = like, -1 = dislike"
    }

    GROUP_POST_COMMENTS {
        int id PK
        int group_post_id FK
        int user_id FK
        datetime created_at
        string text
    }

    GROUP_POST_COMMENT_REACTIONS {
        int group_post_comment_id PK,FK
        int user_id PK,FK
        int is_like "1 = like, -1 = dislike"
    }

    GROUP_EVENTS {
        int id PK
        int group_id FK
        int creator_id FK
        string title
        string description
        datetime event_time
        datetime created_at
    }

    EVENT_RESPONSES {
        int event_id PK,FK
        int user_id PK,FK
        enum status "going | not_going"
    }

    NOTIFICATIONS {
        int id PK
        int user_id FK
        int actor_id FK
        string type "post_reaction | comment | follow_request | follow_accepted | group_invite | group_join_request | group_event"
        string object_type
        int object_id
        bool is_read
        datetime created_at
    }

    WS_TICKETS {
        string ticket PK
        int user_id FK
        datetime created_at
    }

    %% ─── Relationships ───

    USERS ||--o{ SESSIONS : "has"
    USERS ||--o{ POSTS : "creates"
    USERS ||--o{ COMMENTS : "writes"
    USERS ||--o{ FOLLOWS : "follower"
    USERS ||--o{ FOLLOWS : "following"
    USERS ||--o{ NOTIFICATIONS : "receives"
    USERS ||--o{ NOTIFICATIONS : "actor"
    USERS ||--o{ CONVERSATIONS : "user1"
    USERS ||--o{ CONVERSATIONS : "user2"
    USERS ||--o{ MESSAGES : "sender"
    USERS ||--o{ GROUP_MEMBERS : "belongs"
    USERS ||--o{ GROUP_INVITES : "inviter"
    USERS ||--o{ GROUP_INVITES : "invited"
    USERS ||--o{ GROUP_REQUESTS : "requester"
    USERS ||--o{ GROUP_MESSAGES : "sender"
    USERS ||--o{ GROUP_POSTS : "author"
    USERS ||--o{ GROUP_POST_COMMENTS : "author"
    USERS ||--o{ WS_TICKETS : "gets"

    POSTS ||--o{ COMMENTS : "has"
    POSTS ||--o{ POST_REACTIONS : "has"
    POSTS ||--o{ POST_ALLOWED_USERS : "visible by"
    POSTS ||--o{ POST_CATEGORY : "categorized"

    COMMENTS ||--o{ COMMENT_REACTIONS : "has"

    CATEGORY ||--o{ POST_CATEGORY : "belongs to"

    CONVERSATIONS ||--o{ MESSAGES : "contains"

    GROUPS ||--o{ GROUP_MEMBERS : "has"
    GROUPS ||--o{ GROUP_INVITES : "sent"
    GROUPS ||--o{ GROUP_REQUESTS : "received"
    GROUPS ||--o{ GROUP_MESSAGES : "chat"
    GROUPS ||--o{ GROUP_POSTS : "contains"
    GROUPS ||--o{ GROUP_EVENTS : "hosts"
    GROUPS ||--o{ GROUP_MESSAGE_READS : "tracks reads"

    GROUP_POSTS ||--o{ GROUP_POST_REACTIONS : "has"
    GROUP_POSTS ||--o{ GROUP_POST_COMMENTS : "has"

    GROUP_POST_COMMENTS ||--o{ GROUP_POST_COMMENT_REACTIONS : "has"

    GROUP_EVENTS ||--o{ EVENT_RESPONSES : "responses"
```
