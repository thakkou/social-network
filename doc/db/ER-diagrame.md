# Database Schema — Entity Relationship Diagram

## Overview

The database uses SQLite with 25+ tables covering users, posts, groups, messaging, and notifications. Migrations are applied using [golang-migrate](https://github.com/golang-migrate/migrate).

## Tables

### Core Tables

#### `USERS`
| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | Unique user ID |
| created_at | DATETIME | NOT NULL | Account creation timestamp |
| firstname | TEXT | NOT NULL | User's first name |
| lastname | TEXT | NOT NULL | User's last name |
| email | TEXT | NOT NULL UNIQUE | Login email |
| password | TEXT | NOT NULL | Bcrypt-hashed password |
| birthdate | TEXT | NOT NULL | Date of birth |
| nickname | TEXT | UNIQUE | Display handle (optional) |
| aboutme | TEXT | | Short biography (optional) |
| avatar | TEXT | | Avatar image path (optional) |
| is_private | INTEGER | DEFAULT 0 | Profile privacy flag |

#### `SESSIONS`
| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | TEXT | PRIMARY KEY | UUID session token |
| expires_at | DATETIME | NOT NULL | Session expiration |
| user_id | INTEGER | NOT NULL, FK → USERS(id) | Associated user |

### Posts & Social

#### `POSTS`
| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | Post ID |
| user_id | INTEGER | NOT NULL, FK → USERS(id) | Author |
| created_at | DATETIME | NOT NULL | Creation timestamp |
| title | TEXT | NOT NULL | Post title |
| text | TEXT | | Post content |
| image | TEXT | | Image path |
| privacy | TEXT | NOT NULL | `public`, `almost_private`, `private` |

#### `POST_ALLOWED_USERS`
| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| post_id | INTEGER | PK, FK → POSTS(id) | Post reference |
| user_id | INTEGER | PK, FK → USERS(id) | Allowed user |

#### `POST_REACTIONS`
| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| user_id | INTEGER | PK, FK → USERS(id) | Reacting user |
| post_id | INTEGER | PK, FK → POSTS(id) | Target post |
| is_like | INTEGER | NOT NULL | 1 = like, -1 = dislike |

#### `CATEGORY`
| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | Category ID |
| name | TEXT | NOT NULL UNIQUE | Category name |

#### `POST_CATEGORY`
| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| post_id | INTEGER | PK, FK → POSTS(id) | Post reference |
| category_id | INTEGER | PK, FK → CATEGORY(id) | Category reference |

#### `COMMENTS`
| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | Comment ID |
| user_id | INTEGER | NOT NULL, FK → USERS(id) | Author |
| post_id | INTEGER | NOT NULL, FK → POSTS(id) | Parent post |
| created_at | DATETIME | NOT NULL | Creation timestamp |
| text | TEXT | NOT NULL | Comment content |
| image | TEXT | | Image path |

#### `COMMENT_REACTIONS`
| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| user_id | INTEGER | PK, FK → USERS(id) | Reacting user |
| comment_id | INTEGER | PK, FK → COMMENTS(id) | Target comment |
| is_like | INTEGER | NOT NULL | 1 = like, -1 = dislike |

### Follow System

#### `FOLLOWS`
| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| follower_id | INTEGER | PK, FK → USERS(id) | User who follows |
| following_id | INTEGER | PK, FK → USERS(id) | User being followed |
| status | TEXT | NOT NULL | `pending`, `accepted` |
| created_at | DATETIME | NOT NULL | Follow timestamp |

### Groups

#### `GROUPS`
| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | Group ID |
| creator_id | INTEGER | NOT NULL, FK → USERS(id) | Group creator |
| title | TEXT | NOT NULL | Group name |
| description | TEXT | | Group description |
| logo | TEXT | | Group logo image path |
| background | TEXT | | Group banner image path |
| created_at | DATETIME | NOT NULL | Creation timestamp |

#### `GROUP_MEMBERS`
| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| group_id | INTEGER | PK, FK → GROUPS(id) | Group reference |
| user_id | INTEGER | PK, FK → USERS(id) | Member |
| role | TEXT | NOT NULL | `member`, `admin` |
| joined_at | DATETIME | NOT NULL | Join timestamp |

#### `GROUP_INVITES`
| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | Invite ID |
| group_id | INTEGER | NOT NULL, FK → GROUPS(id) | Target group |
| inviter_id | INTEGER | NOT NULL, FK → USERS(id) | Who invited |
| invited_user_id | INTEGER | NOT NULL, FK → USERS(id) | Invited user |
| status | TEXT | NOT NULL | `pending`, `accepted`, `rejected` |
| created_at | DATETIME | NOT NULL | Invite timestamp |

#### `GROUP_REQUESTS`
| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | Request ID |
| group_id | INTEGER | NOT NULL, FK → GROUPS(id) | Target group |
| user_id | INTEGER | NOT NULL, FK → USERS(id) | Requester |
| status | TEXT | NOT NULL DEFAULT 'pending' | `pending` |
| created_at | DATETIME | NOT NULL | Request timestamp |

#### `GROUP_POSTS`
| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | Post ID |
| group_id | INTEGER | NOT NULL, FK → GROUPS(id) | Parent group |
| user_id | INTEGER | NOT NULL, FK → USERS(id) | Author |
| created_at | DATETIME | NOT NULL | Creation timestamp |
| title | TEXT | NOT NULL | Post title |
| text | TEXT | | Post content |
| image | TEXT | | Image path |

#### `GROUP_POST_REACTIONS`
| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| group_post_id | INTEGER | PK, FK → GROUP_POSTS(id) | Target post |
| user_id | INTEGER | PK, FK → USERS(id) | Reacting user |
| is_like | INTEGER | NOT NULL | 1 = like, -1 = dislike |

#### `GROUP_POST_COMMENTS`
| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | Comment ID |
| group_post_id | INTEGER | NOT NULL, FK → GROUP_POSTS(id) | Parent post |
| user_id | INTEGER | NOT NULL, FK → USERS(id) | Author |
| created_at | DATETIME | NOT NULL | Creation timestamp |
| text | TEXT | NOT NULL | Comment content |

#### `GROUP_POST_COMMENT_REACTIONS`
| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| group_post_comment_id | INTEGER | PK, FK → GROUP_POST_COMMENTS(id) | Target comment |
| user_id | INTEGER | PK, FK → USERS(id) | Reacting user |
| is_like | INTEGER | NOT NULL | 1 = like, -1 = dislike |

#### `GROUP_EVENTS`
| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | Event ID |
| group_id | INTEGER | NOT NULL, FK → GROUPS(id) | Hosting group |
| creator_id | INTEGER | NOT NULL, FK → USERS(id) | Event creator |
| title | TEXT | NOT NULL | Event title |
| description | TEXT | | Event description |
| event_time | DATETIME | NOT NULL | Scheduled date/time |
| created_at | DATETIME | NOT NULL | Creation timestamp |

#### `EVENT_RESPONSES`
| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| event_id | INTEGER | PK, FK → GROUP_EVENTS(id) | Target event |
| user_id | INTEGER | PK, FK → USERS(id) | Responding user |
| status | TEXT | NOT NULL | `going`, `not_going` |

### Messaging

#### `CONVERSATIONS`
| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | Conversation ID |
| user1_id | INTEGER | NOT NULL, FK → USERS(id) | Participant 1 |
| user2_id | INTEGER | NOT NULL, FK → USERS(id) | Participant 2 |
| last_message | TEXT | | Preview of last message |
| last_message_at | DATETIME | | Timestamp of last message |
| user1_last_read_message_id | INTEGER | | Last read message by user1 |
| user2_last_read_message_id | INTEGER | | Last read message by user2 |
| created_at | DATETIME | NOT NULL | Creation timestamp |

#### `MESSAGES`
| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | Message ID |
| conversation_id | INTEGER | NOT NULL, FK → CONVERSATIONS(id) | Parent conversation |
| sender_id | INTEGER | NOT NULL, FK → USERS(id) | Sender |
| text | TEXT | | Message content |
| created_at | DATETIME | NOT NULL | Send timestamp |
| is_read | INTEGER | DEFAULT 0 | Read status |

#### `GROUP_MESSAGES`
| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | Message ID |
| group_id | INTEGER | NOT NULL, FK → GROUPS(id) | Target group |
| sender_id | INTEGER | NOT NULL, FK → USERS(id) | Sender |
| text | TEXT | NOT NULL | Message content |
| created_at | DATETIME | NOT NULL | Send timestamp |

#### `GROUP_MESSAGE_READS`
| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| group_id | INTEGER | PK, FK → GROUPS(id) | Group reference |
| user_id | INTEGER | PK, FK → USERS(id) | Reader |
| last_read_message_id | INTEGER | FK → GROUP_MESSAGES(id) | Last read marker |

### Notifications & WebSocket

#### `NOTIFICATIONS`
| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | INTEGER | PRIMARY KEY AUTOINCREMENT | Notification ID |
| user_id | INTEGER | NOT NULL, FK → USERS(id) | Recipient |
| actor_id | INTEGER | NOT NULL, FK → USERS(id) | Who triggered it |
| type | TEXT | NOT NULL | `follow_request`, `follow_accepted`, `new_follower`, `group_invite`, `group_join_request`, `group_event`, `post_reaction`, `comment` |
| object_type | TEXT | | Related object type |
| object_id | INTEGER | | Related object ID |
| is_read | INTEGER | DEFAULT 0 | Read status |
| created_at | DATETIME | NOT NULL | Creation timestamp |

#### `WS_TICKETS`
| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| ticket | TEXT | PRIMARY KEY | UUID ticket |
| user_id | INTEGER | NOT NULL, FK → USERS(id) | Authenticated user |
| created_at | DATETIME | NOT NULL | Creation timestamp |

## Indexes

The database uses SQLite's automatic indexes on PRIMARY KEY and UNIQUE columns. Foreign keys are enforced via `PRAGMA foreign_keys = ON`.

## Migrations

Migrations are located in `backend/pkg/db/migrations/sqlite/`:

| # | Name | Files |
|---|------|-------|
| 1 | Create users & sessions | `000001_create_users_sessions.up/down.sql` |
| 2 | Create follows | `000002_create_follows.up/down.sql` |
| 3 | Create posts | `000003_create_posts.up/down.sql` |
| 4 | Create groups | `000004_create_groups.up/down.sql` |
| 5 | Create chat | `000005_create_chat.up/down.sql` |
| 6 | Create notifications & events | `000006_create_notifications_events.up/down.sql` |
| 7 | Add event is_finished | `000007_add_event_is_finished.up/down.sql` |
| 8 | Add event image | `000008_add_event_image.up/down.sql` |

## Seeder

The seeder (`backend/pkg/db/sqlite/seeder/`) creates test data:

- **6 users** — Alice, Bob, Chloe, Farid, Isabella, Jack
- **Posts & comments** — Various sample posts with privacy levels and comments
- **Follows** — Mixed follow relationships
- **Groups** — Multiple groups with members and invites
- **Events** — Sample events with responses
- **Messages** — Conversation history
- **Notifications** — Pre-seeded notification history
