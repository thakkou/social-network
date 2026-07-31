# Backend Test Suite

Integration and functional tests for the social-network backend API.

## Quick Start

```bash
# Run ALL tests via the main test (from the backend/ directory)
RUN_ALL_SUITES=1 go test ./tests/ -run TestRunAllSuites -count=1 -timeout 1800s -v

# Run all tests directly (equivalent, runs every suite in parallel)
go test ./tests/... -count=1 -timeout 180s

# Run a specific test suite
go test ./tests/auth/... -count=1 -timeout 60s

# Run a single test
go test ./tests/auth/... -run TestLogin_Success -count=1 -timeout 30s -v

# Run tests with verbose output
go test ./tests/followers/... -count=1 -timeout 60s -v
```

> **Note:** The `-count=1` flag disables result caching, ensuring tests always run fresh.

> **Main test** (`tests/main_test.go`): `TestRunAllSuites` is a meta-test that
> runs every suite (auth, users, followers, posts, groups, chat, notifications,
> integration) one after another and prints a per-suite summary. A failure in
> one suite does not stop the others from running. It requires `RUN_ALL_SUITES=1`
> (otherwise it is skipped), so a plain `go test ./tests/...` does not spawn
> every suite a second time.

---

## Test Architecture

```
tests/
├── README.md              ← This file
├── setup/                 ← Test scaffolding (DB, server, helpers)
│   ├── testdb.go          ← Creates temp-file SQLite databases
│   ├── migrate.go         ← Runs database migrations
│   ├── seed.go            ← Loads seeder_test.sql
│   ├── seeder_test.sql    ← 10 users, 63 posts, 88 comments, groups, etc.
│   ├── router.go          ← Builds a full test HTTP server
│   ├── request.go         ← Fluent HTTP request builder (JSON, POST, etc.)
│   ├── auth.go            ← Session/cookie creation helpers
│   ├── assertions.go      ← Custom assertion functions
│   └── cleanup.go         ← Test DB & server teardown
├── fixtures/              ← Typed constants for seeded data
│   ├── users.go           ← Alice, Bob, Chloe, ... + AllUsers(), CommonPassword
│   ├── posts.go           ← PostsByAuthor, PublicPostIDs(), ...
│   ├── groups.go          ← GophersUnited, SportsClub, ...
│   ├── follows.go         ← FollowRelationships, IsFollowing(), ...
│   ├── comments.go        ← CommentsByPost map
│   ├── messages.go        ← ConvAliceBob, MessagesByConversation
│   └── events.go          ← GoMeetup, EventResponses, ...
├── auth/                  ← Authentication tests
├── users/                 ← Profile & privacy tests
├── followers/             ← Follow/unfollow/accept tests
├── posts/                 ← Post CRUD, comments, reactions
├── groups/                ← Groups, members, events, invites
├── chat/                  ← Private & group messaging
├── notifications/         ← Notification CRUD & flows
└── integration/           ← Multi-step integration scenarios
```

---

## Test Patterns

### 1. Standard test structure (HTTP integration test)

Every test follows this pattern:

```go
func TestSomething(t *testing.T) {
    // 1. Create a test server with fresh DB, migrations, and seeded data
    server, db, err := setup.NewTestServerWithDefaults()
    if err != nil {
        t.Fatal(err)
    }
    defer setup.CleanupTestServer(server, db)

    // 2. Create a session for the authenticated user
    sessionID, err := setup.CreateSession(db, fixtures.Alice.ID)
    if err != nil {
        t.Fatal(err)
    }

    // 3. Make HTTP requests using the fluent builder
    body, err := setup.GET(server, "/api/posts").
        WithAuth(sessionID).
        DoStatus(200)   // asserts 200 status

    // 4. Parse and assert on the response
    var wrapper struct {
        Message string `json:"message"`
    }
    setup.MustUnmarshalJSON(t, body, &wrapper)
}
```

### 2. Request builder methods

| Method | Usage |
|--------|-------|
| `setup.GET(server, path)` | GET request |
| `setup.POST(server, path)` | POST request |
| `setup.PUT(server, path)` | PUT request |
| `setup.DELETE(server, path)` | DELETE request |
| `setup.JSON(server, method, path, payload)` | Request with JSON body |
| `.WithAuth(sessionID)` | Attach session cookie |
| `.WithHeader(key, value)` | Add custom header |
| `.WithMultipart(fields, files)` | Form data / file upload |
| `.DoStatus(expectedCode)` | Execute + assert status |
| `.DoOK()` | Execute + assert 2xx |
| `.Do()` | Execute + return raw response |

### 3. Using fixtures (typed constants)

Instead of hardcoding IDs, use the fixtures package:

```go
// ❌ Hardcoded — fragile, unclear
session, _ := setup.CreateSession(db, 1)
setup.PUT(server, "/api/follow/follow/3").WithAuth(session)

// ✅ Typed — readable, maintainable
session, _ := setup.CreateSession(db, fixtures.Alice.ID)
setup.PUT(server, "/api/follow/follow/"+strconv.Itoa(fixtures.Chloe.ID)).
    WithAuth(session)
```

Available fixtures:

| Constant | ID | Notes |
|----------|----|-------|
| `fixtures.Alice` | 1 | Public, nickname "ali_m" |
| `fixtures.Bob` | 2 | Public, no nickname |
| `fixtures.Chloe` | 3 | **Private** profile |
| `fixtures.David` | 4 | Public, nickname "david" |
| `fixtures.Emma` | 5 | **Private** profile |
| `fixtures.Farid` | 6 | Public |
| `fixtures.Grace` | 7 | Public |
| `fixtures.Hugo` | 8 | Public, nickname "costa77" |
| `fixtures.Bella` | 9 | Public, nickname "bella" |
| `fixtures.Jack` | 10 | Public |
| `fixtures.CommonPassword` | — | `"password123"` (all seeded users) |

---

## Test Suites

### `tests/auth/` — Authentication

| Test File | What It Tests |
|-----------|---------------|
| `register_test.go` | User registration: success, missing fields, duplicate email, invalid email, short password |
| `login_test.go` | Login with email/nickname, wrong password, non-existent user, wrong HTTP method |
| `logout_test.go` | Logout with valid/invalid session, unauthenticated |
| `session_test.go` | Session validation (`/api/session/validate`), `/api/me` endpoint |
| `liddlware_test.go` | Auth middleware: blocks unauthenticated requests, allows authenticated requests |

### `tests/users/` — Profile & Privacy

| Test File | What It Tests |
|-----------|---------------|
| `profile_test.go` | View own profile, another user's profile, private user profile, non-existent user |
| `privacy_test.go` | Toggle profile private/public, invalid values, unauthenticated |
| `update_profile_test.go` | Update nickname + about me, unauthenticated |
| `avatar_test.go` | Update avatar with multipart upload |

### `tests/followers/` — Follow System

| Test File | What It Tests |
|-----------|---------------|
| `follow_test.go` | Follow public user, private user (pending), self-follow (rejected), already following (idempotent) |
| `unfollow_test.go` | Unfollow an followed user, unfollow when not following |
| `accept_decline_test.go` | Accept pending follow, reject pending follow, accept non-existent request |
| `follow_request.go` | Full follow flow: request → accept → unfollow |

### `tests/posts/` — Posts & Comments

| Test File | What It Tests |
|-----------|---------------|
| `create_post_test.go` | Create post with categories, missing title, missing categories, invalid privacy, private post with allowed users |
| `posts_test.go` | Get feed, category filter, liked-by-me filter, posted-by-me filter, pagination, single post, like, dislike, delete own/others' post, categories list |
| `comments_test.go` | Create comment, empty text, invalid post, like comment, delete own comment |
| `events_test.go` | Event response (going), events list |
| `event_response_test.go` | Event response (not_going) |
| `invitations_test.go` | Placeholder for post invitation flows |

### `tests/groups/` — Groups

| Test File | What It Tests |
|-----------|---------------|
| `create_group_test.go` | Create group, missing title, unauthenticated |
| `members_test.go` | List groups, get my groups, join and leave, kick member |
| `join_requests_test.go` | View pending requests (admin), accept request, reject request, non-admin access (blocked) |
| `invitations_test.go` | Accept group invite, reject group invite |
| `events_test.go` | Event response (going), invalid response status |
| `event_responses_test.go` | Non-member event response |
| `posts_test.go` | Get group content, react to group post |
| `comments_test.go` | Create group post comment, react to group post comment |

### `tests/chat/` — Messaging

| Test File | What It Tests |
|-----------|---------------|
| `private_chat_test.go` | Send direct message, empty text, self-message (blocked), no follow relationship (blocked), get conversations |
| `group_chat_test.go` | Send group message, non-member (blocked), get group messages |
| `websocket_test.go` | Create WS ticket, get online users |
| `permissions_test.go` | Conversation access control, unauthenticated message (blocked) |

### `tests/notifications/` — Notifications

| Test File | What It Tests |
|-----------|---------------|
| `unread_test.go` | Get all notifications, get unread only, mark single read, mark all read, delete single, delete all, unauthenticated |
| `follow_notifications_test.go` | Follow request notification, follow accepted notification |
| `group_notifications_test.go` | Group invite notification, group join request notification |
| `event_notifications_test.go` | Event notification exists, multiple users have event notifications |

### `tests/integration/` — End-to-End Flows

| Test File | What It Tests |
|-----------|---------------|
| `social_flow_test.go` | Register → login → create post → view profile; Follow private user → accept → send message |
| `group_flow_test.go` | Create group → join request → accept → send group message; Invite user → accept invite |
| `notification_flow_test.go` | Follow → new follower notification; Comment → comment notification; Like → reaction notification |
| `websocket_flow_test.go` | WS ticket creation (authorized), ticket creation (unauthorized) |

---

## Seeded Test Data

The `setup/seeder_test.sql` file contains a comprehensive dataset:

| Entity | Count | Notes |
|--------|-------|-------|
| Users | 10 | 2 private (Chloe, Emma), 8 public |
| Follows | 14 | 12 accepted, 2 pending |
| Posts | 63 | Public, almost_private, and private |
| Comments | 88 | Across many posts |
| Post Reactions | 50+ | Likes and dislikes |
| Comment Reactions | 32+ | Likes and dislikes |
| Conversations | 2 | Alice↔Bob, Alice↔Farid |
| Messages | 4 | 2 per conversation |
| Groups | 4 | With members, admins |
| Group Invites | 4 | Pending invitations |
| Group Join Requests | 2 | Pending requests |
| Group Messages | 12 | 3 per group |
| Group Posts | 8 | 2 per group |
| Events | 4 | 1 per group |
| Event Responses | 6 | Going/not_going |
| Notifications | 22 | All types (follow, reaction, comment, group, event) |

All users share the same password hash for `"password123"`.

---

## Writing New Tests

### Adding a new test to an existing suite

```go
// tests/mysuite/my_test.go
package mysuite

import (
    "testing"
    "01social/tests/setup"
    "01social/tests/fixtures"
)

func TestMyNewFeature(t *testing.T) {
    server, db, err := setup.NewTestServerWithDefaults()
    if err != nil {
        t.Fatal(err)
    }
    defer setup.CleanupTestServer(server, db)

    // Use fixtures for readability
    sessionID, _ := setup.CreateSession(db, fixtures.Alice.ID)

    body, _ := setup.GET(server, "/api/my-endpoint").
        WithAuth(sessionID).
        DoStatus(200)

    var wrapper struct {
        Message string `json:"message"`
    }
    setup.MustUnmarshalJSON(t, body, &wrapper)
}
```

### When to create a new test suite

Add a new directory under `tests/` when testing a new API domain. Follow the existing pattern:

1. Create the directory
2. Add `*_test.go` files with `package <name>`
3. Use `setup.NewTestServerWithDefaults()` for server setup
4. Use `fixtures.*` constants for readable test data references

---

## Troubleshooting

| Symptom | Likely Cause |
|---------|--------------|
| `no such table: ...` | Database connection pool issue. Temp file DB ensures all connections see the same data. |
| `crypto/bcrypt: hashedPassword is not the hash` | Wrong password. All seeded users use `fixtures.CommonPassword` (`"password123"`). |
| `409 Conflict` on login endpoints | You're sending a valid session cookie to an endpoint with `requiresAuth=false`. The middleware returns 409 when an already-authenticated user hits a no-auth endpoint. |
| Tests are slow | Each test rebuilds the full database. This is intentional for isolation. Use `-run <TestName>` to run specific tests faster. |

### Test database

Each test gets a **temporary file-based SQLite database** (`/tmp/social-test-*.db`) that is automatically cleaned up when `CleanupTestServer` runs. The file-based approach avoids SQLite's per-connection `:memory:` isolation issue.
