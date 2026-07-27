# Architecture

## System Overview

The social network follows a **client-server architecture** with three main layers:

1. **Frontend** — Next.js application serving the UI and handling client-side logic
2. **Backend** — Go HTTP server handling business logic and data persistence
3. **Database** — SQLite for persistent storage

## Frontend Architecture

### Next.js App Router

The frontend uses Next.js 14 App Router with a hybrid approach:

- **Server Actions** (`'use server'` functions in `_services/crud/`) — handle data mutations by making HTTP requests to the Go backend
- **Client Components** — interactive UI components with state management
- **Providers** — React context providers for session, WebSocket, and chat state

### Key Components

| Component | Description |
|-----------|-------------|
| `Sidebar` | Left sidebar with navigation, my groups, following, followers (with avatars) |
| `Sidebar2` | Right sidebar for messages (conversations, discover, shared groups) |
| `Header` | Top navigation bar with search, notifications badge, profile avatar |
| `MobileNav` | Bottom navigation for mobile viewports |
| `SearchInput` | Search-as-you-type for users and groups |
| `EmojiPicker` | Emoji selector used in chat and posts |
| `Toast` | Toast notification system |

### Data Flow

```
User Action → Client Component → Server Action ('use server')
  → fetchApi() → HTTP Request → Go Backend → SQLite
  → JSON Response → Server Action → Client Component → UI Update
```

Real-time updates flow through WebSocket:
```
Go Backend → WebSocket Hub → Frontend WebSocket Provider → Component State → UI
```

## Backend Architecture

### HTTP Layer

The Go backend uses the standard `net/http` library with a custom router:

- **Router** (`pkg/routes/routes.go`) — defines all API routes and middleware chains
- **Middleware** — session cookie validation, rate limiting per endpoint
- **Handlers** — per-feature HTTP handlers in `pkg/handlers/`

### Repository Pattern

Database access follows the repository pattern:

```
Handlers → Repository Interfaces → SQLite Implementation
```

Each feature has its own repository file in `pkg/repository/`:
- `user.go` — user CRUD and session management
- `post.go` — posts with privacy filtering
- `follow.go` — follow relationships with status management
- `group_repo.go` — groups, members, invites, requests
- `group_events.go` — events and responses
- `group_posts.go` — group posts, comments, reactions
- `conversation.go` — direct messages and conversations
- `notification.go` — notifications CRUD
- `reaction.go` — post/comment reactions
- `category.go` — post categories
- `profile.go` — profile data aggregation

### WebSocket System

The WebSocket system (`pkg/ws/`) manages real-time communication:

- **Ticket-based Authentication** — clients obtain a WS ticket via `/api/ws-ticket`, then connect with `ws://host/ws?ticket=<ticket>`
- **Client Manager** — tracks connected clients per user ID
- **Event Broadcasting** — typed events (new_message, new_post, notification, typing, etc.)
- **User-specific Notifications** — `NotifyUser(userId, event, data)` sends to specific connected clients

### Middleware Stack

1. **Rate Limiter** — configurable requests per second per endpoint
2. **Auth Middleware** — validates session cookie, sets user context
3. **Handler** — executes the request logic

## Database Architecture

### Migration System

Migrations use golang-migrate with SQL files in `pkg/db/migrations/sqlite/`:

```
000001_create_users_sessions.up/down.sql
000002_create_follows.up/down.sql
000003_create_posts.up/down.sql
000004_create_groups.up/down.sql
000005_create_chat.up/down.sql
000006_create_notifications_events.up/down.sql
000007_add_event_is_finished.up/down.sql
000008_add_event_image.up/down.sql
```

### Seeder

The database can be seeded with test data using `pkg/db/sqlite/seeder/`:
- Creates 6 test users (Alice, Bob, Chloe, Farid, Isabella, Jack)
- Generates sample posts, comments, follows, groups, events, messages, and notifications
- Triggered with `-r` flag: `go run . -r`

## Authentication Flow

```
1. User submits login form → NextAuth credentials provider
2. NextAuth calls backend /api/auth/session/validate
3. Backend validates credentials, creates session, returns user data
4. NextAuth sets httpOnly session cookie
5. Subsequent requests validate session via middleware
```

## Container Architecture

```
┌─────────────────────────────────────┐
│         Docker Network              │
│  ┌──────────────┐  ┌────────────┐  │
│  │  Frontend    │  │  Backend   │  │
│  │  :3000       │  │  :8080     │  │
│  │  Next.js     │  │  Go        │  │
│  └──────┬───────┘  └─────┬──────┘  │
│         │                 │         │
│         └─────────────────┘         │
└─────────────────────────────────────┘
```

## Security

- **Passwords** — hashed with bcrypt
- **Sessions** — UUID-based session tokens stored in httpOnly cookies
- **Rate Limiting** — per-endpoint rate limits to prevent abuse
- **Private Profiles** — unauthenticated/non-follower users cannot view private profiles or posts
- **Post Privacy** — three-tier privacy: public, almost private (followers only), private (specific users)
- **Message Restrictions** — users can only message users they follow or who follow them
