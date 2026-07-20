"use client";

import { useState, useEffect } from "react";
import { useSession } from "next-auth/react";
import { useChat } from "~/app/_providers/chatProvider";
import {
  getConversationById,
  sendMessage,
  type ConversationType,
  type ConversationMessage,
} from "~/app/api/crud/conversations";

const AVATAR_COLORS = ["#FBEAF0", "#EAF3FB", "#EAFBEF", "#FFF3E8", "#F3EAFB"];
const GROUP_COLORS = ["#D4537E", "#1D9E75", "#3B82F6", "#F59E0B", "#8B5CF6"];

function getInitials(name: string) {
  if (!name) return "??";
  return name
    .split(" ")
    .filter(Boolean)
    .slice(0, 2)
    .map((part) => part[0]?.toUpperCase())
    .join("");
}

function colorFor(id: number | string, palette: string[]) {
  const numId = typeof id === "number" ? id : parseInt(id, 10) || 0;
  return palette[Math.abs(numId) % palette.length];
}

function timeAgo(dateStr: string): string {
  if (!dateStr) return "";
  const now = new Date();
  const date = new Date(dateStr);
  const diffMs = now.getTime() - date.getTime();
  const diffMin = Math.floor(diffMs / 60000);
  if (diffMin < 1) return "just now";
  if (diffMin < 60) return `${diffMin}m ago`;
  const diffHrs = Math.floor(diffMin / 60);
  if (diffHrs < 24) return `${diffHrs}h ago`;
  const diffDays = Math.floor(diffHrs / 24);
  if (diffDays < 7) return `${diffDays}d ago`;
  return date.toLocaleDateString();
}

interface DisplayMessage {
  id: number;
  type: "me" | "them";
  text: string;
  senderId: number;
  nickname: string;
  createdAt: string;
  timeAgo: string;
  isSending?: boolean; // Track optimistic state
}

export default function Chat() {
  const { data: session } = useSession();
  const currentUserId = Number(session?.user?.id ?? 0);

  const { selectedChat } = useChat();
  const [message, setMessage] = useState("");

  // Real messages fetched from API
  const [messages, setMessages] = useState<DisplayMessage[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Extract raw payload from context
  const chatData = selectedChat?.data || {};
  const isGroup = selectedChat?.type === "group" || chatData.type === "group";

  // Map conversation type to API string
  const convType: ConversationType = isGroup ? "group" : "direct";
  // Extract conversation ID directly from selectedChat.id
  const convId = selectedChat?.id;
  // Read metadata dynamically from payload
  const displayName = chatData.display_name || "Select a conversation";
  const avatarUrl = chatData.avatar || null;
  const memberCount = chatData.member_count ?? null;
  const initials = getInitials(displayName);

  // Background color calculation if no image avatar is present
  const fallbackBg = isGroup
    ? colorFor(convId ?? 0, GROUP_COLORS)
    : colorFor(convId ?? 0, AVATAR_COLORS);

  // Fetch messages from backend when selectedChat changes
  useEffect(() => {
    if (!convId) {
      setMessages([]);
      return;
    }

    async function loadMessages() {
      setIsLoading(true);
      setError(null);

      const response = await getConversationById(convType, convId as string | number, 0, 30);

      if (response.error) {
        setError(response.error);
        setMessages([]);
      } else if (response.success && response.messages) {
        // Reverse array to display chronologically (oldest top, newest bottom)
        const formatted = [...response.messages]
          .reverse()
          .map((msg: ConversationMessage) => ({
            id: msg.id,
            type: (msg.sender_id === currentUserId ? "me" : "them") as
              | "me"
              | "them",
            text: msg.text,
            senderId: msg.sender_id,
            nickname: msg.nickname || "unknown",
            createdAt: msg.created_at,
            timeAgo: timeAgo(msg.created_at),
          }));
        setMessages(formatted);
      }
      setIsLoading(false);
    }

    loadMessages();
  }, [convType, convId, currentUserId]);

  async function sendMsg() {
    const val = message.trim();
    if (!val || !selectedChat) return;

    const tempId = Date.now();

    // 1. Optimistic UI update
    const tempMsg: DisplayMessage = {
      id: tempId,
      type: "me",
      text: val,
      senderId: currentUserId,
      nickname: "you",
      createdAt: new Date().toISOString(),
      timeAgo: "just now",
      isSending: true,
    };

    setMessages((prev) => [...prev, tempMsg]);
    setMessage("");

    // 2. Build payload based on direct vs group chat
    const payload =
      convType === "group"
        ? {
            type: "group" as const,
            text: val,
            group_id: Number(convId as string | number),
          }
        : {
            type: "direct" as const,
            text: val,
            receiver_id: chatData.other_user_id || Number(convId as string | number),
            conversation_id: Number(convId as string | number),
          };

    // 3. Send message request to backend
    const res = await sendMessage(payload);

    if (res.error) {
      // Rollback optimistic message on failure
      setMessages((prev) => prev.filter((m) => m.id !== tempId));
      setError(res.error);
    } else if (res.success && res.data) {
      // 4. Update message ID with real DB message_id upon success
      setMessages((prev) =>
        prev.map((m) =>
          m.id === tempId
            ? { ...m, id: res.data.message_id, isSending: false }
            : m
        )
      );
    }
  }

  return (
    <main
      className="main"
      style={{
        padding: 0,
        gap: 0,
        display: "flex",
        height: "100%",
        minHeight: "480px",
      }}
    >
      {/* Main Chat View Container */}
      <div
        style={{
          flex: 1,
          display: "flex",
          flexDirection: "column",
          position: "relative",
        }}
      >
        {/* Header */}
        <div
          style={{
            background: "var(--color-background-primary)",
            borderBottom: "0.5px solid var(--color-border-tertiary)",
            padding: "10px 14px",
            display: "flex",
            alignItems: "center",
            gap: "8px",
          }}
        >
          {/* Avatar Rendering */}
          {avatarUrl ? (
            <img
              src={avatarUrl}
              alt={displayName}
              style={{
                width: "28px",
                height: "28px",
                borderRadius: isGroup ? "6px" : "50%",
                objectFit: "cover",
              }}
            />
          ) : (
            <div
              className="av"
              style={{
                width: "28px",
                height: "28px",
                background: fallbackBg,
                color: isGroup ? "#FFFFFF" : "#993556",
                fontSize: "11px",
                borderRadius: isGroup ? "6px" : "50%",
                display: "flex",
                alignItems: "center",
                justifyContent: "center",
                fontWeight: 600,
              }}
            >
              {initials}
            </div>
          )}

          <div>
            <p style={{ fontSize: "12px", fontWeight: 500 }}>{displayName}</p>
            <p style={{ fontSize: "10px", color: "var(--color-text-tertiary)" }}>
              <span className="online-dot" />
              {isGroup
                ? `channel active · ${memberCount ? `${memberCount} members` : "group"}`
                : "online now · websocket"}
            </p>
          </div>
        </div>

        {/* Messages Body */}
        <div
          style={{
            flex: 1,
            padding: "12px",
            display: "flex",
            flexDirection: "column",
            gap: "8px",
            background: "var(--color-background-tertiary)",
            minHeight: "340px",
            overflowY: "auto",
          }}
        >
          {selectedChat ? (
            <>
              {isLoading ? (
                <div
                  style={{
                    margin: "auto",
                    fontSize: "11px",
                    color: "var(--color-text-tertiary)",
                  }}
                >
                  Loading messages...
                </div>
              ) : error ? (
                <div
                  style={{
                    margin: "auto",
                    fontSize: "11px",
                    color: "var(--color-error, #e53e3e)",
                  }}
                >
                  {error}
                </div>
              ) : messages.length === 0 ? (
                <div
                  style={{
                    margin: "auto",
                    fontSize: "11px",
                    color: "var(--color-text-tertiary)",
                  }}
                >
                  No messages here yet. Say hello!
                </div>
              ) : (
                messages.map((msg) => (
                  <div
                    key={msg.id}
                    style={{
                      display: "flex",
                      flexDirection: "column",
                      alignItems: msg.type === "me" ? "flex-end" : "flex-start",
                      opacity: msg.isSending ? 0.6 : 1,
                      maxWidth: "80%",
                      alignSelf: msg.type === "me" ? "flex-end" : "flex-start",
                    }}
                  >
                    {/* Nickname + timestamp for "them" messages */}
                    {msg.type === "them" && (
                      <div
                        style={{
                          display: "flex",
                          alignItems: "center",
                          gap: "6px",
                          marginBottom: "2px",
                          paddingLeft: "4px",
                        }}
                      >
                        <span
                          style={{
                            fontSize: "10px",
                            fontWeight: 600,
                            color: "#D4537E",
                          }}
                        >
                          {msg.nickname}
                        </span>
                        <span
                          style={{
                            fontSize: "9px",
                            color: "#6b6760",
                          }}
                        >
                          {msg.timeAgo}
                        </span>
                      </div>
                    )}
                    {/* Timestamp for "me" messages */}
                    {msg.type === "me" && (
                      <div
                        style={{
                          display: "flex",
                          alignItems: "center",
                          gap: "6px",
                          marginBottom: "2px",
                          paddingRight: "4px",
                        }}
                      >
                        <span
                          style={{
                            fontSize: "9px",
                            color: "#6b6760",
                          }}
                        >
                          {msg.isSending ? "sending..." : msg.timeAgo}
                        </span>
                        <span
                          style={{
                            fontSize: "10px",
                            fontWeight: 600,
                            color: "var(--color-text-primary)",
                          }}
                        >
                          you
                        </span>
                      </div>
                    )}
                    <div
                      className={
                        msg.type === "me" ? "msg-bubble-me" : "msg-bubble-them"
                      }
                    >
                      {msg.text}
                    </div>
                  </div>
                ))
              )}
            </>
          ) : (
            <div
              style={{
                margin: "auto",
                fontSize: "12px",
                color: "var(--color-text-tertiary)",
              }}
            >
              Select a user or group from the sidebar to view conversation
            </div>
          )}
        </div>

        {/* Input Bar */}
        <div
          style={{
            background: "var(--color-background-primary)",
            borderTop: "0.5px solid var(--color-border-tertiary)",
            padding: "8px 12px",
            display: "flex",
            gap: "6px",
            alignItems: "center",
          }}
        >
          <input
            className="inp"
            disabled={!selectedChat || isLoading}
            value={message}
            onChange={(e) => setMessage(e.target.value)}
            onKeyDown={(e) => {
              if (e.key === "Enter") {
                void sendMsg();
              }
            }}
            placeholder={
              selectedChat ? `message ${displayName}...` : "Select a conversation..."
            }
            style={{
              flex: 1,
              fontSize: "12px",
            }}
          />

          <button
            className="btn btn-p"
            onClick={sendMsg}
            disabled={!selectedChat || isLoading || !message.trim()}
          >
            <i className="ti ti-send" />
          </button>
        </div>
      </div>
    </main>
  );
}