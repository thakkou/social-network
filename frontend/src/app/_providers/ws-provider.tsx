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

// ─── Types ───

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
  on: (() => noop) as (event: string, handler: WSEventHandler) => () => void,
  send: noop,
});

export function useWS() {
  return useContext(WSContext);
}

// ─── Event label helpers ───

function toastLabel(eventType: string): string {
  switch (eventType) {
    case "new_message":
      return "💬 New message";
    case "new_group_message":
      return "💬 Group message";
    case "like_posts":
      return "❤️ Post liked";
    case "new_comments":
      return "💬 New comment";
    case "new_posts":
      return "📝 New post";
    case "group_event":
      return "📅 Group event";
    case "group_invite":
      return "👋 Group invite";
    case "group_join_request":
      return "🔔 Join request";
    case "client_connect":
      return "🟢 User online";
    case "client_disconnect":
      return "🔴 User offline";
    default:
      return "🔔 Notification";
  }
}

function toastMessage(eventType: string, data: any): string | undefined {
  if (typeof data === "string") {
    // For client_connect / client_disconnect — data is the user ID string
    return data ? `User #${data}` : undefined;
  }
  if (data?.text) return data.text.slice(0, 80);
  if (data?.nickname) return `from ${data.nickname}`;
  if (data?.title) return data.title;
  return undefined;
}

// ─── Provider ───

export function WSProvider({ children }: { children: ReactNode }) {
  const { data: session } = useSession();
  const { addToast } = useToast();

  const socketRef = useRef<WebSocket | null>(null);
  const [connected, setConnected] = useState(false);
  const [onlineUsers, setOnlineUsers] = useState<string[]>([]);
  const handlersRef = useRef<Map<string, Set<WSEventHandler>>>(new Map());

  // Reconnect info
  const reconnectAttempt = useRef(0);
  const maxReconnect = 10;

  const on = useCallback((event: string, handler: WSEventHandler): (() => void) => {
    const handlers = handlersRef.current;
    if (!handlers.has(event)) {
      handlers.set(event, new Set());
    }
    const eventHandlers = handlers.get(event);
    if (eventHandlers) {
      eventHandlers.add(handler);
    }
    return () => {
      handlers.get(event)?.delete(handler);
    };
  }, []);

  const send = useCallback((eventType: string, data: any) => {
    const ws = socketRef.current;
    if (ws?.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify({ event_type: eventType, data }));
    }
  }, []);

  // ── Connect / disconnect based on session ──
  useEffect(() => {
    if (!session?.user?.session_id) return;

    const backendUrl = process.env.NEXT_PUBLIC_GO_BACKEND_URL || process.env.GO_BACKEND_URL;
    if (!backendUrl) return;

    const wsBase = backendUrl.replace(/^http/, "ws");
    const wsUrl = `${wsBase}/ws?session_id=${session.user.session_id}`;

    let ws: WebSocket | null = null;
    let reconnectTimer: ReturnType<typeof setTimeout> | null = null;

    function connect() {
      ws = new WebSocket(wsUrl);

      ws.onopen = () => {
        reconnectAttempt.current = 0;
        setConnected(true);
        socketRef.current = ws;
      };

      ws.onclose = () => {
        setConnected(false);
        socketRef.current = null;

        // Auto-reconnect
        if (reconnectAttempt.current < maxReconnect) {
          reconnectAttempt.current++;
          const delay = Math.min(1000 * 2 ** reconnectAttempt.current, 30000);
          reconnectTimer = setTimeout(connect, delay);
        }
      };

      ws.onmessage = (event: MessageEvent) => {
        try {
          const msg: Record<string, any> = JSON.parse(event.data as string);
          const eventType: string = msg.event_type ?? "";
          const data: any = msg.data;

          // ── Handle built-in events ──

          // "init" — contains the online user list
          if (eventType === "init" && Array.isArray(data)) {
            setOnlineUsers(data);
            return;
          }

          // client_connect / client_disconnect — update online user list
          if (eventType === "client_connect" && typeof data === "string") {
            const uid: string = data;
            setOnlineUsers((prev) =>
              prev.includes(uid) ? prev : [...prev, uid]
            );
          }
          if (eventType === "client_disconnect" && typeof data === "string") {
            const uid: string = data;
            setOnlineUsers((prev) => prev.filter((id) => id !== uid));
          }

          // ── Show a toast for relevant events ──
          const skipToasts = ["typing:start", "typing:stop", "test_event", "client_connect", "client_disconnect"];
          if (!skipToasts.includes(eventType)) {
            addToast({
              type: "info",
              title: toastLabel(eventType),
              message: toastMessage(eventType, data),
              duration: 4000,
            });
          }

          // ── Dispatch to custom handlers ──
          const handlers = handlersRef.current.get(eventType);
          if (handlers) {
            for (const handler of handlers) {
              try {
                handler(data);
              } catch (e) {
                console.error(`[WS] handler error for ${eventType}:`, e);
              }
            }
          }
        } catch (e) {
          console.error("[WS] error parsing message:", e);
        }
      };

      ws.onerror = () => {
        // onclose will fire after this
      };
    }

    connect();

    return () => {
      if (reconnectTimer) clearTimeout(reconnectTimer);
      reconnectAttempt.current = maxReconnect; // prevent reconnect on unmount
      if (ws) {
        ws.onclose = null; // prevent reconnect logic during cleanup
        ws.close();
      }
    };
  }, [session?.user?.session_id, addToast]);

  return (
    <WSContext.Provider value={{ socket: socketRef.current, connected, onlineUsers, on, send }}>
      {children}
    </WSContext.Provider>
  );
}
