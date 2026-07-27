# API Reference

Base URL: `http://localhost:8080/api`

## Authentication

All authenticated endpoints require a valid session cookie.

### Public Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/auth/register` | Register a new user |
| POST | `/api/auth/login` | Log in with email/nickname + password |
| POST | `/api/auth/session/validate` | Validate a session |
| POST | `/api/auth/logout` | Log out (destroy session) |
| POST | `/api/auth/oauth/github` | GitHub OAuth login |
| POST | `/api/auth/oauth/github/callback` | GitHub OAuth callback |

### Registration Fields

| Field | Type | Required |
|-------|------|----------|
| email | string | ✅ |
| password | string | ✅ |
| firstname | string | ✅ |
| lastname | string | ✅ |
| birthDate | string (date) | ✅ |
| avatar | file | ❌ (optional) |
| nickname | string | ❌ (optional) |
| aboutme | string | ❌ (optional) |

## Users & Profiles

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/profile/{id}` | Get user profile (public info + posts + followers/following) |
| GET | `/api/profile/my` | Get current user's profile |
| PUT | `/api/profile/update` | Update profile (avatar, nickname, about me) |
| PUT | `/api/profile/privacy` | Toggle profile privacy (public/private) |

### Profile Response

```json
{
  "status_code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "firstname": "Amir",
    "lastname": "Kader",
    "email": "amir@example.com",
    "nickname": "@amir",
    "aboutme": "Hello world!",
    "avatar": "/uploads/avatars/uuid.jpg",
    "birthdate": "2000-01-15",
    "is_private": false,
    "created_at": "2025-01-01T00:00:00Z",
    "following_status": "none",
    "followers": [...],
    "following": [...],
    "posts": [...]
  }
}
```

## Posts

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/posts/create` | Create a post (multipart: title, text, image, privacy, category_ids, allowed_users) |
| GET | `/api/posts/{id}` | Get a single post with comments |
| GET | `/api/posts/feed/{last_id}` | Get paginated feed |
| DELETE | `/api/posts/{id}/delete` | Delete own post |
| POST | `/api/posts/{id}/reaction` | Like/dislike (body: `{is_like: 1 | -1}`) |

### Privacy Levels

| Value | Description |
|-------|-------------|
| `public` | Visible to all users |
| `almost_private` | Visible to followers only |
| `private` | Visible only to specified users |

## Comments

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/posts/{id}/comments` | Create comment (body: `{text, image?}`) |
| DELETE | `/api/posts/{id}/comments/{commentId}` | Delete own comment |
| POST | `/api/posts/{id}/comments/{commentId}/reaction` | Like/dislike comment |

## Follows

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/follow/{targetId}` | Follow/unfollow a user |
| GET | `/api/follow/requests` | Get pending follow requests |
| POST | `/api/follow/{requesterId}/accept` | Accept follow request |
| POST | `/api/follow/{requesterId}/reject` | Reject follow request |

## Groups

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/groups/create` | Create group (multipart: title, description, logo, background, invite_ids) |
| GET | `/api/groups/public/{id}` | Get public group info |
| GET | `/api/groups/content/{id}` | Get group feed (posts + events) |
| POST | `/api/groups/{id}/invite` | Invite user to group |
| POST | `/api/groups/{id}/join` | Request to join group |
| POST | `/api/groups/{id}/leave` | Leave group |
| PUT | `/api/groups/{id}/update` | Update group settings |
| POST | `/api/groups/{id}/invites/accept` | Accept group invitation |
| POST | `/api/groups/{id}/invites/reject` | Reject group invitation |
| GET | `/api/groups/{id}/requests` | Get pending join requests (creator only) |
| POST | `/api/groups/{id}/requests/{userId}/accept` | Accept join request |
| POST | `/api/groups/{id}/requests/{userId}/reject` | Reject join request |
| GET | `/api/groups/members/{id}` | Get group members |
| GET | `/api/groups/invite-candidates/{id}` | Get users who can be invited |
| POST | `/api/groups/{id}/members/{userId}/kick` | Remove member from group |

### Group Posts

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/groups/{id}/posts` | Create group post |
| POST | `/api/groups/{id}/posts/{postId}/reaction` | Like/dislike group post |
| POST | `/api/groups/{id}/posts/{postId}/comments` | Comment on group post |
| POST | `/api/groups/{id}/posts/{postId}/comments/{commentId}/reaction` | Like/dislike comment |
| POST | `/api/groups/{id}/posts/{postId}/delete` | Delete group post |
| POST | `/api/groups/{id}/posts/{postId}/comments/{commentId}/delete` | Delete group comment |

### Group Events

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/groups/{id}/events` | Create event (title, description, event_time, image?) |
| POST | `/api/groups/{id}/events/{eventId}/respond` | Respond to event (body: `{status: "going" | "not_going"}`) |
| POST | `/api/groups/{id}/events/{eventId}/delete` | Delete event |

## Messaging

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/messages/{otherUserId}` | Get conversation with another user |
| GET | `/api/conversations` | Get all conversations + groups |
| POST | `/api/messages` | Send a message |
| POST | `/api/messages/typing` | Send typing indicator |

## Notifications

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/notifications` | Get all notifications |
| POST | `/api/notifications/{id}/read` | Mark notification as read |
| POST | `/api/notifications/read-all` | Mark all as read |
| DELETE | `/api/notifications/{id}` | Delete single notification |
| DELETE | `/api/notifications/clear-all` | Delete all notifications |

### Notification Types

| Type | Description |
|------|-------------|
| `follow_request` | Someone wants to follow you (private profile) |
| `follow_accepted` | Someone accepted your follow request |
| `new_follower` | A public user started following you |
| `group_invite` | You've been invited to a group |
| `group_join_request` | Someone wants to join your group |
| `group_event` | An event was created in your group |
| `post_reaction` | Someone liked/disliked your post |
| `comment` | Someone commented on your post |

## Search

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/search?q={query}` | Search users and groups |

## WebSocket

### Connection

```
ws://localhost:8080/ws?ticket=<ticket>
```

### Getting a Ticket

```
POST /api/ws-ticket
Cookie: session=<session_cookie>
→ { "ticket": "uuid-ticket" }
```

### Events (Client → Server)

| Event | Data | Description |
|-------|------|-------------|
| `typing` | `{conversation_id, user_id}` | User is typing |

### Events (Server → Client)

| Event | Data | Description |
|-------|------|-------------|
| `new_posts` | `{post_id, title, author}` | New post created |
| `new_message` | `{conversation_id, text, sender}` | New direct message |
| `new_group_message` | `{group_id, text, sender}` | New group message |
| `follow_request` | `{actor_id, nickname}` | Follow request received |
| `follow_accepted` | `{actor_id, nickname}` | Follow request accepted |
| `group_invite` | `{group_id, title, inviter}` | Group invitation |
| `group_join_request` | `{user_id, nickname}` | Join request for your group |
| `group_event` | `{group_id, title}` | New event in group |
| `notification` | `{id, type, actor, object}` | Generic notification event |

## User Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/users/groups` | Get current user's groups |
| GET | `/api/users/{id}` | Get user data by ID |
| GET | `/api/users/{id}/followers` | Get user's followers |
| GET | `/api/users/{id}/following` | Get users a user follows |
