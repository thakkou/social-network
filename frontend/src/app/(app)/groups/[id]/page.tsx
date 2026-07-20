"use client";

import React, { useState, useEffect } from "react";
import Link from "next/link";
import { useParams } from "next/navigation";
import { useSession } from "next-auth/react";
import {
  getGroupPublic,
  getGroupContent,
  createGroupPost,
  toggleGroupPostReaction,
  createGroupPostComment,
  createGroupEvent,
  respondToEvent,
  requestToJoinGroup,
  getPendingRequests,
  acceptJoinRequest,
  rejectJoinRequest,
  leaveGroup,
  type GroupPublic,
  type GroupFeedItem,
  type PendingRequest,
} from "~/app/api/crud/groups";

type FeedFilter = "all" | "posts" | "events";

export default function GroupDetailPage() {
  const params = useParams();
  const groupId = params?.id as string;

  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [group, setGroup] = useState<GroupPublic | null>(null);
  const [isMember, setIsMember] = useState(false);
  const [feed, setFeed] = useState<GroupFeedItem[]>([]);
  const [filter, setFilter] = useState<FeedFilter>("all");

  // ── Create post state ──
  const [newPostText, setNewPostText] = useState("");
  const [newPostTitle, setNewPostTitle] = useState("");
  const [newPostImage, setNewPostImage] = useState<File | null>(null);
  const [newPostPreview, setNewPostPreview] = useState<string | null>(null);
  const { data: session } = useSession();
  const currentUserId = Number(session?.user?.id ?? 0);

  const [posting, setPosting] = useState(false);

  // ── Create event state ──
  const [showEventForm, setShowEventForm] = useState(false);
  const [eventTitle, setEventTitle] = useState("");
  const [eventDesc, setEventDesc] = useState("");
  const [eventTime, setEventTime] = useState("");
  const [creatingEvent, setCreatingEvent] = useState(false);

  // ── Comment state ──
  const [commentText, setCommentText] = useState<Record<number, string>>({});
  const [commentingPost, setCommentingPost] = useState<Record<number, boolean>>({});

  // ── Join / Leave state ──
  const [joining, setJoining] = useState(false);
  const [joinMessage, setJoinMessage] = useState<string | null>(null);
  const [showLeaveConfirm, setShowLeaveConfirm] = useState(false);
  const [pendingRequests, setPendingRequests] = useState<PendingRequest[]>([]);
  const [requestError, setRequestError] = useState<string | null>(null);

  useEffect(() => {
    if (!isMember || !groupId || group?.creator_id === undefined) return;
    const fetchRequests = async () => {
      const res = await getPendingRequests(groupId);
      if (res.success) setPendingRequests(res.data);
    };
    void fetchRequests();
  }, [isMember, groupId, group?.creator_id]);

  const isCreator = isMember && currentUserId > 0 && currentUserId === group?.creator_id;

  useEffect(() => {
    if (!groupId) return;

    const fetchData = async () => {
      setLoading(true);
      setError(null);

      try {
        const [groupRes, contentRes] = await Promise.all([
          getGroupPublic(groupId),
          getGroupContent(groupId),
        ]);

        if (!groupRes.success) {
          setError(groupRes.error ?? "Failed to load group");
          return;
        }

        setGroup(groupRes.data);
        setIsMember(groupRes.isMember ?? false);

        // Only fetch content if member
        if (groupRes.isMember) {
          if (!contentRes.success) {
            setError(contentRes.error ?? "Failed to load content");
            return;
          }
          setFeed(contentRes.data);
        }
      } catch (err) {
        setError(err instanceof Error ? err.message : "Failed to load group");
      } finally {
        setLoading(false);
      }
    };

    void fetchData();
  }, [groupId]);

  const refreshFeed = async () => {
    const contentRes = await getGroupContent(groupId);
    if (contentRes.success) {
      setFeed(contentRes.data);
    }
  };

  // ── Handlers ──

  const handleCreatePost = async () => {
    if (!newPostText.trim() && !newPostTitle.trim()) return;
    setPosting(true);
    try {
      const res = await createGroupPost(groupId, {
        title: newPostTitle.trim() || undefined,
        text: newPostText.trim() || undefined,
        image: newPostImage ?? undefined,
      });
      if (res.success) {
        setNewPostText("");
        setNewPostTitle("");
        setNewPostImage(null);
        setNewPostPreview(null);
        await refreshFeed();
      }
    } finally {
      setPosting(false);
    }
  };

  const handleReaction = async (item: GroupFeedItem, isLike: number) => {
    const res = await toggleGroupPostReaction(groupId, item.id, isLike);
    if (res.success) {
      await refreshFeed();
    }
  };

  const handleComment = async (postId: number) => {
    const text = commentText[postId]?.trim();
    if (!text) return;

    setCommentingPost((prev) => ({ ...prev, [postId]: true }));
    try {
      const res = await createGroupPostComment(groupId, postId, text);
      if (res.success) {
        setCommentText((prev) => ({ ...prev, [postId]: "" }));
        await refreshFeed();
      }
    } finally {
      setCommentingPost((prev) => ({ ...prev, [postId]: false }));
    }
  };

  const handleCreateEvent = async () => {
    if (!eventTitle.trim() || !eventTime) return;
    setCreatingEvent(true);
    try {
      const formattedTime = eventTime.replace("T", " ") + ":00";
      const res = await createGroupEvent(groupId, {
        title: eventTitle.trim(),
        description: eventDesc.trim() || undefined,
        event_time: formattedTime,
      });
      if (res.success) {
        setEventTitle("");
        setEventDesc("");
        setEventTime("");
        setShowEventForm(false);
        await refreshFeed();
      }
    } finally {
      setCreatingEvent(false);
    }
  };

  const handleRequestToJoin = async () => {
    setJoining(true);
    setJoinMessage(null);
    try {
      const res = await requestToJoinGroup(groupId);
      if (res.success) {
        setJoinMessage("request sent!");
      } else {
        setJoinMessage(res.error ?? "failed");
      }
    } catch {
      setJoinMessage("something went wrong");
    } finally {
      setJoining(false);
    }
  };

  const handleEventResponse = async (eventId: number, status: string) => {
    const res = await respondToEvent(groupId, eventId, status);
    if (res.success) {
      await refreshFeed();
    }
  };

  // ── Leave group handlers ──
  const [leaving, setLeaving] = useState(false);

  const handleLeaveGroup = async () => {
    setLeaving(true);
    const res = await leaveGroup(groupId);
    if (res.success) {
      setShowLeaveConfirm(false);
      // Reset to non-member view
      setIsMember(false);
      setFeed([]);
    }
    setLeaving(false);
  };

  // ── Formatting helpers ──

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

  const formatEventTime = (dateStr: string) => {
    const date = new Date(dateStr);
    if (isNaN(date.getTime())) return "";
    return date.toLocaleDateString("en-US", {
      weekday: "short",
      month: "short",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
  };

  const filteredFeed = feed.filter((item) => {
    if (filter === "all") return true;
    if (filter === "posts") return item.type === "post";
    if (filter === "events") return item.type === "event";
    return false;
  });

  const createTime = group?.created_at
    ? new Date(group.created_at).toLocaleDateString("en-US", {
        month: "short",
        year: "numeric",
      })
    : "";

  const initials = group?.title
    ? group.title
        .split(" ")
        .map((w) => w[0])
        .join("")
        .toUpperCase()
        .slice(0, 2)
    : "—";

  // ── Render ──

  if (loading) {
    return (
      <main className="main">
        <BackLink />
        <div className="card" style={{ textAlign: "center", padding: "32px" }}>
          <p style={{ fontSize: 12, color: "#a09c94" }}>Loading group...</p>
        </div>
      </main>
    );
  }

  if (error || !group) {
    return (
      <main className="main">
        <BackLink />
        <div className="card" style={{ textAlign: "center", padding: "32px" }}>
          <p style={{ fontSize: 14, fontWeight: 500, color: "#e8e4dc" }}>
            Group not found
          </p>
          <p style={{ fontSize: 11, color: "#a09c94", marginTop: 4 }}>
            {error || "The group you're looking for doesn't exist or has been removed."}
          </p>
        </div>
      </main>
    );
  }

  return (
    <main className="main">
      <BackLink />

      {/* ── HERO ── */}
      <div className="card" style={{ padding: 0, overflow: "hidden" }}>
        <div
          style={{
            height: 120,
            background: group.background
              ? `url(${group.background}) center/cover`
              : "linear-gradient(135deg,#2e2b27 0%, #272420 60%, #2e1e24 100%)",
            borderBottom: "1px solid #3a3733",
          }}
        />

        <div style={{ padding: "0 20px 20px" }}>
          <div
            className="av"
            style={{
              width: 72,
              height: 72,
              marginTop: -36,
              borderRadius: "50%",
              background: group.logo ? `url(${group.logo}) center/cover` : "#D4537E",
              fontSize: group.logo ? 0 : 26,
              color: "#fff",
              border: "4px solid #272420",
            }}
          >
            {!group.logo && initials}
          </div>

          <div
            style={{
              display: "flex",
              justifyContent: "space-between",
              alignItems: "flex-start",
              marginTop: 12,
              flexWrap: "wrap",
              gap: 16,
            }}
          >
            <div>
              <h2 style={{ fontSize: 22, marginBottom: 6 }}>{group.title}</h2>

              {group.description && (
                <p
                  style={{
                    color: "#a09c94",
                    maxWidth: 650,
                    lineHeight: 1.6,
                    fontSize: 13,
                  }}
                >
                  {group.description}
                </p>
              )}

              <div style={{ marginTop: 12, display: "flex", gap: 8, flexWrap: "wrap" }}>
                <span className="tag tag-teal">public</span>
                {createTime && <span className="tag tag-gray">created {createTime}</span>}
              </div>
            </div>

            <div style={{ display: "flex", gap: 8, flexWrap: "wrap" }}>
              {isMember ? (
                <>
                  <button className="btn btn-p">
                    <i className="ti ti-message" /> chat
                  </button>
                  <button className="btn btn-g">
                    <i className="ti ti-user-plus" /> invite
                  </button>
                  <button
                    className="btn btn-red"
                    onClick={() => setShowLeaveConfirm(true)}
                  >
                    <i className="ti ti-door-exit" /> leave
                  </button>
                </>
              ) : (
                <>
                  <button
                    className="btn btn-p"
                    disabled={joining}
                    onClick={() => void handleRequestToJoin()}
                  >
                    {joining
                      ? "requesting..."
                      : joinMessage
                        ? joinMessage
                        : <><i className="ti ti-user-plus" /> request to join</>}
                  </button>
                  <span className="tag tag-gray" style={{ alignSelf: "center" }}>
                    not a member
                  </span>
                </>
              )}
            </div>
          </div>
        </div>
      </div>

      {/* ── FILTER + CREATE TOOLBAR ── */}
      <div
        className="card"
        style={{
          display: "flex",
          gap: 10,
          padding: 8,
          flexWrap: "wrap",
        }}
      >
        {(["all", "posts", "events"] as FeedFilter[]).map((f) => (
          <button
            key={f}
            className={filter === f ? "btn btn-p" : "btn btn-g"}
            onClick={() => setFilter(f)}
          >
            {f === "all" ? "all" : f}
          </button>
        ))}
        {isMember && (
          <div style={{ marginLeft: "auto", display: "flex", gap: 8 }}>
            <button
              className="btn btn-t"
              onClick={() => setShowEventForm((v) => !v)}
            >
              <i className="ti ti-calendar-plus" /> {showEventForm ? "cancel" : "event"}
            </button>
          </div>
        )}
      </div>

      {/* ── PENDING REQUESTS (creator only) ── */}
      {isCreator && pendingRequests.length > 0 && (
        <div className="card">
          <p style={{ fontSize: 12, color: "#a09c94", marginBottom: 12 }}>
            Pending join requests ({pendingRequests.length})
          </p>
          {pendingRequests.map((req) => (
            <div
              key={req.id}
              style={{
                display: "flex",
                alignItems: "center",
                gap: 10,
                padding: "6px 0",
                borderBottom: "0.5px solid #3a3733",
              }}
            >
              <div
                className="av"
                style={{
                  width: 28,
                  height: 28,
                  borderRadius: "50%",
                  background: req.avatar ? `url(${req.avatar}) center/cover` : "#2e1e24",
                  color: "#D4537E",
                  fontSize: 11,
                }}
              >
                {!req.avatar && (req.nickname?.[0]?.toUpperCase() || req.firstname?.[0]?.toUpperCase() || "?")}
              </div>
              <div style={{ flex: 1, fontSize: 12 }}>
                {req.nickname || `${req.firstname} ${req.lastname}`.trim()}
              </div>
              <button
                className="btn btn-t"
                style={{ fontSize: 10 }}
                onClick={async () => {
                  setRequestError(null);
                  const res = await acceptJoinRequest(groupId, req.user_id);
                  if (res.success) setPendingRequests((prev) => prev.filter((r) => r.id !== req.id));
                  else setRequestError(res.error ?? "accept failed");
                }}
              >
                accept
              </button>
              <button
                className="btn btn-red"
                style={{ fontSize: 10 }}
                onClick={async () => {
                  setRequestError(null);
                  const res = await rejectJoinRequest(groupId, req.user_id);
                  if (res.success) setPendingRequests((prev) => prev.filter((r) => r.id !== req.id));
                  else setRequestError(res.error ?? "reject failed");
                }}
              >
                decline
              </button>
            </div>
          ))}
          {requestError && (
            <p style={{ fontSize: 11, color: "#e07070", marginTop: 8 }}>{requestError}</p>
          )}
        </div>
      )}

      {/* ── CREATE EVENT FORM ── */}
      {showEventForm && (
        <div className="card">
          <p style={{ fontSize: 12, color: "#a09c94", marginBottom: 12 }}>
            Create a new event
          </p>
          <div className="form-row">
            <label className="form-label">Title</label>
            <input
              className="inp"
              value={eventTitle}
              onChange={(e) => setEventTitle(e.target.value)}
              placeholder="Event title"
            />
          </div>
          <div className="form-row">
            <label className="form-label">Description</label>
            <textarea
              className="inp"
              rows={3}
              value={eventDesc}
              onChange={(e) => setEventDesc(e.target.value)}
              placeholder="Event description..."
              style={{ resize: "none" }}
            />
          </div>
          <div className="form-row">
            <label className="form-label">Date & Time</label>
            <input
              className="inp"
              type="datetime-local"
              value={eventTime}
              onChange={(e) => setEventTime(e.target.value)}
            />
          </div>
          <button
            className="btn btn-p"
            disabled={creatingEvent || !eventTitle.trim() || !eventTime}
            onClick={() => void handleCreateEvent()}
          >
            {creatingEvent ? "creating..." : <><i className="ti ti-send" /> create event</>}
          </button>
        </div>
      )}

      {/* ── CREATE POST (members only) ── */}
      {isMember && (
        <div className="card">
          <p style={{ fontSize: 12, color: "#a09c94", marginBottom: 12 }}>
            Create a new group post
          </p>
          <div className="form-row">
            <label className="form-label">Title (optional)</label>
            <input
              className="inp"
              value={newPostTitle}
              onChange={(e) => setNewPostTitle(e.target.value)}
              placeholder="Post title..."
            />
          </div>
          <textarea
            className="inp"
            rows={3}
            value={newPostText}
            onChange={(e) => setNewPostText(e.target.value)}
            placeholder="Share something with the group..."
            style={{ resize: "none" }}
          />
          {/* Image preview */}
          {newPostPreview && (
            <div style={{ position: "relative", marginTop: 8 }}>
              <img
                src={newPostPreview}
                alt="Preview"
                style={{
                  width: "100%",
                  maxHeight: 200,
                  objectFit: "cover",
                  border: "0.5px solid #3a3733",
                }}
              />
              <button
                onClick={() => {
                  setNewPostImage(null);
                  setNewPostPreview(null);
                }}
                style={{
                  position: "absolute",
                  top: 4,
                  right: 4,
                  background: "#2a1818",
                  border: "0.5px solid #7a2c2c",
                  color: "#e07070",
                  cursor: "pointer",
                  width: 22,
                  height: 22,
                  display: "flex",
                  alignItems: "center",
                  justifyContent: "center",
                  fontSize: 12,
                }}
              >
                ✕
              </button>
            </div>
          )}

          <div style={{ display: "flex", justifyContent: "space-between", marginTop: 12 }}>
            <label className="btn btn-g" style={{ cursor: "pointer", display: "flex", alignItems: "center", gap: 4 }}>
              <i className="ti ti-photo" />
              {newPostImage ? "change" : "image"}
              <input
                type="file"
                accept="image/*"
                style={{ display: "none" }}
                onChange={(e) => {
                  const file = e.target.files?.[0];
                  if (file) {
                    setNewPostImage(file);
                    setNewPostPreview(URL.createObjectURL(file));
                  }
                }}
              />
            </label>
            <button
              className="btn btn-p"
              disabled={posting || (!newPostText.trim() && !newPostTitle.trim())}
              onClick={() => void handleCreatePost()}
            >
              {posting ? "posting..." : <><i className="ti ti-send" /> publish</>}
            </button>
          </div>
        </div>
      )}

      {/* ── FEED ── */}
      {filteredFeed.length === 0 && (
        <div className="card" style={{ textAlign: "center", padding: "24px" }}>
          <p style={{ fontSize: 12, color: "#6b6760" }}>
            {filter === "events"
              ? "No events yet. Create one!"
              : filter === "posts"
                ? "No posts yet. Write something!"
                : "No content yet."}
          </p>
        </div>
      )}

      {filteredFeed.map((item) =>
        item.type === "event" ? (
          /* ── EVENT CARD ── */
          <div key={`event-${item.id}`} className="event-card">
            <div style={{ display: "flex", justifyContent: "space-between", marginBottom: 8 }}>
              <strong>{item.title}</strong>
              <span className="tag tag-amber">{formatEventTime(item.event_time)}</span>
            </div>

            {item.description && (
              <p style={{ marginBottom: 10, color: "#a09c94", lineHeight: 1.6 }}>
                {item.description}
              </p>
            )}

            {/* Event creator */}
            {item.author && (
              <div
                style={{
                  display: "flex",
                  alignItems: "center",
                  gap: 6,
                  marginBottom: 10,
                  fontSize: 11,
                  color: "#6b6760",
                }}
              >
                <div
                  className="av"
                  style={{
                    width: 20,
                    height: 20,
                    borderRadius: "50%",
                    background: item.author.avatar
                      ? `url(${item.author.avatar}) center/cover`
                      : "#2e1e24",
                    color: "#D4537E",
                    fontSize: 9,
                  }}
                >
                  {!item.author.avatar &&
                    (item.author.nickname?.[0]?.toUpperCase() ||
                      item.author.firstname?.[0]?.toUpperCase() ||
                      "?")}
                </div>
                <span>
                  {item.author.nickname ||
                    `${item.author.firstname} ${item.author.lastname}`.trim() ||
                    "Unknown"}
                </span>
              </div>
            )}

            {/* Response buttons */}
            <div style={{ display: "flex", gap: 8 }}>
              <button
                className="btn btn-t"
                onClick={() => void handleEventResponse(item.id, "going")}
              >
                <i className="ti ti-check" /> (
                {item.event_responses?.filter((r) => r.status === "going").length || 0} going)
              </button>
              <button
                className="btn btn-g"
                onClick={() => void handleEventResponse(item.id, "not_going")}
              >
                (
                {item.event_responses?.filter((r) => r.status === "not_going").length || 0} not
                going)
              </button>
            </div>

            {/* Responders */}
            {item.event_responses && item.event_responses.length > 0 && (
              <div
                style={{
                  marginTop: 10,
                  display: "flex",
                  gap: 6,
                  flexWrap: "wrap",
                }}
              >
                {item.event_responses.slice(0, 5).map((r) => (
                  <span key={r.user_id} className="tag tag-gray" style={{ fontSize: 10 }}>
                    {r.nickname || r.firstname} — {r.status}
                  </span>
                ))}
                {item.event_responses.length > 5 && (
                  <span className="tag tag-gray" style={{ fontSize: 10 }}>
                    +{item.event_responses.length - 5} more
                  </span>
                )}
              </div>
            )}
          </div>
        ) : (
          /* ── POST CARD ── */
          <div key={`post-${item.id}`} className="card">
            {/* Author */}
            <div style={{ display: "flex", alignItems: "center", gap: 10 }}>
              <div
                className="av"
                style={{
                  width: 36,
                  height: 36,
                  borderRadius: "50%",
                  background: item.author?.avatar
                    ? `url(${item.author.avatar}) center/cover`
                    : "#2e1e24",
                  color: "#D4537E",
                  fontSize: 13,
                }}
              >
                {!item.author?.avatar &&
                  (item.author?.nickname?.[0]?.toUpperCase() ||
                    item.author?.firstname?.[0]?.toUpperCase() ||
                    "?")}
              </div>

              <div>
                <strong style={{ fontSize: 13 }}>
                  {item.author?.nickname ||
                    `${item.author?.firstname || ""} ${item.author?.lastname || ""}`.trim() ||
                    "Unknown"}
                </strong>
                <div style={{ color: "#6b6760", fontSize: 11 }}>
                  {formatTimeAgo(item.created_at)}
                </div>
              </div>
            </div>

            <div className="divider" />

            {item.title && (
              <h3 style={{ marginBottom: 8, fontSize: 15 }}>{item.title}</h3>
            )}

            {item.text && (
              <p
                style={{
                  color: "#a09c94",
                  lineHeight: 1.7,
                  fontSize: 13,
                  whiteSpace: "pre-wrap",
                }}
              >
                {item.text}
              </p>
            )}

            {item.image && (
              <img
                src={item.image}
                alt="Post image"
                style={{
                  width: "100%",
                  maxHeight: 300,
                  objectFit: "cover",
                  marginTop: 10,
                  border: "0.5px solid #3a3733",
                }}
              />
            )}

            {/* Action buttons */}
            <div style={{ display: "flex", gap: 8, marginTop: 14 }}>
              <button
                className={item.is_liked ? "btn btn-p" : "btn btn-g"}
                style={{ display: "flex", alignItems: "center", gap: 4 }}
                onClick={() => void handleReaction(item, 1)}
              >
                <i className="ti ti-thumb-up" /> {item.likes_count}
              </button>
              <button
                className="btn btn-g"
                style={{ display: "flex", alignItems: "center", gap: 4 }}
                onClick={() => void handleReaction(item, -1)}
              >
                <i className="ti ti-thumb-down" /> {item.dislikes_count}
              </button>
              <button
                className="btn btn-g"
                style={{ display: "flex", alignItems: "center", gap: 4 }}
              >
                <i className="ti ti-message-circle" /> {item.comments_count}
              </button>
            </div>

            {/* Comment input */}
            <div className="divider" />
            <div style={{ display: "flex", gap: 6, alignItems: "center" }}>
              <input
                className="inp"
                value={commentText[item.id] || ""}
                onChange={(e) =>
                  setCommentText((prev) => ({ ...prev, [item.id]: e.target.value }))
                }
                placeholder="Write a comment..."
                onKeyDown={(e) => {
                  if (e.key === "Enter") {
                    e.preventDefault();
                    void handleComment(item.id);
                  }
                }}
                style={{ fontSize: 12 }}
              />
              <button
                className="btn btn-p"
                disabled={commentingPost[item.id] || !commentText[item.id]?.trim()}
                onClick={() => void handleComment(item.id)}
                style={{ fontSize: 11, whiteSpace: "nowrap" }}
              >
                {commentingPost[item.id] ? "..." : "send"}
              </button>
            </div>
          </div>
        )
      )}
      {/* ── LEAVE CONFIRMATION MODAL ── */}
      {showLeaveConfirm && (
        <div
          style={{
            position: "fixed",
            inset: 0,
            background: "rgba(0,0,0,0.6)",
            display: "flex",
            alignItems: "center",
            justifyContent: "center",
            zIndex: 100,
          }}
          onClick={() => setShowLeaveConfirm(false)}
        >
          <div
            className="card"
            style={{
              maxWidth: 380,
              width: "90%",
              textAlign: "center",
            }}
            onClick={(e) => e.stopPropagation()}
          >
            <div
              style={{
                width: 48,
                height: 48,
                borderRadius: "50%",
                background: "#3a1e24",
                color: "#D4537E",
                display: "flex",
                alignItems: "center",
                justifyContent: "center",
                margin: "0 auto 12px",
                fontSize: 20,
              }}
            >
              <i className="ti ti-alert-triangle" />
            </div>
            <p
              style={{
                fontSize: 14,
                fontWeight: 600,
                color: "#e8e4dc",
                marginBottom: 8,
              }}
            >
              Leave &quot;{group.title}&quot;?
            </p>
            <p
              style={{
                fontSize: 12,
                color: "#a09c94",
                lineHeight: 1.5,
                marginBottom: 16,
              }}
            >
              All your posts, comments, messages and reactions in this group
              will be permanently deleted. This action cannot be undone.
            </p>
            <div
              style={{
                display: "flex",
                gap: 8,
                justifyContent: "center",
              }}
            >
              <button
                className="btn btn-g"
                style={{ fontSize: 11 }}
                onClick={() => setShowLeaveConfirm(false)}
              >
                cancel
              </button>
              <button
                className="btn btn-red"
                style={{
                  fontSize: 11,
                  display: "flex",
                  alignItems: "center",
                  gap: 4,
                  opacity: leaving ? 0.6 : 1,
                }}
                disabled={leaving}
                onClick={() => void handleLeaveGroup()}
              >
                <i className="ti ti-door-exit" />{" "}
                {leaving ? "leaving..." : "leave group"}
              </button>
            </div>
          </div>
        </div>
      )}
    </main>
  );
}

function BackLink() {
  return (
    <Link
      href="/groups"
      style={{
        display: "flex",
        alignItems: "center",
        gap: 6,
        color: "#a09c94",
        fontSize: 12,
        textDecoration: "none",
        width: "fit-content",
      }}
    >
      <i className="ti ti-arrow-left" />
      back to groups
    </Link>
  );
}
