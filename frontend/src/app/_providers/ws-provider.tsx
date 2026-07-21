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
import { useSession } from "next-auth/react";
import { useToast } from "~/app/_components/Toast";

// ─── Types & Context Setup ───
type WSEventHandler = (data: any) => void;

interface WSContextType {
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

export function WSProvider({ children }: { children: ReactNode }) {
  const { data: session } = useSession();
  const { addToast } = useToast();

  const socketRef = useRef<WebSocket | null>(null);
  const [connected, setConnected] = useState(false);
  const [onlineUsers, setOnlineUsers] = useState<string[]>([]);
  const handlersRef = useRef<Map<string, Set<WSEventHandler>>>(new Map());

  const reconnectAttempt = useRef(0);
  const maxReconnect = 10;

  // ── Listener Registration Log ──
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

  // ── Outbound Message Logs ──
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
          console.error("[WS] ❌ Ticket fetch failed:", res.status);
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

          // Handle online status events
          if (eventType === "init" && Array.isArray(data)) {
            console.log("[WS] 👥 Initialized online users list:", data);
            setOnlineUsers(data);
            return;
          }

          if (eventType === "client_connect" && typeof data === "string") {
            console.log(`[WS] 🟢 User came online: ${data}`);
            setOnlineUsers((prev) => (prev.includes(data) ? prev : [...prev, data]));
          }

          if (eventType === "client_disconnect" && typeof data === "string") {
            console.log(`[WS] 🔴 User went offline: ${data}`);
            setOnlineUsers((prev) => prev.filter((id) => id !== data));
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
  }, [session?.user?.session_id]);

  return (
    <WSContext.Provider value={{ socket: socketRef.current, connected, onlineUsers, on, send }}>
      {children}
    </WSContext.Provider>
  );
}