"use client";

import React, { useEffect, useState } from "react";
import { getConversations, ConversationFeedItem } from "~/app/api/crud/conversations";

const AVATAR_COLORS = ["#FBEAF0", "#EAF3FB", "#EAFBEF", "#FFF3E8", "#F3EAFB"];
const GROUP_COLORS = ["#D4537E", "#1D9E75", "#3B82F6", "#F59E0B", "#8B5CF6"];

function getInitials(name: string) {
  return name
    .split(" ")
    .filter(Boolean)
    .slice(0, 2)
    .map((part) => part[0]?.toUpperCase())
    .join("");
}

function colorFor(id: number, palette: string[]) {
  return palette[Math.abs(id) % palette.length];
}

interface MessagesSidebarProps {
  activeId?: string | number;
  onSelect?: (item: ConversationFeedItem) => void;
}

export const MessagesSidebar: React.ComponentType<MessagesSidebarProps> = ({
  activeId,
  onSelect,
}) => {
  const [users, setUsers] = useState<ConversationFeedItem[]>([]);
  const [groups, setGroups] = useState<ConversationFeedItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let cancelled = false;

    async function load() {
      const result = await getConversations();

      if (cancelled) return;

      if ("error" in result) {
        setError(result.error as string);
        setLoading(false);
        return;
      }

      setUsers(result.direct ?? []);
      setGroups(result.groups ?? []);
      setLoading(false);
    }

    load();
    return () => {
      cancelled = true;
    };
  }, []);

  if (loading) {
    return (
      <aside className="sidebar2">
        <p className="sec-label">loading...</p>
      </aside>
    );
  }

  if (error) {
    return (
      <aside className="sidebar2">
        <p className="sec-label" style={{ color: "var(--color-text-secondary)" }}>
          {error}
        </p>
      </aside>
    );
  }

  return (
    <aside className="sidebar2">
      {/* Section 1: Users */}
      <p className="sec-label" style={{ padding: 0, marginBottom: "8px" }}>
        users
      </p>
      <div style={{ fontSize: "11px", display: "flex", flexDirection: "column", gap: "6px" }}>
        {users.map((user) => {
          const isSelected = activeId === user.id;
          return (
            <div
              key={`direct-${user.id}`}
              onClick={() => onSelect?.(user)}
              style={{
                display: "flex",
                alignItems: "center",
                gap: "8px",
                padding: "6px",
                cursor: "pointer",
                background: isSelected
                  ? "var(--color-background-tertiary)"
                  : "var(--color-background-secondary)",
                border: "0.5px solid var(--color-border-tertiary)",
                fontWeight: isSelected ? 500 : "normal",
              }}
            >
              <div
                style={{
                  width: "16px",
                  height: "16px",
                  borderRadius: "50%",
                  background: colorFor(user.id, AVATAR_COLORS),
                  color: "#993556",
                  fontSize: "8px",
                  display: "flex",
                  alignItems: "center",
                  justifyContent: "center",
                  fontWeight: 600,
                  flexShrink: 0,
                }}
              >
                {getInitials(user.display_name)}
              </div>
              <span style={{ color: "var(--color-text-primary)" }}>
                {user.display_name}
              </span>
              {user.unread_count > 0 && (
                <span
                  style={{
                    marginLeft: "auto",
                    fontSize: "9px",
                    color: "var(--color-text-secondary)",
                  }}
                >
                  {user.unread_count}
                </span>
              )}
            </div>
          );
        })}
      </div>

      <div className="divider" style={{ margin: "12px 0" }}></div>

      {/* Section 2: Shared Groups */}
      <p className="sec-label" style={{ padding: 0, marginBottom: "8px" }}>
        shared groups
      </p>
      <div style={{ fontSize: "11px", display: "flex", flexDirection: "column", gap: "6px" }}>
        {groups.map((group) => {
          const isSelected = activeId === group.id;
          return (
            <div
              key={`group-${group.id}`}
              onClick={() => onSelect?.(group)}
              style={{
                display: "flex",
                alignItems: "center",
                gap: "6px",
                padding: "6px",
                cursor: "pointer",
                background: isSelected
                  ? "var(--color-background-tertiary)"
                  : "var(--color-background-secondary)",
                border: "0.5px solid var(--color-border-tertiary)",
                fontWeight: isSelected ? 500 : "normal",
              }}
            >
              <span
                style={{
                  width: "6px",
                  height: "6px",
                  background: colorFor(group.id, GROUP_COLORS),
                  flexShrink: 0,
                }}
              ></span>
              <span style={{ color: "var(--color-text-primary)" }}>
                {group.display_name}
              </span>
              {group.unread_count > 0 && (
                <span
                  style={{
                    marginLeft: "auto",
                    fontSize: "9px",
                    color: "var(--color-text-secondary)",
                  }}
                >
                  {group.unread_count}
                </span>
              )}
            </div>
          );
        })}
      </div>
    </aside>
  );
};