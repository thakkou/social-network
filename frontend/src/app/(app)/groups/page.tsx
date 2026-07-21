"use client";

import { useState, useEffect } from "react";
import Link from "next/link";
import { useSession } from "next-auth/react";
import { getUserGroups, getGroupContent, type GroupPublic, type GroupFeedItem } from "~/app/api/crud/groups";
import SearchInput from "~/app/_components/SearchInput";
import CreateGroupForm from "~/app/_components/Forms/CreateGroupForm";

export default function Groups() {
  const { data: session } = useSession();
  const [myGroups, setMyGroups] = useState<GroupPublic[]>([]);
  const [groupFeeds, setGroupFeeds] = useState<{ group: GroupPublic; items: GroupFeedItem[] }[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const load = async () => {
      setLoading(true);
      const res = await getUserGroups();
      if (res.success) {
        setMyGroups(res.data);

        // Fetch content from each group in parallel
        const feeds = await Promise.all(
          res.data.map(async (g) => {
            const contentRes = await getGroupContent(String(g.id), { limit: 5 });
            return {
              group: g,
              items: contentRes.success ? contentRes.data : [],
            };
          })
        );
        setGroupFeeds(feeds);
      }
      setLoading(false);
    };
    void load();
  }, []);

  const formatTimeAgo = (dateStr: string) => {
    const date = new Date(dateStr);
    const now = new Date();
    const diffMs = now.getTime() - date.getTime();
    const diffMins = Math.floor(diffMs / 60000);
    const diffHrs = Math.floor(diffMins / 60);
    const diffDays = Math.floor(diffHrs / 24);
    if (diffMins < 1) return "just now";
    if (diffMins < 60) return `${diffMins}m ago`;
    if (diffHrs < 24) return `${diffHrs}h ago`;
    if (diffDays < 7) return `${diffDays}d ago`;
    return date.toLocaleDateString("en-US", { month: "short", day: "numeric" });
  };

  return (
    <main className="main">
      <CreateGroupForm />

      <div style={{ marginBottom: "16px" }}>
        <SearchInput placeholder="Search groups..." typeSearch="groups" />
      </div>

      {loading ? (
        <div className="card" style={{ textAlign: "center", padding: "24px" }}>
          <p style={{ fontSize: 12, color: "#a09c94" }}>Loading your groups...</p>
        </div>
      ) : groupFeeds.length === 0 ? (
        <div className="card" style={{ textAlign: "center", padding: "24px" }}>
          <p style={{ fontSize: 14, color: "#e8e4dc" }}>Not a member of any group yet</p>
          <p style={{ fontSize: 11, color: "#6b6760", marginTop: 4 }}>
            Create a group or accept an invite to see content here.
          </p>
        </div>
      ) : (
        groupFeeds.map(({ group, items }) => (
          <div key={group.id} className="card">
            {/* Group header */}
            <Link
              href={`/groups/${group.id}`}
              style={{ textDecoration: "none" }}
            >
              <div
                style={{
                  display: "flex",
                  alignItems: "center",
                  gap: 8,
                  marginBottom: items.length > 0 ? 12 : 0,
                }}
              >
                <div
                  className="av"
                  style={{
                    width: 24,
                    height: 24,
                    borderRadius: "50%",
                    background: group.logo
                      ? `url(${group.logo}) center/cover`
                      : "#D4537E",
                    fontSize: group.logo ? 0 : 10,
                    color: "#fff",
                    flexShrink: 0,
                  }}
                >
                  {!group.logo &&
                    group.title
                      .split(" ")
                      .map((w) => w[0])
                      .join("")
                      .toUpperCase()
                      .slice(0, 2)}
                </div>
                <div style={{ flex: 1 }}>
                  <p style={{ fontSize: 13, fontWeight: 600, color: "#e8e4dc", margin: 0 }}>
                    {group.title}
                  </p>
                  {group.description && (
                    <p style={{ fontSize: 11, color: "#6b6760", margin: 0 }}>
                      {group.description}
                    </p>
                  )}
                </div>
                <span className="tag tag-pink" style={{ fontSize: 9 }}>
                  member
                </span>
              </div>
            </Link>

            {/* Feed items from this group */}
            {items.length > 0 ? (
              <div style={{ display: "flex", flexDirection: "column", gap: 8 }}>
                {items.map((item) => (
                  <div key={`${item.type}-${item.id}`} className="event-card" style={{ fontSize: 12 }}>
                    {item.type === "post" ? (
                      <div>
                        <div style={{ display: "flex", justifyContent: "space-between", marginBottom: 4 }}>
                          <strong style={{ fontSize: 12 }}>
                            {item.title || "Untitled"}
                          </strong>
                          <span style={{ fontSize: 10, color: "#6b6760" }}>
                            {formatTimeAgo(item.created_at)}
                          </span>
                        </div>
                        {item.text && (
                          <p style={{ fontSize: 11, color: "#a09c94", margin: 0 }}>
                            {item.text.slice(0, 150)}
                            {item.text.length > 150 ? "..." : ""}
                          </p>
                        )}
                        <div style={{ display: "flex", gap: 6, marginTop: 6, fontSize: 10, color: "#6b6760" }}>
                          <span>
                            <i className="ti ti-thumb-up" /> {item.likes_count || 0}
                          </span>
                          <span>
                            <i className="ti ti-message-circle" /> {item.comments_count || 0}
                          </span>
                        </div>
                      </div>
                    ) : (
                      <div>
                        <div style={{ display: "flex", justifyContent: "space-between", marginBottom: 4 }}>
                          <strong style={{ fontSize: 12 }}>{item.title}</strong>
                          <span className="tag tag-amber" style={{ fontSize: 9 }}>
                            {item.event_time
                              ? new Date(item.event_time).toLocaleDateString("en-US", {
                                  month: "short",
                                  day: "numeric",
                                })
                              : ""}
                          </span>
                        </div>
                        {item.description && (
                          <p style={{ fontSize: 11, color: "#a09c94", margin: 0 }}>
                            {item.description.slice(0, 120)}
                            {item.description.length > 120 ? "..." : ""}
                          </p>
                        )}
                      </div>
                    )}
                  </div>
                ))}
              </div>
            ) : (
              <p style={{ fontSize: 11, color: "#6b6760", margin: 0 }}>
                No recent activity in this group.
              </p>
            )}

            {/* Quick actions */}
            <div style={{ display: "flex", gap: 6, marginTop: 10 }}>
              <Link
                className="btn btn-p"
                href={`/groups/${group.id}`}
                style={{ fontSize: 11, display: "flex", alignItems: "center", gap: 3 }}
              >
                <i className="ti ti-arrow-right" style={{ fontSize: 12 }} /> view group
              </Link>
              <Link
                className="btn btn-g"
                href={`/messages`}
                style={{ fontSize: 11, display: "flex", alignItems: "center", gap: 3 }}
              >
                <i className="ti ti-message" style={{ fontSize: 12 }} /> chat
              </Link>
            </div>
          </div>
        ))
      )}
    </main>
  );
}
