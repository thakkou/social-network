"use client";

import { useState, useEffect, useCallback, useRef } from "react";
import Image from "next/image";
import { useRouter } from "next/navigation";
import { useSession } from "next-auth/react";
import {
  getFeedPosts,
  likePost,
  dislikePost,
  createPost,
  type FeedPost,
} from "~/app/_services/crud/post";
import { getProfileData } from "~/app/_services/crud/getProfile";

const CATEGORIES = [
  "General", "Lifestyle", "Health & Fitness", "Travel",
  "Food & Cooking", "Education", "Business", "Finance",
  "Entertainment", "Sports", "Personal Dev", "Culture", "News",
];

type InviteUser = {
  id: number;
  nickname: string;
  firstname: string;
  lastname: string;
  avatar: string;
}

export default function Home() {
  const router = useRouter();
  const [posts, setPosts] = useState<FeedPost[]>([]);
  const [loading, setLoading] = useState(true);

  // Create post form
  const [showForm, setShowForm] = useState(false);
  const [formTitle, setFormTitle] = useState("");
  const [formText, setFormText] = useState("");
  const [formPrivacy, setFormPrivacy] = useState("public");
  const [formCategories, setFormCategories] = useState<string[]>(["General"]);
  const [formImage, setFormImage] = useState<File | null>(null);
  const [formImagePreview, setFormImagePreview] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const { data: session } = useSession();
  const currentUserId = session?.user?.id;

  // ── Allowed users for private posts ──
  const [selectedUsers, setSelectedUsers] = useState<InviteUser[]>([]);
  const [inviteQuery, setInviteQuery] = useState("");
  const [searchResults, setSearchResults] = useState<InviteUser[]>([]);
  const [searchOpen, setSearchOpen] = useState(false);
  const [connections, setConnections] = useState<InviteUser[]>([]);
  const debounceRef = useRef<NodeJS.Timeout | null>(null);

  // ── Category filter state ──
  const [filterCategories, setFilterCategories] = useState<string[]>([]);
  const [hasMore, setHasMore] = useState(true);
  const [loadingMore, setLoadingMore] = useState(false);
  const sentinelRef = useRef<HTMLDivElement>(null);

  const fetchPosts = useCallback(async (selectedCategories?: string[], reset = true) => {
    if (reset) {
      setLoading(true);
      setHasMore(true);
    }
    const cats = selectedCategories ?? filterCategories;
    const lastId = reset ? 0 : (posts.length > 0 ? posts[posts.length - 1]!.id : 0);
    const res = await getFeedPosts({
      limit: 20,
      last_id: lastId,
      ...(cats.length > 0 ? { categories: cats } : {}),
    });
    if (res.success) {
      if (reset) {
        setPosts(res.data);
      } else {
        setPosts((prev) => [...prev, ...res.data]);
      }
      if (res.data.length < 20) {
        setHasMore(false);
      }
    }
    if (reset) setLoading(false);
    setLoadingMore(false);
  }, [filterCategories, posts.length]);

  // ── Infinite scroll observer ──
  useEffect(() => {
    const sentinel = sentinelRef.current;
    if (!sentinel) return;
    const observer = new IntersectionObserver(
      (entries) => {
        if (entries[0]?.isIntersecting && hasMore && !loadingMore && !loading) {
          setLoadingMore(true);
          void fetchPosts(undefined, false);
        }
      },
      { rootMargin: "200px" }
    );
    observer.observe(sentinel);
    return () => observer.disconnect();
  }, [hasMore, loadingMore, loading, fetchPosts]);

  const toggleFilterCategory = (cat: string) => {
    setFilterCategories((prev) => {
      const next = prev.includes(cat)
        ? prev.filter((c) => c !== cat)
        : [...prev, cat];
      void fetchPosts(next, true);
      return next;
    });
  };

  // ── Fetch connections (followers + following) when form opens ──
  useEffect(() => {
    if (!showForm || !currentUserId) return;
    const load = async () => {
      const res = await getProfileData(currentUserId);
      if (res.success) {
        const merged = new Map<number, InviteUser>();
        for (const u of [...(res.data.followers ?? []), ...(res.data.following ?? [])]) {
          if (!merged.has(u.id)) {
            merged.set(u.id, {
              id: u.id,
              nickname: u.nickname || `${u.firstname} ${u.lastname}`.trim(),
              firstname: u.firstname,
              lastname: u.lastname,
              avatar: u.avatar || "",
            });
          }
        }
        setConnections(Array.from(merged.values()));
      }
    };
    void load();
  }, [showForm, currentUserId]);

  // ── Filter connections by name (client-side, no API call) ──
  useEffect(() => {
    if (debounceRef.current) clearTimeout(debounceRef.current);
    if (!inviteQuery.trim() || formPrivacy !== "private") {
      setSearchResults([]);
      setSearchOpen(false);
      return;
    }
    const q = inviteQuery.trim().toLowerCase();
    debounceRef.current = setTimeout(() => {
      const filtered = connections.filter(
        (u) =>
          !selectedUsers.some((s) => s.id === u.id) &&
          (u.nickname.toLowerCase().includes(q) ||
            `${u.firstname} ${u.lastname}`.toLowerCase().includes(q))
      );
      setSearchResults(filtered);
      setSearchOpen(filtered.length > 0);
    }, 300);
    return () => {
      if (debounceRef.current) clearTimeout(debounceRef.current);
    };
  }, [inviteQuery, selectedUsers, formPrivacy, connections]);

  const addUser = (user: InviteUser) => {
    setSelectedUsers((prev) => [...prev, user]);
    setInviteQuery("");
    setSearchResults([]);
    setSearchOpen(false);
  };

  const removeUser = (userId: number) => {
    setSelectedUsers((prev) => prev.filter((u) => u.id !== userId));
  };

  useEffect(() => {
    void fetchPosts(undefined, true);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const handleCreatePost = async () => {
    if (!formTitle.trim() || !formText.trim()) return;
    if (formPrivacy === "private" && selectedUsers.length === 0) return;
    setSubmitting(true);

    const fd = new FormData();
    fd.set("title", formTitle.trim());
    fd.set("text", formText.trim());
    fd.set("privacy", formPrivacy);
    formCategories.forEach((cat) => fd.append("categories", cat));
    if (formImage) fd.set("image", formImage);
    if (formPrivacy === "private" && selectedUsers.length > 0) {
      fd.set("allowed_user_ids", selectedUsers.map((u) => u.id).join(","));
    }

    const res = await createPost(fd);
    if (res.success) {
      setShowForm(false);
      setFormTitle("");
      setFormText("");
      setFormPrivacy("public");
      setFormCategories(["General"]);
      setFormImage(null);
      setFormImagePreview(null);
      setSelectedUsers([]);
      setInviteQuery("");
      void fetchPosts();
    }
    setSubmitting(false);
  };

  const toggleCategory = (cat: string) => {
    setFormCategories((prev) =>
      prev.includes(cat) ? prev.filter((c) => c !== cat) : [...prev, cat]
    );
  };

  const handleLikeToggle = async (postId: number, e: React.MouseEvent) => {
    e.stopPropagation();
    const res = await likePost(postId);
    if (res.success) {
      setPosts((prev) =>
        prev.map((p) =>
          p.id === postId
            ? {
                ...p,
                is_liked: res.data!.is_liked,
                like_count: res.data!.likes,
                dislike_count: res.data!.dislikes,
              }
            : p
        )
      );
    }
  };

  const handleDislikeToggle = async (postId: number, e: React.MouseEvent) => {
    e.stopPropagation();
    const res = await dislikePost(postId);
    if (res.success) {
      setPosts((prev) =>
        prev.map((p) =>
          p.id === postId
            ? {
                ...p,
                is_liked: res.data!.is_liked,
                like_count: res.data!.likes,
                dislike_count: res.data!.dislikes,
              }
            : p
        )
      );
    }
  };

  return (
    <main className="main">
      {/* Create Post Form Toggle */}
      <div className="card" style={{ marginBottom: showForm ? 0 : undefined }}>
        {!showForm ? (
          <button
            className="btn btn-p"
            style={{
              width: "100%",
              display: "flex",
              alignItems: "center",
              gap: "6px",
              fontSize: "12px",
              justifyContent: "center",
            }}
            onClick={() => setShowForm(true)}
          >
            <i className="ti ti-plus" style={{ fontSize: "14px" }} /> new post
          </button>
        ) : (
          <div>
            <p
              style={{
                fontSize: "12px",
                color: "#a09c94",
                marginBottom: "10px",
                fontWeight: 500,
              }}
            >
              Create a post
            </p>
            <input
              className="inp"
              placeholder="Title"
              value={formTitle}
              onChange={(e) => setFormTitle(e.target.value)}
              style={{ fontSize: "12px", marginBottom: "8px" }}
              maxLength={255}
            />
            <textarea
              className="inp"
              rows={4}
              placeholder="What's on your mind?"
              value={formText}
              onChange={(e) => setFormText(e.target.value)}
              style={{ resize: "none", fontSize: "12px", marginBottom: "8px" }}
              maxLength={2000}
            />
            {/* Image preview */}
            {formImagePreview && (
              <div
                style={{
                  marginBottom: 8,
                  position: "relative",
                  display: "inline-block",
                }}
              >
                <Image
                  src={formImagePreview}
                  alt="preview"
                  width={0}
                  height={0}
                  sizes="100px"
                  style={{
                    maxHeight: 100,
                    borderRadius: 4,
                    border: "0.5px solid #3a3733",
                    width: "auto",
                    height: "auto",
                  }}
                  unoptimized
                />
                <button
                  className="btn btn-red"
                  style={{
                    position: "absolute",
                    top: -6,
                    right: -6,
                    width: 20,
                    height: 20,
                    padding: 0,
                    display: "flex",
                    alignItems: "center",
                    justifyContent: "center",
                    fontSize: 10,
                  }}
                  onClick={() => {
                    setFormImage(null);
                    setFormImagePreview(null);
                    if (fileInputRef.current) fileInputRef.current.value = "";
                  }}
                >
                  x
                </button>
              </div>
            )}
            <input
              ref={fileInputRef}
              type="file"
              accept="image/*"
              style={{ display: "none" }}
              onChange={(e) => {
                const file = e.target.files?.[0];
                if (file) {
                  setFormImage(file);
                  const reader = new FileReader();
                  reader.onload = () =>
                    setFormImagePreview(reader.result as string);
                  reader.readAsDataURL(file);
                }
              }}
            />
            {/* Privacy selector */}
            <div
              style={{
                display: "flex",
                gap: "4px",
                marginBottom: "8px",
                flexWrap: "wrap",
              }}
            >
              {["public", "almost_private", "private"].map((p) => (
                <button
                  key={p}
                  className={`btn ${formPrivacy === p ? "btn-t" : "btn-g"}`}
                  style={{ fontSize: "10px", padding: "3px 8px" }}
                  onClick={() => {
                    setFormPrivacy(p);
                    if (p !== "private") {
                      setSelectedUsers([]);
                      setInviteQuery("");
                    }
                  }}
                >
                  {p.replace("_", " ")}
                </button>
              ))}
            </div>

            {/* Allowed users selection for private posts */}
            {formPrivacy === "private" && (
              <div style={{ marginBottom: "10px" }}>
                <p style={{ fontSize: "10px", color: "#a09c94", marginBottom: "6px", fontWeight: 500 }}>
                  only these users can see this post
                </p>

                {/* Selected users chips */}
                {selectedUsers.length > 0 && (
                  <div
                    style={{
                      display: "flex",
                      flexWrap: "wrap",
                      gap: 6,
                      marginBottom: 8,
                    }}
                  >
                    {selectedUsers.map((u) => (
                      <span
                        key={u.id}
                        className="tag tag-pink"
                        style={{ fontSize: 10, cursor: "default", display: "flex", alignItems: "center", gap: 4 }}
                      >
                        {u.nickname || `${u.firstname} ${u.lastname}`.trim()}
                        <span
                          style={{ cursor: "pointer", marginLeft: 2 }}
                          onClick={() => removeUser(u.id)}
                        >
                          ×
                        </span>
                      </span>
                    ))}
                  </div>
                )}

                <div style={{ position: "relative" }}>
                  <input
                    className="inp"
                    placeholder="search by name..."
                    value={inviteQuery}
                    onChange={(e) => setInviteQuery(e.target.value)}
                    onFocus={() => searchResults.length > 0 && setSearchOpen(true)}
                    style={{ fontSize: "11px" }}
                  />

                  {searchOpen && searchResults.length > 0 && (
                    <div
                      className="card"
                      style={{
                        position: "absolute",
                        top: "calc(100% + 2px)",
                        left: 0,
                        right: 0,
                        zIndex: 10,
                        padding: "4px 0",
                        maxHeight: 180,
                        overflowY: "auto",
                      }}
                    >
                      {searchResults.map((u) => (
                        <div
                          key={u.id}
                          onClick={() => addUser(u)}
                          style={{
                            display: "flex",
                            alignItems: "center",
                            gap: 8,
                            padding: "6px 10px",
                            cursor: "pointer",
                            fontSize: 11,
                          }}
                          onMouseEnter={(e) => {
                            (e.currentTarget as HTMLElement).style.background = "#2e2b27";
                          }}
                          onMouseLeave={(e) => {
                            (e.currentTarget as HTMLElement).style.background = "transparent";
                          }}
                        >
                          <div
                            className="av"
                            style={{
                              width: 22,
                              height: 22,
                              borderRadius: "50%",
                              background: u.avatar ? `url(${u.avatar}) center/cover` : "#2e1e24",
                              color: "#D4537E",
                              fontSize: 9,
                            }}
                          >
                            {!u.avatar && (u.nickname?.[0]?.toUpperCase() || u.firstname?.[0]?.toUpperCase() || "?")}
                          </div>
                          <span>{u.nickname || `${u.firstname} ${u.lastname}`.trim()}</span>
                          <span style={{ marginLeft: "auto", color: "#D4537E", fontSize: 10 }}>+ add</span>
                        </div>
                      ))}
                    </div>
                  )}
                  {searchOpen && searchResults.length === 0 && connections.length > 0 && (
                    <div
                      className="card"
                      style={{
                        position: "absolute",
                        top: "calc(100% + 2px)",
                        left: 0,
                        right: 0,
                        zIndex: 10,
                        padding: "6px 10px",
                        fontSize: 10,
                        color: "#6b6760",
                      }}
                    >
                      no matching connections found
                    </div>
                  )}
                </div>
              </div>
            )}
            {/* Categories */}
            <div
              style={{
                display: "flex",
                gap: "4px",
                marginBottom: "10px",
                flexWrap: "wrap",
              }}
            >
              {CATEGORIES.map((cat) => (
                <span
                  key={cat}
                  className={`tag ${
                    formCategories.includes(cat) ? "tag-pink" : "tag-gray"
                  }`}
                  style={{ cursor: "pointer", fontSize: "10px" }}
                  onClick={() => toggleCategory(cat)}
                >
                  {cat}
                </span>
              ))}
            </div>
            {/* Actions */}
            <div
              style={{
                display: "flex",
                justifyContent: "space-between",
              }}
            >
              <div style={{ display: "flex", gap: "4px" }}>
                <button
                  className="btn btn-g"
                  style={{
                    fontSize: "11px",
                    display: "flex",
                    alignItems: "center",
                    gap: 4,
                  }}
                  onClick={() => fileInputRef.current?.click()}
                >
                  <i className="ti ti-photo" style={{ fontSize: "12px" }} />{" "}
                  image
                </button>
              </div>
              <div style={{ display: "flex", gap: "4px" }}>
                <button
                  className="btn btn-g"
                  style={{ fontSize: "11px" }}
                  onClick={() => {
                    setShowForm(false);
                    setFormImage(null);
                    setFormImagePreview(null);
                  }}
                >
                  cancel
                </button>
                <button
                  className="btn btn-p"
                  style={{
                    fontSize: "11px",
                    display: "flex",
                    alignItems: "center",
                    gap: 4,
                    opacity:
                      !formTitle.trim() || !formText.trim() || submitting
                        ? 0.6
                        : 1,
                  }}
                  disabled={
                    !formTitle.trim() || !formText.trim() || submitting
                  }
                  onClick={() => void handleCreatePost()}
                >
                  <i className="ti ti-send" style={{ fontSize: "12px" }} />{" "}
                  {submitting ? "posting..." : "publish"}
                </button>
              </div>
            </div>
          </div>
        )}
      </div>

      {/* ── Category Filter Bar ── */}
      <div
        className="card"
        style={{
          display: "flex",
          gap: "4px",
          flexWrap: "wrap",
          padding: "8px 12px",
        }}
      >
        <span
          className={`tag ${filterCategories.length === 0 ? "tag-pink" : "tag-gray"}`}
          style={{ cursor: "pointer", fontSize: "10px" }}
          onClick={() => {
            setFilterCategories([]);
            void fetchPosts([]);
          }}
        >
          all
        </span>
        {CATEGORIES.map((cat) => (
          <span
            key={cat}
            className={`tag ${filterCategories.includes(cat) ? "tag-pink" : "tag-gray"}`}
            style={{ cursor: "pointer", fontSize: "10px" }}
            onClick={() => toggleFilterCategory(cat)}
          >
            {cat}
          </span>
        ))}
      </div>

      {loading ? (
        <div className="card" style={{ textAlign: "center", padding: "24px" }}>
          <p style={{ fontSize: "12px", color: "#a09c94" }}>
            Loading feed...
          </p>
        </div>
      ) : posts.length === 0 ? (
        <div className="card" style={{ textAlign: "center", padding: "24px" }}>
          <p style={{ fontSize: "14px", color: "#e8e4dc" }}>No posts yet</p>
          <p
            style={{
              fontSize: "11px",
              color: "#6b6760",
              marginTop: "4px",
            }}
          >
            Be the first to create a post!
          </p>
        </div>
      ) : (
        posts.map((post) => {
          const firstInitial =
            post.firstname?.charAt(0)?.toUpperCase() ?? "";
          const lastInitial = post.lastname?.charAt(0)?.toUpperCase() ?? "";
          const initials = firstInitial + lastInitial || post.nickname?.charAt(0)?.toUpperCase() || "?";

          return (
            <div
              key={post.id}
              className="card"
              style={{ cursor: "pointer" }}
              onClick={() => router.push(`/posts/${post.id}`)}
            >
              {/* Author header */}
              <div
                style={{
                  display: "flex",
                  alignItems: "center",
                  gap: "8px",
                  marginBottom: "8px",
                }}
              >
                <div
                  className="av"
                  style={{
                    width: "30px",
                    height: "30px",
                    background: "#FBEAF0",
                    color: "#993556",
                    fontSize: "11px",
                    overflow: "hidden",
                  }}
                >
                  {post.avatar ? (
                    <Image
                      src={post.avatar}
                      alt="avatar"
                      width={30}
                      height={30}
                      style={{ objectFit: "cover" }}
                    />
                  ) : (
                    initials
                  )}
                </div>
                <div style={{ flex: 1 }}>
                  <p
                    style={{
                      fontSize: "12px",
                      fontWeight: 500,
                      color: "#e8e4dc",
                    }}
                  >
                    {post.nickname ||
                      `${post.firstname || ""} ${post.lastname || ""}`.trim()}
                    {post.privacy && (
                      <span
                        className={`tag ${
                          post.privacy === "public"
                            ? "tag-teal"
                            : post.privacy === "almost_private"
                            ? "tag-gray"
                            : "tag-pink"
                        }`}
                        style={{ marginLeft: 6 }}
                      >
                        {post.privacy.replace("_", " ")}
                      </span>
                    )}
                  </p>
                  <p
                    style={{
                      fontSize: "10px",
                      color: "#6b6760",
                    }}
                  >
                    {post.time_ago ||
                      new Date(post.created_at).toLocaleDateString()}
                    {post.categories?.length > 0 && (
                      <> · {post.categories.join(", ")}</>
                    )}
                  </p>
                </div>
              </div>

              {/* Title */}
              {post.title && (
                <h3
                  style={{
                    fontSize: "15px",
                    fontWeight: 600,
                    marginBottom: "6px",
                    color: "#e8e4dc",
                  }}
                >
                  {post.title}
                </h3>
              )}

              {/* Text */}
              <p
                style={{
                  fontSize: "12px",
                  color: "#e8e4dc",
                  lineHeight: 1.6,
                  marginBottom: post.image ? "8px" : "0",
                }}
              >
                {post.text}
              </p>

              {/* Image */}
              {post.image && (
                <div
                  style={{
                    marginBottom: "10px",
                    borderRadius: "6px",
                    overflow: "hidden",
                    border: "0.5px solid #3a3733",
                    position: "relative",
                    height: "280px",
                    background: "#2a2824",
                  }}
                >
                  <Image
                    src={post.image}
                    alt="post image"
                    fill
                    sizes="100vw"
                    style={{ objectFit: "cover" }}
                    unoptimized
                  />
                </div>
              )}

              <div className="divider" />

              {/* Actions */}
              <div
                style={{
                  display: "flex",
                  gap: "6px",
                  alignItems: "center",
                }}
              >
                <button
                  className={`btn ${post.is_liked === 1 ? "btn-t" : "btn-g"}`}
                  style={{
                    display: "flex",
                    alignItems: "center",
                    gap: "3px",
                    fontSize: "11px",
                  }}
                  onClick={(e) => void handleLikeToggle(post.id, e)}
                >
                  <i
                    className="ti ti-thumb-up"
                    style={{
                      fontSize: "12px",
                      color: post.is_liked === 1 ? "#1D9E75" : undefined,
                    }}
                  />{" "}
                  {post.like_count || 0}
                </button>
                <button
                  className={`btn ${post.is_liked === -1 ? "btn-red" : "btn-g"}`}
                  style={{
                    display: "flex",
                    alignItems: "center",
                    gap: "3px",
                    fontSize: "11px",
                  }}
                  onClick={(e) => void handleDislikeToggle(post.id, e)}
                >
                  <i
                    className="ti ti-thumb-down"
                    style={{
                      fontSize: "12px",
                      color: post.is_liked === -1 ? "#D4537E" : undefined,
                    }}
                  />{" "}
                  {post.dislike_count || 0}
                </button>
                <button
                  className="btn btn-g"
                  style={{
                    display: "flex",
                    alignItems: "center",
                    gap: "3px",
                    fontSize: "11px",
                  }}
                  onClick={(e) => {
                    e.stopPropagation();
                    router.push(`/posts/${post.id}`);
                  }}
                >
                  <i
                    className="ti ti-message-circle"
                    style={{ fontSize: "12px" }}
                  />{" "}
                  {post.comment_count || 0}
                </button>
              </div>
            </div>
          );
        })
      )}

      {/* Infinite scroll sentinel */}
      <div ref={sentinelRef} />
      {loadingMore && (
        <div className="card" style={{ textAlign: "center", padding: "16px" }}>
          <p style={{ fontSize: "11px", color: "#a09c94" }}>Loading more posts...</p>
        </div>
      )}
      {!hasMore && posts.length > 0 && (
        <div className="card" style={{ textAlign: "center", padding: "12px" }}>
          <p style={{ fontSize: "10px", color: "#6b6760" }}>
            — you've reached the end —
          </p>
        </div>
      )}
    </main>
  );
}
