# social-network

# floow ws-ticket

[Browser]                 [Next.js]                   [Go Backend]
    │                         │                            │
    │   1. POST /api/ws-ticket│                            │
    │ ───────────────────────►│                            │
    │                         │  2. POST /api/ws-ticket    │
    │                         │  (with session_id cookie)  │
    │                         │ ──────────────────────────►│
    │                         │                            │
    │                         │  3. Create one-time ticket │
    │                         │    (in-memory map,         │
    │                         │     expires after 30s)     │
    │                         │◄───────────────────────────│
    │  4. Returns { ticket }  │                            │
    │◄────────────────────────│                            │
    │                         │                            │
    │  5. new WebSocket(      │                            │
    │     "ws://backend/ws?   │                            │
    │     ticket=xyz...")     │                            │
    │ ────────────────────────────────────────────────────►│
    │                         │                            │
    │                         │  6. Ticket redeemed,       │
    │                         │     connection upgraded    │
    │                         │     via gorilla/websocket  │
    │                         │◄───────────────────────────│
    │                         │                            │
    │  7. ✓ Connected!        │                            │
    │◄─────────────────────────────────────────────────────│