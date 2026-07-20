"use client";

import { useState, useEffect } from "react";
import { useChat } from "~/app/_providers/chatProvider"; // Adjust path if needed

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

interface Message {
  type: "me" | "them";
  text: string;
  senderName?: string;
}

export default function Chat() {
  const { selectedChat } = useChat();
  const [message, setMessage] = useState("");

  // In-memory chat storage keyed by selection string (e.g. "user-2" or "group-1")
  const [chats, setChats] = useState<Record<string, Message[]>>({
    "user-2": [
      { type: "them", text: "Sounds good, thanks!" }
    ],
    "group-1": [
      { type: "them", text: "Welcome to Gophers United!", senderName: "System" }
    ]
  });

  // Extract raw payload from context
  const chatData = selectedChat?.data || {};
  const isGroup = selectedChat?.type === "group" || chatData.type === "group";
  const chatId = selectedChat?.id ?? String(chatData.id ?? "");
  const activeKey = `${selectedChat?.type || "user"}-${chatId}`;

  // Read metadata dynamically from payload
  const displayName = chatData.display_name || "Select a conversation";
  const avatarUrl = chatData.avatar || null;
  const memberCount = chatData.member_count ?? null;
  const initials = getInitials(displayName);

  // Background color calculation if no image avatar is present
  const fallbackBg = isGroup 
    ? colorFor(chatId, GROUP_COLORS) 
    : colorFor(chatId, AVATAR_COLORS);

  function sendMsg() {
    const val = message.trim();
    if (!val || !selectedChat) return;

    setChats((prev) => ({
      ...prev,
      [activeKey]: [
        ...(prev[activeKey] || []),
        {
          type: "me",
          text: val,
        },
      ],
    }));

    setMessage("");
  }

  const currentMessages = chats[activeKey] || [];

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
          {/* Avatar Rendering: Image or Styled Initials Box */}
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
              <div
                style={{
                  textAlign: "center",
                  fontSize: "10px",
                  color: "var(--color-text-tertiary)",
                }}
              >
                today · 14:22
              </div>

              {currentMessages.length === 0 ? (
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
                currentMessages.map((msg, index) => (
                  <div
                    key={index}
                    style={{
                      display: "flex",
                      flexDirection: "column",
                      alignItems: msg.type === "me" ? "flex-end" : "flex-start",
                    }}
                  >
                    {isGroup && msg.type === "them" && msg.senderName && (
                      <span
                        style={{
                          fontSize: "9px",
                          color: "var(--color-text-tertiary)",
                          margin: "0 4px 2px 4px",
                        }}
                      >
                        {msg.senderName}
                      </span>
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
            disabled={!selectedChat}
            value={message}
            onChange={(e) => setMessage(e.target.value)}
            onKeyDown={(e) => {
              if (e.key === "Enter") {
                sendMsg();
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
            disabled={!selectedChat}
          >
            <i className="ti ti-send" />
          </button>
        </div>
      </div>
    </main>
  );
}