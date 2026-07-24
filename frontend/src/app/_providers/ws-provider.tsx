"use client";

import {
  createContext,
  useContext,
  useEffect,
  useRef,
  useState,
  useCallback,
  type ReactNode,
} from "react";
import { useSession, signOut } from "next-auth/react";
import { useToast } from "~/app/_components/Toast";

type WSEventHandler = (data: any) => void;

type WSContextType = {
  socket: WebSocket | null;
  connected: boolean;
  onlineUsers: string[];
  on: (event: string, handler: WSEventHandler) => () => void;
  send: (eventType: string, data: any) => void;
}

const noop = () => {};

const WSContext = createContext<WSContextType>({
  socket: null,
  connected: false,
  onlineUsers: [],
  on: () => noop,
  send: noop,
});

export function useWS() {
  return useContext(WSContext);
}

// ─── Toast helper ───
const NOTIF_TOAST_MAP: Record<string, { title: string; type: "info" | "success" | "error" | "warning"; message: (data: any) => string }> = {
  like_posts: {
    title: "❤️ New Reaction",
    type: "info",
    message: (d) => `${d.nickname || "Someone"} liked your post`,
  },
  new_comments: {
    title: "💬 New Comment",
    type: "info",
    message: (d) => `${d.nickname || "Someone"} commented on your post`,
  },
  new_posts: {
    title: "📝 New Post",
    type: "info",
    message: (d) => `${d.nickname || "Someone"} posted: "${(d.title || "").slice(0, 40)}"`,
  },
  new_follower: {
    title: "👥 New Follower",
    type: "info",
    message: (d) => `${d.nickname || "Someone"} started following you`,
  },
  follow_request: {
    title: "👥 Follow Request",
    type: "warning",
    message: (d) => `${d.nickname || "Someone"} wants to follow you`,
  },
  follow_accepted: {
    title: "✅ Follow Accepted",
    type: "success",
    message: (d) => `${d.nickname || "Someone"} accepted your follow request`,
  },
  group_invite: {
    title: "📨 Group Invite",
    type: "info",
    message: (d) => `${d.nickname || "Someone"} invited you to join "${d.group_name || "a group"}"`,
  },
  group_join_request: {
    title: "📥 Join Request",
    type: "info",
    message: (d) => `${d.nickname || "Someone"} wants to join your group`,
  },
  group_event: {
    title: "📅 New Event",
    type: "info",
    message: (d) => `${d.group_name ? `New event in "${d.group_name}"` : "A new event was created in your group"}`,
  },
  new_message: {
    title: "💬 New Message",
    type: "info",
    message: (d) => `${d.nickname || "Someone"}: "${(d.text || "").slice(0, 50)}"`,
  },
  new_group_message: {
    title: "💬 Group Message",
    type: "info",
    message: (d) => `${d.nickname || "Someone"} in group: "${(d.text || "").slice(0, 50)}"`,
  },
  // Backend sends 'group_message' from groups.go messages endpoint
  group_message: {
    title: "💬 Group Message",
    type: "info",
    message: (d) => `${d.nickname || "Someone"} sent a message in group`,
  },
};

// Events that should show a toast
const TOAST_EVENTS = Object.keys(NOTIF_TOAST_MAP);

export function WSProvider({ children }: { children: ReactNode }) {
  const { data: session } = useSession();
  const { addToast } = useToast();

  const socketRef = useRef<WebSocket | null>(null);
  const [connected, setConnected] = useState(false);
  const [onlineUsers, setOnlineUsers] = useState<string[]>([]);
  const handlersRef = useRef<Map<string, Set<WSEventHandler>>>(new Map());

  const reconnectAttempt = useRef(0);
  const maxReconnect = 10;

  // ── Listener Registration ──
  const on = useCallback((event: string, handler: WSEventHandler): (() => void) => {
    console.log(`[WS] 🎧 Registering listener for event: "${event}"`);
    const handlers = handlersRef.current;
    if (!handlers.has(event)) {
      handlers.set(event, new Set());
    }
    handlers.get(event)?.add(handler);

    return () => {
      console.log(`[WS] 🔕 Unsubscribing listener for event: "${event}"`);
      handlers.get(event)?.delete(handler);
    };
  }, []);

  // ── Outbound Message ──
  const send = useCallback((eventType: string, data: any) => {
    const ws = socketRef.current;
    if (ws?.readyState === WebSocket.OPEN) {
      console.log(`[WS] 📤 Sending message: "${eventType}"`, data);
      ws.send(JSON.stringify({ event_type: eventType, data }));
    } else {
      console.warn(
        `[WS] ⚠️ Cannot send "${eventType}": WebSocket is not connected (State: ${ws?.readyState ?? "NO_SOCKET"})`
      );
    }
  }, []);

  // ── Toast + currentUserId refs — avoids stale closures in onmessage ──
  const toastRef = useRef<((eventType: string, data: any) => void) | null>(null);
  const currentUserIdRef = useRef<string | undefined>(undefined);

  // Keep refs current on every render
  currentUserIdRef.current = session?.user?.id;

  toastRef.current = (eventType: string, data: any) => {
    const config = NOTIF_TOAST_MAP[eventType];
    if (!config) return;
    addToast({
      title: config.title,
      message: config.message(data),
      type: config.type,
      duration: 4000,
    });
  };

  // ── Connect / Disconnect Lifecycle ──
  useEffect(() => {
    const sessionId = session?.user?.session_id;

    if (!sessionId) {
      console.log("[WS] ⏳ Waiting for user session_id (User not logged in or session loading)...");
      return;
    }

    const backendUrl = process.env.NEXT_PUBLIC_GO_BACKEND_URL;

    if (!backendUrl) {
      console.error("[WS] ❌ Missing NEXT_PUBLIC_GO_BACKEND_URL environment variable.");
      return;
    }

    const wsBase = backendUrl.replace(/^http/, "ws");

    let ws: WebSocket | null = null;
    let reconnectTimer: ReturnType<typeof setTimeout> | null = null;

    async function fetchTicket(): Promise<string | null> {
      try {
        const res = await fetch("/api/ws-ticket", { method: "POST" });
        if (!res.ok) {
          console?.error("[WS] ❌ Ticket fetch failed:", res.status);
          return null;
        }
        const body = await res.json();
        return body.ticket ?? null;
      } catch (err) {
        console.error("[WS] ❌ Ticket fetch error:", err);
        return null;
      }
    }

    async function connect() {
      const ticket = await fetchTicket();
      if (!ticket) {
        console.error("[WS] 🛑 Could not obtain WS ticket, skipping connect.");
        return;
      }

      const wsUrl = `${wsBase}/ws?ticket=${encodeURIComponent(ticket)}`;

      console.log(
        `[WS] 🔌 Connecting to Go WebSocket (Attempt ${reconnectAttempt.current + 1}/${maxReconnect})...`
      );
      console.log(`[WS] 🔗 URL: ${wsUrl}`);

      ws = new WebSocket(wsUrl);

      ws.onopen = () => {
        console.log("[WS] ✅ Connected successfully!");
        reconnectAttempt.current = 0;
        setConnected(true);
        socketRef.current = ws;
      };

      ws.onclose = (event) => {
        console.log(
          `[WS] ❌ Connection closed. Code: ${event.code}, Reason: "${event.reason || "None"}", Clean: ${event.wasClean}`
        );
        setConnected(false);
        socketRef.current = null;

        if (reconnectAttempt.current < maxReconnect) {
          reconnectAttempt.current++;
          const delay = Math.min(1000 * 2 ** reconnectAttempt.current, 30000);
          console.log(`[WS] ⏳ Reconnecting in ${delay / 1000}s...`);
          reconnectTimer = setTimeout(() => { connect(); }, delay);
        } else {
          console.error("[WS] 🛑 Reached maximum reconnection attempts. Giving up.");
        }
      };

      ws.onerror = (err) => {
        console.error("[WS] ⚠️ Connection error encountered:", err);
      };

      ws.onmessage = (event: MessageEvent) => {
        try {
          const msg = JSON.parse(event.data as string);
          const eventType: string = msg.event_type ?? "";
          const data: any = msg.data;

          console.log(`[WS] 📩 Received event: "${eventType}"`, data);

          // Handle force_logout — session was revoked on the backend (e.g. new login)
          if (eventType === "force_logout") {
            console.log("[WS] 🚪 Force logout received — session was revoked on the backend");
            signOut({ callbackUrl: "/login" });
            return;
          }

          // Handle online status events — no early return so custom handlers still fire
          if (eventType === "init" && Array.isArray(data)) {
            console.log("[WS] 👥 Initialized online users list:", data);
            setOnlineUsers(data);
          }

          if (eventType === "client_connect" && typeof data === "string") {
            console.log(`[WS] 🟢 User came online: ${data}`);
            setOnlineUsers((prev) => (prev.includes(data) ? prev : [...prev, data]));
          }

          if (eventType === "client_disconnect" && typeof data === "string") {
            console.log(`[WS] 🔴 User went offline: ${data}`);
            setOnlineUsers((prev) => prev.filter((id) => id !== data));
          }

          // Show toast for notification-type events (skip if current user is the actor)
          if (TOAST_EVENTS.includes(eventType)) {
            const actorId = data?.user_id ?? data?.sender_id ?? data?.actor_id;
            if (actorId && String(actorId) !== currentUserIdRef.current) {
              toastRef.current?.(eventType, data);
            }
            // For events without an actor field (e.g. group_event from backend), always show
            if (!actorId && !data?.user_id && !data?.sender_id) {
              toastRef.current?.(eventType, data);
            }
          }

          // Trigger custom registered event listeners
          const handlers = handlersRef.current.get(eventType);
          if (handlers && handlers.size > 0) {
            console.log(`[WS] 🚀 Dispatching "${eventType}" to ${handlers.size} listener(s)`);
            for (const handler of handlers) {
              handler(data);
            }
          } else {
            console.log(`[WS] ℹ️ No custom listener registered for event "${eventType}"`);
          }
        } catch (e) {
          console.error("[WS] ❌ Error parsing incoming message:", e, "Raw payload:", event.data);
        }
      };
    }

    connect();

    return () => {
      console.log("[WS] 🛑 Cleaning up WebSocket connection (Route change or logout)...");
      if (reconnectTimer) clearTimeout(reconnectTimer);
      reconnectAttempt.current = maxReconnect; // Prevent reconnect on cleanup
      if (ws) {
        ws.onclose = null;
        ws.close();
        console.log("[WS] 🔒 Socket closed cleanly.");
      }
    };
  }, [session?.user?.session_id]); // only reconnect when session_id changes — toast/online logic uses refs

  // ── Periodic online users sync (self-corrects stale entries) ──
  useEffect(() => {
    if (!connected) return;
    const backendUrl = process.env.NEXT_PUBLIC_GO_BACKEND_URL;
    if (!backendUrl) return;

    const sync = async () => {
      try {
        const res = await fetch(`${backendUrl}/api/online-users`);
        const body = await res.json();
        if (Array.isArray(body?.data)) {
          setOnlineUsers(body.data);
        }
      } catch {
        // silent — WS events will eventually correct the list
      }
    };

    // Sync immediately on connect, then every 45s
    sync();
    const interval = setInterval(sync, 45000);
    return () => clearInterval(interval);
  }, [connected]);

  return (
    <WSContext.Provider value={{ socket: socketRef.current, connected, onlineUsers, on, send }}>
      {children}
    </WSContext.Provider>
  );
}
