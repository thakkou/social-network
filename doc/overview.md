# Social Network — Project Overview

A Facebook-like social network built with a **Go backend** and **Next.js frontend**, featuring real-time chat, user profiles, posts with privacy controls, groups with events, and notifications.

## Tech Stack

| Layer | Technology |
|-------|-----------|
| **Frontend** | Next.js 14 (React, TypeScript) |
| **Backend** | Go 1.24 (net/http, gorilla/websocket) |
| **Database** | SQLite (via mattn/go-sqlite3) |
| **Authentication** | NextAuth with credentials + session cookies |
| **Real-time** | WebSocket (gorilla/websocket) |
| **Containerization** | Docker & Docker Compose |
| **Migrations** | golang-migrate |

## Key Features

- **Authentication** — Email/password registration and login with session-based auth
- **Profiles** — Public and private profiles with avatar, nickname, bio
- **Posts** — Create posts with images/GIFs, privacy levels (public, almost private, private), categories
- **Followers** — Follow/unfollow with requests for private profiles
- **Groups** — Create groups, invite members, request to join, group posts, events
- **Chat** — Real-time direct messaging and group chat via WebSockets
- **Notifications** — Real-time notifications for follows, group invites, events
- **Search** — Search users and groups

## Architecture Overview

```
┌─────────────┐     HTTP/REST     ┌──────────────┐     SQL      ┌────────┐
│             │ ─────────────────> │              │ ──────────> │        │
│  Next.js    │                    │  Go Backend  │             │ SQLite │
│  Frontend   │ <────────────────  │  (Port 8080) │ <────────── │        │
│  (Port 3000)│     JSON Response  │              │             └────────┘
│             │                    └──────────────┘
│             │                    ┌──────────────┐
│             │ ────────────────>  │  WebSocket   │
│             │ <──────────────── │  (ws://)      │
└─────────────┘    Real-time      └──────────────┘
```

## Project Structure

```
social-network/
├── backend/                    # Go backend
│   ├── server.go               # Entry point
│   ├── pkg/
│   │   ├── handlers/           # HTTP handlers
│   │   ├── repository/         # Database queries
│   │   ├── models/             # Data models & representations
│   │   ├── db/                 # Database connection & migrations
│   │   │   ├── migrations/     # SQL migration files
│   │   │   └── sqlite/         # SQLite connection, seeder
│   │   ├── routes/             # Route definitions
│   │   ├── middlewares/        # Auth, rate limiting
│   │   ├── ws/                 # WebSocket hub & client management
│   │   └── utilities/          # Image handling, validators, OAuth
│   ├── go.mod / go.sum
│   └── docs/                   # Swagger docs
├── frontend/                   # Next.js frontend
│   ├── src/
│   │   ├── app/
│   │   │   ├── (app)/          # Authenticated routes (feed, profile, groups, messages, etc.)
│   │   │   ├── (public)/       # Public routes (login, register)
│   │   │   ├── _components/    # Shared components (Sidebar, Header, etc.)
│   │   │   ├── _services/      # Server actions (CRUD operations)
│   │   │   └── _providers/     # Context providers (WS, Chat, Auth)
│   │   ├── middleware.ts       # NextAuth middleware
│   │   └── styles/globals.css  # Global styles
│   ├── package.json
│   └── next.config.js
├── Dockerfile.frontend
├── Dockerfile.backend
├── docker-compose.yml
└── script.sh                   # Build & deployment script
```

## Running the Project

### Local Development

```bash
# Start backend + frontend
./localStart.sh -i    # First time (installs deps)
./localStart.sh       # Subsequent runs
./localStart.sh -r    # Refresh (re-seeds database)
```

### Docker

```bash
# Build and start containers
./script.sh
# or manually:
docker compose up --build
```

- Frontend: http://localhost:3000
- Backend:  http://localhost:8080
