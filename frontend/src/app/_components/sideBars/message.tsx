"use client";

import React, { useEffect, useState, useCallback } from "react";
import { useRouter } from "next/navigation";
import { useSession } from "next-auth/react";
import { getConversations } from "~/app/api/crud/conversations";
import type { ConversationFeedItem } from "~/app/api/crud/conversations";
import { getProfileData } from "~/app/api/crud/getProfile";
import { useChat } from "~/app/_providers/chatProvider";
import { useWS } from "~/app/_providers/ws-provider";

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

function moveConversationToFront(
  list: ConversationFeedItem[],
  convId: number,
  text: string
): ConversationFeedItem[] {
  const idx = list.findIndex((c) => c.id === convId);
  const now = new Date().toISOString();
  if (idx === -1) return list;

  const updated = [...list];
  const old = updated[idx]!;
  const item: ConversationFeedItem = {
    type: old.type,
    id: old.id,
    display_name: old.display_name,
    avatar: old.avatar,
    unread_count: old.unread_count,
    rank: old.rank,
    other_user_id: old.other_user_id,
    member_count: old.member_count,
    last_message: text,
    last_message_at: now,
  };
  updated.splice(idx, 1);
  updated.unshift(item);
  return updated;
}

type FollowUser = { id: number; nickname: string; firstname: string; lastname: string; avatar: string };

interface MessagesSidebarProps {
  onSelect?: (item: ConversationFeedItem) => void;
}

export const MessagesSidebar: React.ComponentType<MessagesSidebarProps> = ({
  onSelect,
}) => {
  const router = useRouter();
  const { data: session } = useSession();
  const currentUserId = session?.user?.id;
  const { selectedChat, selectChat } = useChat();
  const { on, onlineUsers } = useWS();

  const [users, setUsers] = useState<ConversationFeedItem[]>([]);
  const [groups, setGroups] = useState<ConversationFeedItem[]>([]);
  const [following, setFollowing] = useState<FollowUser[]>([]);
  const [followers, setFollowers] = useState<FollowUser[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Load conversations + user's social graph
  useEffect(() => {
    let cancelled = false;

    async function load() {
      const [convResult, profileResult] = await Promise.all([
        getConversations(),
        currentUserId ? getProfileData(currentUserId) : Promise.resolve(null),
      ]);

      if (cancelled) return;

      if ("error" in convResult) {
        setError(convResult.error as string);
        setLoading(false);
        return;
      }

      setUsers(convResult.direct ?? []);
      setGroups(convResult.groups ?? []);

      if (profileResult?.success && profileResult.data) {
        setFollowing(profileResult.data.following ?? []);
        setFollowers(profileResult.data.followers ?? []);
      }

      setLoading(false);
    }

    load();
    return () => {
      cancelled = true;
    };
  }, [currentUserId]);

  const isSelected = useCallback(
    (convId: number, type: "user" | "group") => {
      return selectedChat?.type === type && selectedChat?.id === String(convId);
    },
    [selectedChat]
  );

  // Listen for live messages and update conversation list
  useEffect(() => {
    const unsubs = [
      on("new_message", (data: any) => {
        const convId = Number(data.conversation_id);
        const selected = isSelected(convId, "user");
        setUsers((prev) => {
          const updated = moveConversationToFront(prev, convId, data.text);
          return updated.map((c) =>
            c.id === convId && !selected
              ? { ...c, unread_count: (c.unread_count || 0) + 1 }
              : c
          );
        });
      }),
      on("new_group_message", (data: any) => {
        const groupId = Number(data.group_id);
        const selected = isSelected(groupId, "group");
        setGroups((prev) => {
          const updated = moveConversationToFront(prev, groupId, data.text);
          return updated.map((c) =>
            c.id === groupId && !selected
              ? { ...c, unread_count: (c.unread_count || 0) + 1 }
              : c
          );
        });
      }),
    ];

    return () => unsubs.forEach((fn) => fn());
  }, [on, isSelected]);

  // Build merged user list: conversations first, then non-contacted following + followers
  const buildMergedUsers = useCallback(() => {
    const convUserIds = new Set(users.map((u) => u.other_user_id ?? u.id));
    // Following/followers not yet contacted
    const extraUserIds = new Set<number>();
    const extraUsers: { id: number; nickname: string; display_name: string; avatar: string }[] = [];

    for (const u of [...following, ...followers]) {
      if (extraUserIds.has(u.id) || convUserIds.has(u.id) || u.id === Number(currentUserId)) continue;
      extraUserIds.add(u.id);
      extraUsers.push({
        id: u.id,
        nickname: u.nickname,
        display_name: u.nickname || `${u.firstname} ${u.lastname}`.trim(),
        avatar: u.avatar,
      });
    }

    // Sort extra users alphabetically
    extraUsers.sort((a, b) => a.display_name.localeCompare(b.display_name));

    return { conversations: users, extra: extraUsers };
  }, [users, following, followers, currentUserId]);

  const handleSelect = (item: ConversationFeedItem, type: "user" | "group") => {
    if (type === "user") {
      setUsers((prev) =>
        prev.map((c) => (c.id === item.id ? { ...c, unread_count: 0 } : c))
      );
    } else {
      setGroups((prev) =>
        prev.map((c) => (c.id === item.id ? { ...c, unread_count: 0 } : c))
      );
    }

    selectChat({
      id: String(item.id),
      type: type,
      data: item,
    });

    onSelect?.(item);
    router.push('/messages');
  };

  const handleStartConversation = (targetUser: { id: number; display_name: string; avatar: string }) => {
    const item: ConversationFeedItem = {
      type: "direct",
      id: targetUser.id,
      display_name: targetUser.display_name,
      avatar: targetUser.avatar,
      unread_count: 0,
      rank: 0,
      other_user_id: targetUser.id,
    };
    handleSelect(item, "user");
  };

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

  const { conversations, extra } = buildMergedUsers();

  return (
    <aside className="sidebar2">
      {/* Section 1: Direct conversations (sorted by last message) */}
      {conversations.length > 0 && (
        <>
          <p className="sec-label" style={{ padding: 0, marginBottom: "8px" }}>
            conversations ({conversations.length})
          </p>
          <div style={{ fontSize: "11px", display: "flex", flexDirection: "column", gap: "6px", marginBottom: "12px" }}>
            {conversations.map((user) => {
              const sel =
                selectedChat?.type === "user" && selectedChat?.id === String(user.id);

              return (
                <div
                  key={`direct-${user.id}`}
                  onClick={() => handleSelect(user, "user")}
                  style={{
                    display: "flex",
                    alignItems: "center",
                    gap: "8px",
                    padding: "6px",
                    cursor: "pointer",
                    background: sel
                      ? "var(--color-background-tertiary)"
                      : "var(--color-background-secondary)",
                    border: "0.5px solid var(--color-border-tertiary)",
                    fontWeight: sel ? 500 : "normal",
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
                  <span style={{ color: "var(--color-text-primary)", flex: 1, overflow: "hidden", textOverflow: "ellipsis", whiteSpace: "nowrap" }}>
                    {user.display_name}
                  </span>
                  <span
                    className={onlineUsers.includes(String(user.other_user_id ?? user.id)) ? "online-dot" : "offline-dot"}
                    style={{ flexShrink: 0 }}
                  />
                  {user.unread_count > 0 && (
                    <span
                      style={{
                        fontSize: "9px",
                        color: "var(--color-text-secondary)",
                      }}
                    >
                      {user.unread_count}
                    </span>
                  )}
                  {user.last_message && (
                    <span style={{ fontSize: "8px", color: "var(--color-text-tertiary)", maxWidth: "60px", overflow: "hidden", textOverflow: "ellipsis", whiteSpace: "nowrap" }}>
                      {user.last_message.slice(0, 12)}
                    </span>
                  )}
                </div>
              );
            })}
          </div>
        </>
      )}

      {/* Section 2: Users you can message (sorted alphabetically) */}
      {extra.length > 0 && (
        <>
          <p className="sec-label" style={{ padding: 0, marginBottom: "8px" }}>
            discover ({extra.length})
          </p>
          <div style={{ fontSize: "11px", display: "flex", flexDirection: "column", gap: "6px", marginBottom: "12px" }}>
            {extra.map((u) => (
              <div
                key={`extra-${u.id}`}
                onClick={() => handleStartConversation(u)}
                style={{
                  display: "flex",
                  alignItems: "center",
                  gap: "8px",
                  padding: "6px",
                  cursor: "pointer",
                  background: "var(--color-background-secondary)",
                  border: "0.5px solid var(--color-border-tertiary)",
                }}
              >
                <div
                  style={{
                    width: "16px",
                    height: "16px",
                    borderRadius: "50%",
                    background: colorFor(u.id, AVATAR_COLORS),
                    color: "#993556",
                    fontSize: "8px",
                    display: "flex",
                    alignItems: "center",
                    justifyContent: "center",
                    fontWeight: 600,
                    flexShrink: 0,
                  }}
                >
                  {getInitials(u.display_name)}
                </div>
                <span style={{ color: "var(--color-text-primary)", flex: 1 }}>
                  {u.display_name}
                </span>
                <span
                  className={onlineUsers.includes(String(u.id)) ? "online-dot" : "offline-dot"}
                  style={{ flexShrink: 0 }}
                />
                <i className="ti ti-message-plus" style={{ fontSize: "12px", color: "#D4537E" }} />
              </div>
            ))}
          </div>
        </>
      )}

      {/* If nothing to show */}
      {conversations.length === 0 && extra.length === 0 && (
        <p className="sec-label" style={{ padding: 0, marginBottom: "8px" }}>
          follow users to start messaging
        </p>
      )}

      <div className="divider" style={{ margin: "12px 0" }}></div>

      {/* Section 3: Shared Groups */}
      <p className="sec-label" style={{ padding: 0, marginBottom: "8px" }}>
        shared groups
      </p>
      <div style={{ fontSize: "11px", display: "flex", flexDirection: "column", gap: "6px" }}>
        {groups.length === 0 ? (
          <p style={{ fontSize: "10px", color: "var(--color-text-tertiary)" }}>no groups yet</p>
        ) : (
          groups.map((group) => {
            const sel =
              selectedChat?.type === "group" && selectedChat?.id === String(group.id);

            return (
              <div
                key={`group-${group.id}`}
                onClick={() => handleSelect(group, "group")}
                style={{
                  display: "flex",
                  alignItems: "center",
                  gap: "6px",
                  padding: "6px",
                  cursor: "pointer",
                  background: sel
                    ? "var(--color-background-tertiary)"
                    : "var(--color-background-secondary)",
                  border: "0.5px solid var(--color-border-tertiary)",
                  fontWeight: sel ? 500 : "normal",
                }}
              >
                {group.avatar ? (
                  <img
                    src={group.avatar}
                    alt={group.display_name}
                    style={{
                      width: "20px",
                      height: "20px",
                      borderRadius: "4px",
                      objectFit: "cover",
                      flexShrink: 0,
                    }}
                  />
                ) : (
                  <span
                    style={{
                      width: "6px",
                      height: "6px",
                      background: colorFor(group.id, GROUP_COLORS),
                      flexShrink: 0,
                    }}
                  ></span>
                )}
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
          })
        )}
      </div>
    </aside>
  );
};
