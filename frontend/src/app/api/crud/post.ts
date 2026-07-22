"use server";

import { fetchApi } from "../helper/fetch";

// ─── Types ───

export interface PostComment {
  id: number;
  user_id: number;
  nickname: string;
  created_at: string;
  time_ago: string;
  text: string;
  image: string;
  like_count: number;
  dislike_count: number;
  is_liked: number; // 1=liked, 0=none, -1=disliked
}

export interface FeedPost {
  id: number;
  user_id: number;
  nickname: string;
  firstname: string;
  lastname: string;
  avatar: string;
  created_at: string;
  time_ago: string;
  title: string;
  text: string;
  image: string;
  privacy: string;
  like_count: number;
  dislike_count: number;
  comment_count: number;
  is_liked: number;
  comments: PostComment[];
  categories: string[];
}

interface GoApiResponse<T> {
  status_code: number;
  message: string;
  data: T;
}

// ─── Posts ───

// Fetch feed posts
export async function getFeedPosts(params?: {
  limit?: number;
  last_id?: number;
  categories?: string[];
  liked_by_me?: boolean;
  posted_by_me?: boolean;
}) {
  const searchParams: Record<string, string | number | boolean | undefined> = {
    limit: params?.limit ?? 20,
    last_id: params?.last_id ?? 0,
    ...(params?.liked_by_me ? { liked_by_me: true } : {}),
    ...(params?.posted_by_me ? { posted_by_me: true } : {}),
  };

  // Pass categories as individual query params (backend reads r.Form["categories"])
  if (params?.categories && params.categories.length > 0) {
    // We'll pass them as a comma-separated string since fetchApi uses URLSearchParams
    searchParams.categories = params.categories.join(",");
  }

  const result = await fetchApi<GoApiResponse<FeedPost[]>>("/api/posts", {
    method: "GET",
    searchParams,
  });

  if (!result.success) return { error: result.error };
  return { success: true, data: result.data.data };
}

// Fetch a single post with full comments
export async function getPostById(postId: number) {
  const result = await fetchApi<GoApiResponse<FeedPost>>(
    `/api/posts/${postId}`,
    { method: "GET" }
  );

  if (!result.success) return { error: result.error };
  return { success: true, data: result.data.data };
}

// Like a post
export async function likePost(postId: number) {
  const result = await fetchApi<GoApiResponse<{
    post_id: number;
    likes: number;
    dislikes: number;
    is_liked: number;
  }>>(`/api/posts/${postId}/like`, { method: "POST" });

  if (!result.success) return { error: result.error };
  return { success: true, data: result.data.data };
}

// Dislike a post
export async function dislikePost(postId: number) {
  const result = await fetchApi<GoApiResponse<{
    post_id: number;
    likes: number;
    dislikes: number;
    is_liked: number;
  }>>(`/api/posts/${postId}/dislike`, { method: "POST" });

  if (!result.success) return { error: result.error };
  return { success: true, data: result.data.data };
}

// ─── Comments ───

// Create a comment (supports optional image)
export async function createComment(
  postId: number,
  text: string,
  image?: File
) {
  let body: FormData | { postId: number; text: string };

  if (image) {
    const fd = new FormData();
    fd.set("postId", String(postId));
    fd.set("text", text);
    fd.set("image", image);
    body = fd;
  } else {
    body = { postId, text };
  }

  const result = await fetchApi<GoApiResponse<{
    id: number;
    text: string;
    postId: number;
    userId: number;
    createdAt: string;
    nickname: string;
  }>>("/api/comments/create", { method: "POST", body });

  if (!result.success) return { error: result.error };
  return { success: true, data: result.data.data };
}

// Like a comment
export async function likeComment(commentId: number) {
  const result = await fetchApi<GoApiResponse<{
    comment_id: number;
    likes: number;
    dislikes: number;
    is_liked: number;
  }>>(`/api/comments/${commentId}/like`, { method: "POST" });

  if (!result.success) return { error: result.error };
  return { success: true, data: result.data.data };
}

// Dislike a comment
export async function dislikeComment(commentId: number) {
  const result = await fetchApi<GoApiResponse<{
    comment_id: number;
    likes: number;
    dislikes: number;
    is_liked: number;
  }>>(`/api/comments/${commentId}/dislike`, { method: "POST" });

  if (!result.success) return { error: result.error };
  return { success: true, data: result.data.data };
}

// Delete a comment
export async function deleteComment(commentId: number) {
  const result = await fetchApi<GoApiResponse<null>>(
    `/api/comments/${commentId}/delete`,
    { method: "DELETE" }
  );

  if (!result.success) return { error: result.error };
  return { success: true };
}

// ─── Create Post (for feed) ───

interface CreatePostResponse {
  status_code: number;
  message: string;
  data: {
    post_id: number;
  };
}

// Create a new post with optional image, privacy, categories
export async function createPost(formData: FormData) {
  const result = await fetchApi<CreatePostResponse>("/api/posts/create", {
    method: "POST",
    body: formData,
  });

  if (!result.success) return { error: result.error };
  return { success: true, data: result.data.data };
}
