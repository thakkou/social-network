"use client";

import { useState, useEffect } from "react";
import { useSession } from "next-auth/react";
import { fetchApi } from "~/app/api/helper/fetch";

interface PostUser {
  id: number;
  firstname: string;
  lastname: string;
  nickname: string;
  avatar: string;
}

interface PostComment {
  id: number;
  userId: number;
  nickname: string;
  text: string;
  createdAt: string;
  timeAgo: string;
}

interface FeedPost {
  id: number;
  userId: number;
  firstname: string;
  lastname: string;
  nickname: string;
  avatar: string;
  created_at: string;
  timeAgo: string;
  title: string;
  text: string;
  image: string;
  privacy: string;
  likeCount: number;
  dislikeCount: number;
  isLiked: number;
  commentCount: number;
  comments: PostComment[];
  categories: string[];
}

export default function Home() {
  const { data: session } = useSession();
  const [posts, setPosts] = useState<FeedPost[]>([]);
  const [loading, setLoading] = useState(true);

  const fetchPosts = async () => {
    setLoading(true);
    try {
      const res = await fetchApi<FeedPost[]>("/api/posts?limit=20", {
        method: "GET",
      });
      if (res.success) {
        const data = Array.isArray(res.data) ? res.data : (res.data as any).data ?? [];
        setPosts(data);
      }
    } catch (err) {
      console.error("Failed to fetch posts:", err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    void fetchPosts();
  }, []);

  return (
    <main className="main">
      {loading ? (
        <div className="card" style={{ textAlign: "center", padding: "24px" }}>
          <p style={{ fontSize: "12px", color: "var(--color-text-secondary)" }}>
            Loading feed...
          </p>
        </div>
      ) : posts.length === 0 ? (
        <div className="card" style={{ textAlign: "center", padding: "24px" }}>
          <p style={{ fontSize: "14px", color: "var(--color-text-primary)" }}>
            No posts yet
          </p>
          <p style={{ fontSize: "11px", color: "var(--color-text-secondary)", marginTop: "4px" }}>
            Be the first to create a post!
          </p>
        </div>
      ) : (
        posts.map((post) => (
          <div key={post.id} className="card">
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
                  (post.nickname || post.firstname || "?")[0]?.toUpperCase()
                )}
              </div>
              <div style={{ flex: 1 }}>
                <p
                  style={{
                    fontSize: "12px",
                    fontWeight: 500,
                    color: "var(--color-text-primary)",
                  }}
                >
                  {post.nickname ||
                    `${post.firstname} ${post.lastname}`}
                  <span className={`tag ${post.privacy === "public" ? "tag-teal" : post.privacy === "almost_private" ? "tag-gray" : "tag-pink"}`}>
                    {post.privacy?.replace("_", " ") || "public"}
                  </span>
                </p>
                <p style={{ fontSize: "10px", color: "var(--color-text-tertiary)" }}>
                  {post.timeAgo || new Date(post.created_at).toLocaleDateString()}
                  {post.categories?.length > 0 && (
                    <> · {post.categories.join(", ")}</>
                  )}
                </p>
              </div>
            </div>

            {post.title && (
              <h3
                style={{
                  fontSize: "15px",
                  fontWeight: 600,
                  marginBottom: "6px",
                  color: "var(--color-text-primary)",
                }}
              >
                {post.title}
              </h3>
            )}

            <p
              style={{
                fontSize: "12px",
                color: "var(--color-text-primary)",
                lineHeight: 1.6,
                marginBottom: post.image ? "8px" : "0",
              }}
            >
              {post.text}
            </p>

            {post.image && (
              <div
                style={{
                  marginBottom: "10px",
                  borderRadius: "6px",
                  overflow: "hidden",
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

            <div
              style={{
                display: "flex",
                gap: "6px",
                alignItems: "center",
              }}
            >
              <button
                className={`btn ${post.isLiked === 1 ? "btn-t" : "btn-g"}`}
                style={{
                  display: "flex",
                  alignItems: "center",
                  gap: "3px",
                }}
              >
                <i
                  className="ti ti-thumb-up"
                  style={{ fontSize: "12px" }}
                />{" "}
                {post.likeCount || 0}
              </button>
              <button
                className="btn btn-g"
                style={{
                  display: "flex",
                  alignItems: "center",
                  gap: "3px",
                }}
              >
                <i
                  className="ti ti-message-circle"
                  style={{ fontSize: "12px" }}
                />{" "}
                {post.commentCount || 0}
              </button>
            </div>

            {post.comments && post.comments.length > 0 && (
              <div
                style={{
                  marginTop: "8px",
                  paddingTop: "8px",
                  borderTop: "0.5px solid var(--color-border-tertiary)",
                }}
              >
                {post.comments.slice(0, 3).map((c) => (
                  <div
                    key={c.id}
                    style={{
                      display: "flex",
                      gap: "8px",
                      alignItems: "flex-start",
                      marginBottom: "6px",
                    }}
                  >
                    <div
                      className="av"
                      style={{
                        width: "22px",
                        height: "22px",
                        background: "#E1F5EE",
                        color: "#0F6E56",
                        fontSize: "10px",
                      }}
                    >
                      {(c.nickname || "?")[0]?.toUpperCase()}
                    </div>
                    <div
                      style={{
                        background: "var(--color-background-secondary)",
                        border: "0.5px solid var(--color-border-tertiary)",
                        padding: "5px 8px",
                        fontSize: "11px",
                        flex: 1,
                        color: "var(--color-text-primary)",
                      }}
                    >
                      <strong style={{ fontSize: "10px" }}>
                        {c.nickname}
                      </strong>
                      <br />
                      {c.text}
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        ))
      )}
    </main>
  );
}
