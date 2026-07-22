'use client';

import React, { useState } from "react";
import { useRouter } from "next/navigation";
import { useSession } from "next-auth/react";
import { likePost, dislikePost, deletePost } from "~/app/api/crud/post";
import ConfirmModal from "~/app/_components/ConfirmModal";

interface Post {
  id: number;
  user_id: number;
  nickname: string;
  created_at: string;
  time_ago: string;
  title: string;
  text: string;
  image: string;
  privacy: "public" | "private" | "almost_private";
  like_count: number;
  dislike_count: number;
  comment_count?: number;
  is_liked: number;
  categories: string[] | null;
  comments: unknown[] | null;
}

interface ProfilePostsProps {
  posts: Post[];
}

export default function ProfilePosts({ posts }: ProfilePostsProps) {
  const router = useRouter();
  const { data: session } = useSession();
  const currentUserId = Number(session?.user?.id ?? 0);

  // Local state for optimistic like/dislike updates
  const [localPosts, setLocalPosts] = useState<Post[]>(posts);

  // Sync local state when prop changes
  React.useEffect(() => {
    setLocalPosts(posts);
  }, [posts]);

  if (!localPosts || localPosts.length === 0) {
    return (
      <div
        id="profile-posts"
        style={{
          textAlign: "center",
          padding: "2rem",
          color: "#6b6760",
          border: "0.5px dashed #3a3733",
        }}
      >
        <i
          className="ti ti-notes"
          style={{ fontSize: "20px", display: "block", marginBottom: "4px" }}
        />
        No posts available on this profile yet.
      </div>
    );
  }

  const getPrivacyTagClass = (privacy: Post["privacy"]) => {
    switch (privacy) {
      case "public":
        return "tag-teal";
      case "almost_private":
        return "tag-purple";
      case "private":
        return "tag-gray";
      default:
        return "tag-gray";
    }
  };

  const handleLikeToggle = async (postId: number, e: React.MouseEvent) => {
    e.stopPropagation();
    const res = await likePost(postId);
    if (res.success) {
      setLocalPosts((prev) =>
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
      setLocalPosts((prev) =>
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

  const [deleteTarget, setDeleteTarget] = useState<Post | null>(null);

  const handleDeletePost = async (postId: number) => {
    setDeleteTarget(null);
    const res = await deletePost(postId);
    if (res.success) {
      setLocalPosts((prev) => prev.filter((p) => p.id !== postId));
    }
  };

  const commentCount = (post: Post): number => {
    // Use comment_count if available (from profile API), fall back to comments array length
    if (typeof post.comment_count === "number") {
      return post.comment_count;
    }
    if (post.comments && Array.isArray(post.comments)) {
      return post.comments.length;
    }
    return 0;
  };

  return (
    <div
      id="profile-posts"
      style={{
        display: "flex",
        flexDirection: "column",
        gap: "8px",
      }}
    >
      {localPosts.map((post) => (
        <div
          key={post.id}
          className="card"
          style={{ cursor: "pointer" }}
          onClick={() => router.push(`/posts/${post.id}`)}
        >
          {/* Header row: timestamp + privacy tag */}
          <div
            style={{
              display: "flex",
              alignItems: "center",
              justifyContent: "space-between",
              marginBottom: "6px",
            }}
          >
            <p style={{ fontSize: "10px", color: "#6b6760" }}>
              {post.time_ago ||
                new Date(post.created_at).toLocaleDateString()}
              {post.categories && post.categories.length > 0 && (
                <>
                  {" · "}
                  {post.categories.join(", ")}
                </>
              )}
            </p>
            <span className={`tag ${getPrivacyTagClass(post.privacy)}`}>
              {post.privacy.replace("_", " ")}
            </span>
          </div>

          {/* Title */}
          {post.title && (
            <h4
              style={{
                fontSize: "13px",
                fontWeight: 600,
                color: "#e8e4dc",
                marginBottom: "4px",
              }}
            >
              {post.title}
            </h4>
          )}

          {/* Text */}
          <p
            style={{
              fontSize: "12px",
              color: "#e8e4dc",
              lineHeight: 1.5,
              whiteSpace: "pre-wrap",
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
                alt={post.title || "Post media"}
                style={{
                  width: "100%",
                  maxHeight: "300px",
                  objectFit: "cover",
                }}
              />
            </div>
          )}

          <div className="divider" />

          {/* Action buttons */}
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
              {commentCount(post)}
            </button>
            {/* Delete button — only for own posts */}
            {currentUserId === post.user_id && (
              <button
                className="btn btn-red"
                style={{
                  marginLeft: "auto",
                  display: "flex",
                  alignItems: "center",
                  gap: "3px",
                  fontSize: "11px",
                }}
                onClick={(e) => {
                  e.stopPropagation();
                  setDeleteTarget(post);
                }}
              >
                <i className="ti ti-trash" style={{ fontSize: "12px" }} />
              </button>
            )}
          </div>
        </div>
      ))}

      <ConfirmModal
        open={deleteTarget !== null}
        title="Delete post?"
        message={`Are you sure you want to delete this post${deleteTarget?.title ? `: "${deleteTarget.title}"` : ""}? This cannot be undone.`}
        confirmLabel="delete"
        confirmClass="btn-red"
        onConfirm={() => deleteTarget && void handleDeletePost(deleteTarget.id)}
        onCancel={() => setDeleteTarget(null)}
      />
    </div>
  );
}