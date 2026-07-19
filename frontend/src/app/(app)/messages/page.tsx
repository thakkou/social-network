
"use client";
import { useState } from "react";

export default function Chat() {
  const users = [
    { id: "u1", name: "Selin Rauf", initials: "SR", color: "#FBEAF0", isGroup: false },
    { id: "u2", name: "Adam Smith", initials: "AS", color: "#EAF3FB", isGroup: false },
    { id: "u3", name: "Maya Ali", initials: "MA", color: "#EAFBEF", isGroup: false },
    { id: "u4", name: "John Doe", initials: "JD", color: "#FFF3E8", isGroup: false },
  ];

  const groups = [
    { id: "g1", name: "Engineering Team", initials: "ENG", color: "#EAEFFB", isGroup: true },
    { id: "g2", name: "General Chat", initials: "GEN", color: "#F5EAFB", isGroup: true },
  ];

  // Set default active chat to the first user
  const [activeChat, setActiveChat] = useState(users[0]);
  const [message, setMessage] = useState("");

  const [chats, setChats] = useState({
    u1: [
      { type: "them", text: "hey, tested the websocket handler — looks solid" },
      { type: "me", text: "thanks! just added group broadcast support too" },
      { type: "them", text: "nice. any issue with the gorilla lib version?" },
    ],
    u2: [{ type: "them", text: "hello Adam 👋" }],
    u3: [],
    u4: [],
    g1: [
      { type: "them", text: "System: Deployment to production successful 🚀", senderName: "CI/CD" },
      { type: "them", text: "Anyone reviewing logs for the websocket spike?", senderName: "Selin Rauf" }
    ],
    g2: []
  });

  function sendMsg() {
    const val = message.trim();
    if (!val) return;

    setChats((prev) => ({
      ...prev,
      [activeChat.id]: [
        ...(prev[activeChat.id] || []),
        {
          type: "me",
          text: val,
        },
      ],
    }));

    setMessage("");
  }


  const currentMessages = chats[activeChat.id] || [];

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
          <div
            className="av"
            style={{
              width: "28px",
              height: "28px",
              background: activeChat.color,
              color: activeChat.isGroup ? "#2C5282" : "#993556",
              fontSize: "11px",
              borderRadius: activeChat.isGroup ? "6px" : "50%"
            }}
          >
            {activeChat.initials}
          </div>

          <div>
            <p style={{ fontSize: "12px", fontWeight: 500 }}>{activeChat.name}</p>
            <p style={{ fontSize: "10px", color: "var(--color-text-tertiary)" }}>
              <span className="online-dot" />
              {activeChat.isGroup ? "channel active · group" : "online now · websocket"}
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
          <div style={{ textAlign: "center", fontSize: "10px", color: "var(--color-text-tertiary)" }}>
            today · 14:22
          </div>

          {currentMessages.map((msg, index) => (
            <div key={index} style={{ display: "flex", flexDirection: "column", alignItems: msg.type === "me" ? "flex-end" : "flex-start" }}>
              {/* Add sender alias if it's a group incoming message */}
              {activeChat.isGroup && msg.type === "them" && msg.senderName && (
                <span style={{ fontSize: "9px", color: "var(--color-text-tertiary)", margin: "0 4px 2px 4px" }}>
                  {msg.senderName}
                </span>
              )}
              <div className={msg.type === "me" ? "msg-bubble-me" : "msg-bubble-them"}>
                {msg.text}
              </div>
            </div>
          ))}
        </div>

        {/* Input Bar */}
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
    value={message}
    onChange={(e) => setMessage(e.target.value)}
    onKeyDown={(e) => {
      if (e.key === "Enter") {
        sendMsg();
      }
    }}
    placeholder={`message ${activeChat.name}...`}
    style={{
      flex: 1,
      fontSize: "12px",
    }}
  />

  <button className="btn btn-p" onClick={sendMsg}>
    <i className="ti ti-send" />
  </button>
</div>
      </div>
    </main>
  );
}

