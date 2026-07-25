# social-network

# floow ws-ticket
[Browser]                          [Next.js]                  [Go Backend]
    │                                  │                          │
    │  useSession()                    │                          │
    │  gets session_id                 │                          │
    │────────────────────────────────────────► POST /api/ws-ticket  │
    │                                  │     (with session cookie) │
    │                                  │     Creates one-time ticket│
    │◄──────────────────────────────────────── { ticket: "xyz123" }│
    │                                  │                          │
    │  new WebSocket("ws://.../ws      │                          │
    │    ?ticket=xyz123")              │                          │
    │──────────────────────────────────────── Upgrade request ────►│
    │                                  │     Redeems ticket        │
    │                                  │     Gets userId           │
    │                                  │     Upgrades to WS        │
    │◄──────────────────────────────────── WebSocket connection ───│
    │                                  │                          │
    │  WS sends: {"event_type":        │                          │
    │    "init", "data": [onlineIDs]}  │                          │
    │◄─────────────────────────────────────────────────────────────│
