"use client";

import { useState, useEffect, useCallback } from "react";
import { useRouter } from "next/navigation";
import {
  getFeedPosts,
  likePost,
  dislikePost,
  type FeedPost,
} from "~/app/api/crud/post";

export default function Home() {
  const router = useRouter();
  const [posts, setPosts] = useState<FeedPost[]>([]);
  const [loading, setLoading] = useState(true);

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

  const handleLikeToggle = async (postId: number, e: React.MouseEvent) => {
    e.stopPropagation();
    const post = posts.find((p) => p.id === postId);
    if (!post) return;
    const res =
      post.is_liked === 1 ? await dislikePost(postId) : await likePost(postId);
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
    const post = posts.find((p) => p.id === postId);
    if (!post) return;
    const res =
      post.is_liked === -1
        ? await likePost(postId)
        : await dislikePost(postId);
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
