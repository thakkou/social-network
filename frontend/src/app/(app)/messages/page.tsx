"use client";

import { useState, useEffect, useRef, useCallback } from "react";
import { useSession } from "next-auth/react";
import Image from "next/image";
import { useChat } from "~/app/_providers/chatProvider";
import { useWS } from "~/app/_providers/ws-provider";
import { MessagesSidebar } from "~/app/_components/sideBars/message";
import { EmojiPicker } from "~/app/_components/EmojiPicker";
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
  isSending?: boolean;
}

export default function Chat() {
  const { data: session } = useSession();
  const currentUserId = Number(session?.user?.id ?? 0);

  const { selectedChat, selectChat } = useChat();
  const { send: wsSend, on: wsOn, onlineUsers } = useWS();
  const [message, setMessage] = useState("");

  const [messages, setMessages] = useState<DisplayMessage[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const [otherTyping, setOtherTyping] = useState(false);

  // showUsers: toggles the users sidebar overlay on mobile
  const [showUsers, setShowUsers] = useState(true);

  // Track screen size to render only one MessagesSidebar instance
  const [isDesktop, setIsDesktop] = useState(false);
  useEffect(() => {
    const mq = window.matchMedia("(min-width: 769px)");
    setIsDesktop(mq.matches);
    const handler = (e: MediaQueryListEvent) => setIsDesktop(e.matches);
    mq.addEventListener("change", handler);
    return () => mq.removeEventListener("change", handler);
  }, []);

  const messagesEndRef = useRef<HTMLDivElement>(null);
  const typingTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const isTypingRef = useRef(false);

  const chatData = selectedChat?.data || {};
  const isGroup = selectedChat?.type === "group" || chatData.type === "group";

  const convType: ConversationType = isGroup ? "group" : "direct";
  const convId = selectedChat?.id;
  const displayName = chatData.display_name || "Select a conversation";
  const avatarUrl = chatData.avatar || null;
  const memberCount = chatData.member_count ?? null;
  const initials = getInitials(displayName);

  const fallbackBg = isGroup
    ? colorFor(convId ?? 0, GROUP_COLORS)
    : colorFor(convId ?? 0, AVATAR_COLORS);

  // Auto-scroll
  const scrollToBottom = useCallback(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  }, []);

  useEffect(() => {
    scrollToBottom();
  }, [messages, otherTyping, scrollToBottom]);

  // Fetch messages when selectedChat changes
  useEffect(() => {
    if (!convId) {
      setMessages([]);
      return;
    }

    // If convId matches a user ID (new conversation from discover/profile), skip loading
    const isNewConv = chatData.other_user_id && String(convId) === String(chatData.other_user_id);
    if (isNewConv) {
      setMessages([]);
      setError(null);
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
        const formatted = [...response.messages]
          .reverse()
          .map((msg: ConversationMessage) => ({
            id: msg.id,
            type: (msg.sender_id === currentUserId ? "me" : "them") as "me" | "them",
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
  }, [convType, convId, currentUserId, chatData.other_user_id]);

  // On mobile: when a chat is selected, close the users overlay
  useEffect(() => {
    if (selectedChat) {
      setShowUsers(false);
    }
  }, [selectedChat]);

  // Listen for incoming live messages via WS
  useEffect(() => {
    const unsubs = [
      wsOn("new_message", (data: any) => {
        const incomingConvId = String(data.conversation_id);
        if (incomingConvId !== String(convId)) return;

        setMessages((prev) => {
          const isPendingOptimistic = prev.some(
            (m) => m.senderId === currentUserId && m.isSending
          );
          if (isPendingOptimistic) return prev;
          if (prev.some((m) => m.id === data.message_id)) return prev;

          const type = data.sender_id === currentUserId ? ("me" as const) : ("them" as const);
          return [
            ...prev,
            {
              id: data.message_id,
              type,
              text: data.text,
              senderId: data.sender_id,
              nickname: data.nickname || "unknown",
              createdAt: new Date().toISOString(),
              timeAgo: "just now",
            },
          ];
        });
      }),
      wsOn("new_group_message", (data: any) => {
        const incomingGroupId = String(data.group_id);
        if (incomingGroupId !== String(convId)) return;

        setMessages((prev) => {
          const isPendingOptimistic = prev.some(
            (m) => m.senderId === currentUserId && m.isSending
          );
          if (isPendingOptimistic) return prev;
          if (prev.some((m) => m.id === data.message_id)) return prev;

          const type = data.sender_id === currentUserId ? ("me" as const) : ("them" as const);
          return [
            ...prev,
            {
              id: data.message_id,
              type,
              text: data.text,
              senderId: data.sender_id,
              nickname: data.nickname || "unknown",
              createdAt: new Date().toISOString(),
              timeAgo: "just now",
            },
          ];
        });
      }),
      // Also handle 'group_message' event from groups.go messages endpoint
      wsOn("group_message", (data: any) => {
        const incomingGroupId = String(data.group_id);
        if (incomingGroupId !== String(convId)) return;

        setMessages((prev) => {
          const isPendingOptimistic = prev.some(
            (m) => m.senderId === currentUserId && m.isSending
          );
          if (isPendingOptimistic) return prev;
          if (prev.some((m) => m.id === data.message_id)) return prev;

          const type = data.sender_id === currentUserId ? ("me" as const) : ("them" as const);
          return [
            ...prev,
            {
              id: data.message_id,
              type,
              text: data.text,
              senderId: data.sender_id,
              nickname: data.nickname || "unknown",
              createdAt: new Date().toISOString(),
              timeAgo: "just now",
            },
          ];
        });
      }),
    ];

    return () => unsubs.forEach((fn) => fn());
  }, [wsOn, convId, currentUserId]);

  // Listen for typing indicators from the other user
  useEffect(() => {
    const unsubs = [
      wsOn("typing:start", (data: any) => {
        if (isGroup) return;
        const fromId = String(data.userId);
        if (fromId === String(currentUserId)) return;
        setOtherTyping(true);
      }),
      wsOn("typing:stop", (data: any) => {
        if (isGroup) return;
        const fromId = String(data.userId);
        if (fromId === String(currentUserId)) return;
        setOtherTyping(false);
      }),
    ];

    return () => unsubs.forEach((fn) => fn());
  }, [wsOn, currentUserId, isGroup]);

  // Reset typing state when switching chats
  useEffect(() => {
    setOtherTyping(false);
  }, [convId]);

  // Send typing:start/stop via WS (debounced)
  const handleInputChange = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      setMessage(e.target.value);

      if (!selectedChat || isGroup) return;

      const receiverId = chatData.other_user_id;
      if (!receiverId) return;

      if (!isTypingRef.current) {
        isTypingRef.current = true;
        wsSend("typing:start", {
          conversationId: Number(convId),
          receiverId: Number(receiverId),
          userId: String(currentUserId),
        });
      }

      if (typingTimerRef.current) clearTimeout(typingTimerRef.current);
      typingTimerRef.current = setTimeout(() => {
        isTypingRef.current = false;
        wsSend("typing:stop", {
          conversationId: Number(convId),
          receiverId: Number(receiverId),
          userId: String(currentUserId),
        });
      }, 2000);
    },
    [selectedChat, isGroup, chatData, convId, currentUserId, wsSend]
  );

  // Send typing:stop immediately when sending a message
  async function sendMsg() {
    const val = message.trim();
    if (!val || !selectedChat) return;

    if (typingTimerRef.current) clearTimeout(typingTimerRef.current);
    if (isTypingRef.current) {
      isTypingRef.current = false;
      const receiverId = chatData.other_user_id;
      if (receiverId) {
        wsSend("typing:stop", {
          conversationId: Number(convId),
          receiverId: Number(receiverId),
          userId: String(currentUserId),
        });
      }
    }

    const tempId = Date.now();

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

    // Determine if convId is a real conversation ID or a user ID (from discover/profile)
    const isNewConversation = chatData.other_user_id && String(convId) === String(chatData.other_user_id);

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
            ...(isNewConversation ? {} : { conversation_id: Number(convId as string | number) }),
          };

    const res = await sendMessage(payload);

    if (res.error) {
      setMessages((prev) => prev.filter((m) => m.id !== tempId));
      setError(res.error);
    } else if (res.success && res.data) {
      setMessages((prev) =>
        prev.map((m) =>
          m.id === tempId
            ? { ...m, id: res.data.message_id, isSending: false }
            : m
        )
      );

      // If this was a new conversation, update convId to the real conversation ID
      if (isNewConversation && res.data.conversation_id) {
        selectChat({
          id: String(res.data.conversation_id),
          type: "user",
          data: {
            ...chatData,
          },
        });
      }
    }
  }

  return (
    <main className="main msg-main" style={{ padding: 0, gap: 0 }}>
      {/* Mobile users sidebar overlay — only on mobile */}
      {!isDesktop && (
        <div className={`msgs-sidebar-overlay ${showUsers ? "active" : ""}`}>
          <MessagesSidebar onSelect={() => setShowUsers(false)} />
        </div>
      )}

      {/* Content wrapper: chat + desktop sidebar side-by-side */}
      <div className="msg-content-wrapper">
        {/* Chat panel */}
        <div className="msg-chat-col">
        {/* Mobile header with back + filter button */}
        {selectedChat && (
          <div
            className="mobile-chat-header"
            style={{
              display: "none",
              padding: "6px 8px",
              background: "#211f1c",
              borderBottom: "0.5px solid #3a3733",
            }}
          >
            <div
              style={{
                display: "flex",
                justifyContent: "space-between",
                alignItems: "center",
              }}
            >
              <button
                className="btn btn-g"
                onClick={() => {
                  selectChat(null);
                  setShowUsers(true);
                }}
                style={{
                  display: "flex",
                  alignItems: "center",
                  gap: "4px",
                  fontSize: "11px",
                  padding: "3px 8px",
                }}
              >
                <i className="ti ti-arrow-left" /> back
              </button>
              <button
                className="msgs-filter-btn btn btn-g"
                onClick={() => setShowUsers(true)}
                style={{
                  display: "flex",
                  alignItems: "center",
                  gap: "4px",
                  fontSize: "11px",
                  padding: "3px 8px",
                }}
              >
                <i className="ti ti-users" /> users
              </button>
            </div>
          </div>
        )}

        {/* Chat header with filter/users toggle button */}
        <div
          style={{
            background: "#211f1c",
            borderBottom: "0.5px solid #3a3733",
            padding: "10px 14px",
            display: "flex",
            alignItems: "center",
            gap: "8px",
          }}
        >
          {/* Users toggle — visible on mobile only */}
          <button
            className="msgs-filter-btn"
            onClick={() => setShowUsers(true)}
            style={{
              display: "none",
              alignItems: "center",
              justifyContent: "center",
              width: "28px",
              height: "28px",
              background: "transparent",
              border: "0.5px solid #3a3733",
              borderRadius: "6px",
              cursor: "pointer",
              color: "#a09c94",
              flexShrink: 0,
            }}
            title="Show users"
          >
            <i className="ti ti-menu-2" style={{ fontSize: "15px" }} />
          </button>

          {avatarUrl ? (
            <Image
              src={avatarUrl}
              alt={displayName}
              width={28}
              height={28}
              style={{
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
            <p style={{ fontSize: "10px", color: "#6b6760" }}>
              {isGroup ? (
                <>
                  <span className="online-dot" /> channel active · {memberCount ? `${memberCount} members` : "group"}
                </>
              ) : (
                <>
                  <span className={onlineUsers.includes(String(chatData.other_user_id || convId)) ? "online-dot" : "offline-dot"} />
                  {" "}
                  {onlineUsers.includes(String(chatData.other_user_id || convId)) ? "online now" : "offline"}
                </>
              )}
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
            background: "#1a1917",
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
                    color: "#6b6760",
                  }}
                >
                  Loading messages...
                </div>
              ) : error ? (
                <div
                  style={{
                    margin: "auto",
                    fontSize: "11px",
                    color: "#e53e3e",
                  }}
                >
                  {error}
                </div>
              ) : messages.length === 0 ? (
                <div
                  style={{
                    margin: "auto",
                    fontSize: "11px",
                    color: "#6b6760",
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
                            color: "#e8e4dc",
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

              {/* Typing indicator */}
              {otherTyping && (
                <div
                  style={{
                    display: "flex",
                    alignItems: "center",
                    gap: "6px",
                    paddingLeft: "4px",
                    alignSelf: "flex-start",
                  }}
                >
                  <span
                    style={{
                      fontSize: "10px",
                      color: "#6b6760",
                      fontStyle: "italic",
                    }}
                  >
                    typing
                    <span className="typing-dots">
                      <span>.</span><span>.</span><span>.</span>
                    </span>
                  </span>
                </div>
              )}

              <div ref={messagesEndRef} />
            </>
          ) : (
            <div
              style={{
                margin: "auto",
                fontSize: "12px",
                color: "#6b6760",
                textAlign: "center",
                padding: "2rem",
              }}
            >
              <i className="ti ti-messages" style={{ fontSize: "32px", display: "block", marginBottom: "12px", opacity: 0.4 }} />
              Select a conversation to start chatting<br />
              <span style={{ fontSize: "10px", color: "#a09c94" }}>
                or click the <i className="ti ti-menu-2" /> button to see your contacts
              </span>
            </div>
          )}
        </div>

        {/* Input Bar */}
        <div
          style={{
            background: "#211f1c",
            borderTop: "0.5px solid #3a3733",
            padding: "8px 12px",
            display: "flex",
            gap: "6px",
            alignItems: "center",
          }}
        >
          <EmojiPicker
            onSelect={(emoji) => {
              setMessage((prev) => prev + emoji);
            }}
          />
          <input
            className="inp"
            disabled={!selectedChat || isLoading}
            value={message}
            onChange={handleInputChange}
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

        {/* Desktop sidebar — only on desktop */}
        {isDesktop && (
          <aside className="msg-desktop-sidebar">
            <MessagesSidebar onSelect={() => setShowUsers(false)} />
          </aside>
        )}
      </div>
    </main>
  );
}
