'use client';

import React from "react";
import Image from "next/image";

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
  is_liked: number;
  categories: string[] | null;
  comments: unknown[] | null;
}

interface ProfilePostsProps {
  posts: Post[];
}

export default function ProfilePosts({ posts }: ProfilePostsProps) {

  if (!posts || posts.length === 0) {
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

  return (
    <div
      id="profile-posts"
      style={{
        display: "flex",
        flexDirection: "column",
        gap: "8px",
      }}
    >
      {posts.map((post) => (
        <div key={post.id} className="card">

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
            </p>

            <span className={`tag ${getPrivacyTagClass(post.privacy)}`}>
              {post.privacy.replace("_", " ")}
            </span>
          </div>

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

          <p
            style={{
              fontSize: "12px",
              color: "#e8e4dc",
              lineHeight: 1.5,
              whiteSpace: "pre-wrap",
            }}
          >
            {post.text}
          </p>

          {post.image && (
            <div
              style={{
                marginTop: "8px",
                border: "0.5px solid #3a3733",
                overflow: "hidden",
                background: "#1a1917",
                position: "relative",
                width: "100%",
                height: "280px",
              }}
            >
              <Image
                src={post.image}
                alt={post.title || "Post media attachment"}
                fill
                style={{
                  objectFit: "cover",
                }}
              />
            </div>
          )}

          <div className="divider" style={{ margin: "8px 0 6px 0" }} />

          <div
            style={{
              display: "flex",
              alignItems: "center",
              gap: "16px",
            }}
          >
            <span
              style={{
                fontSize: "11px",
                color: post.is_liked === 1 ? "#D4537E" : "#a09c94",
                display: "flex",
                alignItems: "center",
                gap: "4px",
                cursor: "pointer",
              }}
            >
              <i
                className={`ti ${
                  post.is_liked === 1
                    ? "ti-heart-filled"
                    : "ti-heart"
                }`}
              />
              {post.like_count}
            </span>

            <span
              style={{
                fontSize: "11px",
                color: "#a09c94",
                display: "flex",
                alignItems: "center",
                gap: "4px",
              }}
            >
              <i className="ti ti-thumb-down" />
              {post.dislike_count}
            </span>
          </div>

        </div>
      ))}
    </div>
  );
}