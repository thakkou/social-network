"use client";

import { useState, useEffect, useCallback, useRef } from "react";
import { useRouter } from "next/navigation";
import {
  getFeedPosts,
  likePost,
  dislikePost,
  createPost,
  type FeedPost,
} from "~/app/api/crud/post";

const CATEGORIES = [
  "General", "Lifestyle", "Health & Fitness", "Travel",
  "Food & Cooking", "Education", "Business", "Finance",
  "Entertainment", "Sports", "Personal Dev", "Culture", "News",
];

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

  const fetchPosts = useCallback(async () => {
    setLoading(true);
    const res = await getFeedPosts({ limit: 20 });
    if (res.success) {
      setPosts(res.data);
    }
    setLoading(false);
  }, []);

  useEffect(() => {
    void fetchPosts();
  }, [fetchPosts]);

  const handleCreatePost = async () => {
    if (!formTitle.trim() || !formText.trim()) return;
    setSubmitting(true);

    const fd = new FormData();
    fd.set("title", formTitle.trim());
    fd.set("text", formText.trim());
    fd.set("privacy", formPrivacy);
    formCategories.forEach((cat) => fd.append("categories", cat));
    if (formImage) fd.set("image", formImage);

    const res = await createPost(fd);
    if (res.success) {
      setShowForm(false);
      setFormTitle("");
      setFormText("");
      setFormPrivacy("public");
      setFormCategories(["General"]);
      setFormImage(null);
      setFormImagePreview(null);
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
                <img
                  src={formImagePreview}
                  alt="preview"
                  style={{
                    maxHeight: 100,
                    borderRadius: 4,
                    border: "0.5px solid #3a3733",
                  }}
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
                  onClick={() => setFormPrivacy(p)}
                >
                  {p.replace("_", " ")}
                </button>
              ))}
            </div>
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
                    <img
                      src={post.avatar}
                      alt="avatar"
                      style={{
                        width: "100%",
                        height: "100%",
                        objectFit: "cover",
                      }}
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
                  }}
                >
                  <img
                    src={post.image}
                    alt="post image"
                    style={{
                      width: "100%",
                      maxHeight: "300px",
                      objectFit: "cover",
                    }}
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
    </main>
  );
}
