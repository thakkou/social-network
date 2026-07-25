"use client";

import { useState, useEffect, useRef } from "react";
import { useParams } from "next/navigation";
import Link from "next/link";
import Image from "next/image";
import { useSession } from "next-auth/react";
import {
  getPostById,
  likePost,
  dislikePost,
  createComment,
  likeComment,
  dislikeComment,
  deleteComment,
  type FeedPost,
  type PostComment,
} from "~/app/_services/crud/post";

function formatCommentTime(createdAt: string): string {
  // Always compute time-ago from the timestamp locally rather than trusting
  // the backend's time_ago, which can produce wildly wrong values (e.g.
  // "3558 months ago") when the timestamp parsing goes wrong.
  if (!createdAt) return "";
  const date = new Date(createdAt);
  if (isNaN(date.getTime())) return "";
  const diffMs = Date.now() - date.getTime();
  const diffMin = Math.floor(diffMs / 60000);
  if (diffMin < 1) return "just now";
  if (diffMin < 60) return `${diffMin}m ago`;
  const diffHrs = Math.floor(diffMin / 60);
  if (diffHrs < 24) return `${diffHrs}h ago`;
  const diffDays = Math.floor(diffHrs / 24);
  if (diffDays < 7) return `${diffDays}d ago`;
  return date.toLocaleDateString();
}

export default function PostDetailPage() {
  const params = useParams();
  const { data: session } = useSession();
  const currentUserId = session?.user?.id;

  const postId = Number(params?.id);

  const [post, setPost] = useState<FeedPost | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  // Comment form
  const [commentText, setCommentText] = useState("");
  const [commentImage, setCommentImage] = useState<File | null>(null);
  const [commentImagePreview, setCommentImagePreview] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const [deleteError, setDeleteError] = useState<string | null>(null);

  const fileInputRef = useRef<HTMLInputElement>(null);

  const fetchPost = async () => {
    if (!postId || isNaN(postId)) return;
    setLoading(true);
    const res = await getPostById(postId);
    if (res.success) {
      setPost(res.data);
      setError("");
    } else {
      setError(res.error || "Post not found");
    }
    setLoading(false);
  };

  useEffect(() => {
    void fetchPost();
  }, [postId]);

  const handleLikeToggle = async () => {
    if (!post) return;
    const res = await likePost(post.id);
    if (res.success) {
      setPost((prev) =>
        prev
          ? {
              ...prev,
              is_liked: res.data!.is_liked,
              like_count: res.data!.likes,
              dislike_count: res.data!.dislikes,
            }
          : prev
      );
    }
  };

  const handleDislikeToggle = async () => {
    if (!post) return;
    const res = await dislikePost(post.id);
    if (res.success) {
      setPost((prev) =>
        prev
          ? {
              ...prev,
              is_liked: res.data!.is_liked,
              like_count: res.data!.likes,
              dislike_count: res.data!.dislikes,
            }
          : prev
      );
    }
  };

  const handleCommentLikeToggle = async (comment: PostComment) => {
    const res = await likeComment(comment.id);
    if (res.success) {
      setPost((prev) => {
        if (!prev) return prev;
        return {
          ...prev,
          comments: (prev.comments || []).map((c) =>
            c.id === comment.id
              ? {
                  ...c,
                  is_liked: res.data!.is_liked,
                  like_count: res.data!.likes,
                  dislike_count: res.data!.dislikes,
                }
              : c
          ),
        };
      });
    }
  };

  const handleCommentDislikeToggle = async (comment: PostComment) => {
    const res = await dislikeComment(comment.id);
    if (res.success) {
      setPost((prev) => {
        if (!prev) return prev;
        return {
          ...prev,
          comments: (prev.comments || []).map((c) =>
            c.id === comment.id
              ? {
                  ...c,
                  is_liked: res.data!.is_liked,
                  like_count: res.data!.likes,
                  dislike_count: res.data!.dislikes,
                }
              : c
          ),
        };
      });
    }
  };

  const handleDeleteComment = async (commentId: number) => {
    setDeleteError(null);
    const res = await deleteComment(commentId);
    if (res.success) {
      setPost((prev) =>
        prev
          ? {
              ...prev,
              comments: (prev.comments || []).filter((c) => c.id !== commentId),
              comment_count: Math.max(0, prev.comment_count - 1),
            }
          : prev
      );
    } else {
      setDeleteError(res.error || "Could not delete comment");
    }
  };

  const handleSubmitComment = async () => {
    if (!commentText.trim()) return;
    setSubmitting(true);
    const res = await createComment(
      postId,
      commentText.trim(),
      commentImage ?? undefined
    );
    if (res.success) {
      // Add the new comment locally instead of refetching the entire post (which scrolls to top)
      const newComment: PostComment = {
        id: res.data!.id,
        user_id: Number(currentUserId ?? 0),
        nickname: res.data!.nickname || "",
        created_at: res.data!.createdAt,
        time_ago: "just now",
        text: commentText.trim(),
        image: res.data!.image || "",
        like_count: 0,
        dislike_count: 0,
        is_liked: 0,
      };
      setPost((prev) =>
        prev
          ? {
              ...prev,
              comments: [...(prev.comments || []), newComment],
              comment_count: prev.comment_count + 1,
            }
          : prev
      );
      setCommentText("");
      setCommentImage(null);
      setCommentImagePreview(null);
    }
    setSubmitting(false);
  };

  const handleImageSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      setCommentImage(file);
      const reader = new FileReader();
      reader.onload = () => setCommentImagePreview(reader.result as string);
      reader.readAsDataURL(file);
    }
  };

  if (loading) {
    return (
      <main className="main">
        <div className="card" style={{ textAlign: "center", padding: "24px" }}>
          <p style={{ fontSize: "12px", color: "#a09c94" }}>Loading post...</p>
        </div>
      </main>
    );
  }

  if (error || !post) {
    return (
      <main className="main">
        <Link
          href="/"
          style={{
            display: "flex",
            alignItems: "center",
            gap: 6,
            color: "#a09c94",
            fontSize: 12,
            textDecoration: "none",
            width: "fit-content",
            marginBottom: 12,
          }}
        >
          <i className="ti ti-arrow-left" /> back to feed
        </Link>
        <div className="card" style={{ textAlign: "center", padding: "24px" }}>
          <p style={{ fontSize: "14px", fontWeight: 500, color: "#e8e4dc" }}>
            {error || "Post not found"}
          </p>
          <p style={{ fontSize: "11px", color: "#6b6760", marginTop: 4 }}>
            This post may have been removed or you don't have permission to view it.
          </p>
        </div>
      </main>
    );
  }

  const firstInitial = post.firstname?.charAt(0)?.toUpperCase() ?? "";
  const lastInitial = post.lastname?.charAt(0)?.toUpperCase() ?? "";
  const initials = firstInitial + lastInitial || post.nickname?.charAt(0)?.toUpperCase() || "?";

  return (
    <main className="main">
      {/* Back link */}
      <Link
        href="/"
        style={{
          display: "flex",
          alignItems: "center",
          gap: 6,
          color: "#a09c94",
          fontSize: 12,
          textDecoration: "none",
          width: "fit-content",
          marginBottom: 12,
        }}
      >
        <i className="ti ti-arrow-left" /> back to feed
      </Link>

      {/* Post Card */}
      <div className="card">
        {/* Author header */}
        <div
          style={{
            display: "flex",
            alignItems: "center",
            gap: "10px",
            marginBottom: "12px",
          }}
        >
          <Link
            href={`/profile/${post.user_id}`}
            className="av"
            style={{
              width: "36px",
              height: "36px",
              borderRadius: "50%",
              background: "#2e1e24",
              color: "#D4537E",
              fontSize: "13px",
              textDecoration: "none",
              overflow: "hidden",
            }}
          >
            {post.avatar ? (
              <Image
                src={post.avatar}
                alt="avatar"
                width={36}
                height={36}
                style={{ objectFit: "cover", borderRadius: "50%" }}
              />
            ) : (
              initials
            )}
          </Link>
          <div>
            <Link
              href={`/profile/${post.user_id}`}
              style={{
                fontSize: "13px",
                fontWeight: 600,
                color: "#e8e4dc",
                textDecoration: "none",
              }}
            >
              {post.nickname ||
                `${post.firstname || ""} ${post.lastname || ""}`.trim()}
            </Link>
            <div
              style={{
                display: "flex",
                alignItems: "center",
                gap: 6,
                marginTop: 2,
              }}
            >
              <span style={{ fontSize: "10px", color: "#6b6760" }}>
                {post.time_ago || new Date(post.created_at).toLocaleDateString()}
              </span>
              {post.privacy && (
                <span
                  className={`tag ${
                    post.privacy === "public"
                      ? "tag-teal"
                      : post.privacy === "almost_private"
                      ? "tag-gray"
                      : "tag-pink"
                  }`}
                  style={{ fontSize: "9px" }}
                >
                  {post.privacy.replace("_", " ")}
                </span>
              )}
              {post.categories?.length > 0 && (
                <span style={{ fontSize: "10px", color: "#6b6760" }}>
                  · {post.categories.join(", ")}
                </span>
              )}
            </div>
          </div>
        </div>

        {/* Title */}
        {post.title && (
          <h2
            style={{
              fontSize: "16px",
              fontWeight: 600,
              marginBottom: "8px",
              color: "#e8e4dc",
              lineHeight: 1.4,
            }}
          >
            {post.title}
          </h2>
        )}

        {/* Text */}
        <p
          style={{
            fontSize: "13px",
            color: "#e8e4dc",
            lineHeight: 1.7,
            marginBottom: post.image ? "12px" : 0,
          }}
        >
          {post.text}
        </p>

        {/* Image */}
        {post.image && (
          <div
            style={{
              marginBottom: "12px",
              borderRadius: "6px",
              overflow: "hidden",
              border: "0.5px solid #3a3733",
              position: "relative",
              height: "360px",
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

        {/* Reactions */}
        <div className="divider" />
        <div style={{ display: "flex", gap: "8px", alignItems: "center" }}>
          <button
            className={`btn ${post.is_liked === 1 ? "btn-t" : "btn-g"}`}
            style={{
              display: "flex",
              alignItems: "center",
              gap: "4px",
              fontSize: "11px",
            }}
            onClick={() => void handleLikeToggle()}
          >
            <i
              className="ti ti-thumb-up"
              style={{
                fontSize: "13px",
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
              gap: "4px",
              fontSize: "11px",
            }}
            onClick={() => void handleDislikeToggle()}
          >
            <i
              className="ti ti-thumb-down"
              style={{
                fontSize: "13px",
                color: post.is_liked === -1 ? "#D4537E" : undefined,
              }}
            />{" "}
            {post.dislike_count || 0}
          </button>
          <span
            style={{
              display: "flex",
              alignItems: "center",
              gap: "4px",
              fontSize: "11px",
              color: "#6b6760",
            }}
          >
            <i className="ti ti-message-circle" style={{ fontSize: "13px" }} />{" "}
            {post.comment_count || 0} {post.comment_count === 1 ? "comment" : "comments"}
          </span>
        </div>
      </div>

      {/* Comment Form */}
      <div className="card">
        <p
          style={{
            fontSize: "12px",
            color: "#a09c94",
            marginBottom: "10px",
            fontWeight: 500,
          }}
        >
          Write a comment
        </p>
        <textarea
          className="inp"
          rows={3}
          placeholder="Share your thoughts..."
          value={commentText}
          onChange={(e) => setCommentText(e.target.value)}
          style={{ resize: "none", fontSize: "12px" }}
          maxLength={1000}
        />
        {commentImagePreview && (
          <div
            style={{
              marginTop: 8,
              position: "relative",
              display: "inline-block",
            }}
          >
            <Image
              src={commentImagePreview}
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
                setCommentImage(null);
                setCommentImagePreview(null);
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
          onChange={handleImageSelect}
        />
        <div
          style={{
            display: "flex",
            justifyContent: "space-between",
            marginTop: "10px",
          }}
        >
          <button
            className="btn btn-g"
            style={{ fontSize: "11px", display: "flex", alignItems: "center", gap: 4 }}
            onClick={() => fileInputRef.current?.click()}
          >
            <i className="ti ti-photo" style={{ fontSize: "12px" }} /> image
          </button>
          <button
            className="btn btn-p"
            style={{
              fontSize: "11px",
              display: "flex",
              alignItems: "center",
              gap: 4,
              opacity: !commentText.trim() || submitting ? 0.6 : 1,
            }}
            disabled={!commentText.trim() || submitting}
            onClick={() => void handleSubmitComment()}
          >
            <i className="ti ti-send" style={{ fontSize: "12px" }} />{" "}
            {submitting ? "posting..." : "comment"}
          </button>
        </div>
      </div>

      {/* Delete error message */}
      {deleteError && (
        <div
          className="card"
          style={{
            textAlign: "center",
            padding: "10px",
            border: "0.5px solid #7a2c2c",
            background: "#2a1818",
          }}
        >
          <p style={{ fontSize: "11px", color: "#e07070" }}>
            <i className="ti ti-alert-circle" style={{ marginRight: 4 }} />
            {deleteError}
          </p>
        </div>
      )}

      {/* Comments */}
      <div style={{ display: "flex", flexDirection: "column", gap: "8px" }}>
        {post?.comments?.length === 0 ? (
          <div
            className="card"
            style={{ textAlign: "center", padding: "16px" }}
          >
            <p style={{ fontSize: "12px", color: "#6b6760" }}>
              No comments yet. Be the first!
            </p>
          </div>
        ) : (
          post?.comments?.map((comment) => (
            <div key={comment.id} className="card">
              <div
                style={{
                  display: "flex",
                  alignItems: "flex-start",
                  gap: "10px",
                }}
              >
                <Link
                  href={`/profile/${comment.user_id}`}
                  className="av"
                  style={{
                    width: "28px",
                    height: "28px",
                    borderRadius: "50%",
                    background: "#E1F5EE",
                    color: "#0F6E56",
                    fontSize: "11px",
                    textDecoration: "none",
                    flexShrink: 0,
                  }}
                >
                  {(comment.nickname || "?")[0]?.toUpperCase()}
                </Link>
                <div style={{ flex: 1, minWidth: 0 }}>
                  <div
                    style={{
                      background: "#2e2b27",
                      border: "0.5px solid #3a3733",
                      padding: "8px 10px",
                      fontSize: "12px",
                      color: "#e8e4dc",
                      lineHeight: 1.6,
                    }}
                  >
                    <div
                      style={{
                        display: "flex",
                        justifyContent: "space-between",
                        alignItems: "center",
                        marginBottom: 4,
                      }}
                    >
                      <strong style={{ fontSize: "11px", color: "#D4537E" }}>
                        {comment.nickname || "user"}
                      </strong>
                      <div
                        style={{
                          display: "flex",
                          alignItems: "center",
                          gap: 4,
                        }}
                      >
                        <span
                          style={{ fontSize: "10px", color: "#6b6760" }}
                        >
                          {formatCommentTime(comment.created_at)}
                        </span>
                        {/* Delete button (own comments only) */}
                        {currentUserId &&
                          Number(currentUserId) === comment.user_id && (
                            <button
                              className="btn btn-red"
                              style={{
                                padding: "1px 4px",
                                fontSize: "9px",
                                lineHeight: 1,
                              }}
                              onClick={() => void handleDeleteComment(comment.id)}
                            >
                              <i className="ti ti-trash" />
                            </button>
                          )}
                      </div>
                    </div>
                    <p>{comment.text}</p>
                    {comment.image && (
                      <Image
                        src={comment.image}
                        alt="comment image"
                        width={0}
                        height={0}
                        sizes="120px"
                        style={{
                          maxHeight: 120,
                          borderRadius: 4,
                          marginTop: 6,
                          border: "0.5px solid #3a3733",
                          width: "auto",
                          height: "auto",
                        }}
                        unoptimized
                      />
                    )}
                  </div>

                  {/* Comment reactions */}
                  <div
                    style={{
                      display: "flex",
                      gap: "4px",
                      marginTop: "4px",
                    }}
                  >
                    <button
                      className={`btn-sm ${
                        comment.is_liked === 1 ? "btn-sm-t" : "btn-sm-g"
                      }`}
                      style={{
                        display: "flex",
                        alignItems: "center",
                        gap: "3px",
                        fontSize: "10px",
                        padding: "2px 6px",
                        border: "0.5px solid #3a3733",
                        background:
                          comment.is_liked === 1 ? "#1a3a2e" : "transparent",
                        color:
                          comment.is_liked === 1 ? "#1D9E75" : "#a09c94",
                        cursor: "pointer",
                        borderRadius: "3px",
                      }}
                      onClick={() => void handleCommentLikeToggle(comment)}
                    >
                      <i
                        className="ti ti-thumb-up"
                        style={{
                          fontSize: "11px",
                          color: comment.is_liked === 1 ? "#1D9E75" : undefined,
                        }}
                      />{" "}
                      {comment.like_count || 0}
                    </button>
                    <button
                      className={`btn-sm ${
                        comment.is_liked === -1 ? "btn-sm-red" : "btn-sm-g"
                      }`}
                      style={{
                        display: "flex",
                        alignItems: "center",
                        gap: "3px",
                        fontSize: "10px",
                        padding: "2px 6px",
                        border: "0.5px solid #3a3733",
                        background:
                          comment.is_liked === -1 ? "#3a1e24" : "transparent",
                        color:
                          comment.is_liked === -1 ? "#D4537E" : "#a09c94",
                        cursor: "pointer",
                        borderRadius: "3px",
                      }}
                      onClick={() => void handleCommentDislikeToggle(comment)}
                    >
                      <i
                        className="ti ti-thumb-down"
                        style={{
                          fontSize: "11px",
                          color: comment.is_liked === -1 ? "#D4537E" : undefined,
                        }}
                      />{" "}
                      {comment.dislike_count || 0}
                    </button>
                  </div>
                </div>
              </div>
            </div>
          ))
        )}
      </div>
    </main>
  );
}
