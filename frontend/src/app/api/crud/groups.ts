"use server";

import { fetchApi } from "../helper/fetch";

interface CreateGroupResponse {
  status_code: number;
  message: string;
  data: {
    group_id: number;
    title: string;
  };
}

interface GroupApiResponse<T> {
  status_code: number;
  message: string;
  data: T;
}

// ─── Types for Group Public Details ───
export interface GroupPublic {
  id: number;
  creator_id: number;
  title: string;
  description: string;
  logo: string;
  background: string;
  created_at: string;
}

interface GroupPublicResponse {
  status_code: number;
  message: string;
  data: {
    group: GroupPublic;
    is_member: boolean;
  };
}

// ─── Types for Feed Items (Posts & Events) ───
export interface FeedAuthor {
  id: number;
  nickname: string;
  firstname: string;
  lastname: string;
  avatar: string;
}

export interface EventResponder {
  user_id: number;
  status: string;
  nickname: string;
  firstname: string;
  lastname: string;
  avatar: string;
}

export interface GroupFeedItem {
  id: number;
  type: "post" | "event";
  group_id: number;
  user_id: number;
  title: string;
  text: string;
  image: string;
  description: string;
  event_time: string;
  created_at: string;
  author: FeedAuthor;
  likes_count: number;
  dislikes_count: number;
  is_liked: boolean;
  comments_count: number;
  event_responses: EventResponder[];
}

// Create a new group using FormData (supports file uploads for logo and background)
export async function createGroup(formData: FormData) {
  const titleVal = formData.get("title");
  const descVal = formData.get("description");
  const title = typeof titleVal === "string" ? titleVal.trim() : "";
  const description = typeof descVal === "string" ? descVal.trim() : "";

  if (!title) {
    return { error: "Title is required." };
  }

  // Send FormData directly via fetchApi (Next.js automatically sets multipart headers)
  const result = await fetchApi<CreateGroupResponse>("/api/groups/create", {
    method: "POST",
    body: formData,
  });

  if (!result.success) {
    return { error: result.error };
  }

  return {
    success: true,
    message: result.data.message,
    group: result.data.data,
  };
}

// Fetch public group details by ID
export async function getGroupPublic(groupId: string) {
  if (!groupId) return { error: "Group ID is required." };

  const result = await fetchApi<GroupPublicResponse>(
    `/api/groups/public/${groupId}`,
    { method: "GET" }
  );

  if (!result.success) return { error: result.error };
  return {
    success: true,
    data: result.data.data.group,
    isMember: result.data.data.is_member,
  };
}

// Fetch group content feed (posts + events)
export async function getGroupContent(
  groupId: string,
  params?: { limit?: number; last_id?: number }
) {
  if (!groupId) return { error: "Group ID is required." };

  const result = await fetchApi<GroupApiResponse<GroupFeedItem[]>>(
    `/api/groups/content/${groupId}`,
    {
      method: "GET",
      searchParams: {
        limit: params?.limit ?? 20,
        last_id: params?.last_id ?? 0,
      },
    }
  );

  if (!result.success) return { error: result.error };
  return { success: true, data: result.data.data };
}

// ─── Mutations ───

// Create a group post (supports FormData with image)
export async function createGroupPost(
  groupId: string,
  payload: { title?: string; text?: string; image?: File }
) {
  if (!groupId) return { error: "Group ID is required." };

  let body: FormData | { title?: string; text?: string };

  if (payload.image) {
    const formData = new FormData();
    if (payload.title) formData.set("title", payload.title);
    if (payload.text) formData.set("text", payload.text);
    formData.set("image", payload.image);
    body = formData;
  } else {
    body = { title: payload.title, text: payload.text };
  }

  const result = await fetchApi<GroupApiResponse<{ post_id: number }>>(
    `/api/groups/${groupId}/posts`,
    { method: "POST", body }
  );

  if (!result.success) return { error: result.error };
  return { success: true, data: result.data.data };
}

// Like or dislike a group post (isLike: 1 = like, -1 = dislike)
export async function toggleGroupPostReaction(
  groupId: string,
  postId: number,
  isLike: number
) {
  if (!groupId) return { error: "Group ID is required." };

  const result = await fetchApi<GroupApiResponse<null>>(
    `/api/groups/${groupId}/posts/${postId}/reaction`,
    { method: "POST", body: { is_like: isLike } }
  );

  if (!result.success) return { error: result.error };
  return { success: true };
}

// Comment on a group post
export async function createGroupPostComment(
  groupId: string,
  postId: number,
  text: string
) {
  if (!groupId) return { error: "Group ID is required." };

  const result = await fetchApi<GroupApiResponse<null>>(
    `/api/groups/${groupId}/posts/${postId}/comments`,
    { method: "POST", body: { text } }
  );

  if (!result.success) return { error: result.error };
  return { success: true };
}

// Create a group event
export async function createGroupEvent(
  groupId: string,
  payload: { title: string; description?: string; event_time: string }
) {
  if (!groupId) return { error: "Group ID is required." };

  const result = await fetchApi<GroupApiResponse<{ event_id: number }>>(
    `/api/groups/${groupId}/events`,
    { method: "POST", body: payload }
  );

  if (!result.success) return { error: result.error };
  return { success: true, data: result.data.data };
}

// Request to join a group
export async function requestToJoinGroup(groupId: string) {
  if (!groupId) return { error: "Group ID is required." };

  const result = await fetchApi<GroupApiResponse<null>>(
    `/api/groups/${groupId}/join`,
    { method: "POST" }
  );

  if (!result.success) return { error: result.error };
  return { success: true };
}

// ─── Pending Requests ───
export interface PendingRequest {
  id: number;
  user_id: number;
  nickname: string;
  firstname: string;
  lastname: string;
  avatar: string;
  status: string;
  created_at: string;
}

// Get pending join requests (only group creator)
export async function getPendingRequests(groupId: string) {
  if (!groupId) return { error: "Group ID is required." };

  const result = await fetchApi<GroupApiResponse<PendingRequest[]>>(
    `/api/groups/${groupId}/requests`,
    { method: "GET" }
  );

  if (!result.success) return { error: result.error };
  return { success: true, data: result.data.data };
}

// Accept a join request
export async function acceptJoinRequest(groupId: string, userId: number) {
  const result = await fetchApi<GroupApiResponse<null>>(
    `/api/groups/${groupId}/requests/${userId}/accept`,
    { method: "POST" }
  );
  if (!result.success) return { error: result.error };
  return { success: true };
}

// Reject a join request
export async function rejectJoinRequest(groupId: string, userId: number) {
  const result = await fetchApi<GroupApiResponse<null>>(
    `/api/groups/${groupId}/requests/${userId}/reject`,
    { method: "POST" }
  );
  if (!result.success) return { error: result.error };
  return { success: true };
}

// ─── User's Groups ───
export async function getUserGroups() {
  const result = await fetchApi<GroupApiResponse<GroupPublic[]>>(
    `/api/users/groups`,
    { method: "GET" }
  );
  if (!result.success) return { error: result.error };
  return { success: true, data: result.data.data };
}

// ─── Privacy ───
export async function updateProfilePrivacy(isPrivate: boolean) {
  const result = await fetchApi<GroupApiResponse<null>>(
    `/api/profile/privacy`,
    { method: "PUT", body: { is_private: isPrivate ? 1 : 0 } }
  );
  if (!result.success) return { error: result.error };
  return { success: true };
}

// ─── Group Invite Actions ───

// Accept a group invite
export async function acceptGroupInvite(groupId: string) {
  const result = await fetchApi<GroupApiResponse<null>>(
    `/api/groups/${groupId}/invites/accept`,
    { method: "POST" }
  );
  if (!result.success) return { error: result.error };
  return { success: true };
}

// Reject a group invite
export async function rejectGroupInvite(groupId: string) {
  const result = await fetchApi<GroupApiResponse<null>>(
    `/api/groups/${groupId}/invites/reject`,
    { method: "POST" }
  );
  if (!result.success) return { error: result.error };
  return { success: true };
}

// ─── Profile Update (server action) ───

export async function updateProfileNickname(formData: FormData) {
  const result = await fetchApi<{ status_code: number; message: string }>(
    "/api/profile/update",
    { method: "PUT", body: formData }
  );

  if (!result.success) return { error: result.error };
  return { success: true };
}

// ─── Leave Group ───

export async function leaveGroup(groupId: string) {
  if (!groupId) return { error: "Group ID is required." };

  const result = await fetchApi<GroupApiResponse<null>>(
    `/api/groups/${groupId}/leave`,
    { method: "POST" }
  );

  if (!result.success) return { error: result.error };
  return { success: true };
}

// Respond to a group event (status: "going" | "not_going")
export async function respondToEvent(
  groupId: string,
  eventId: number,
  status: string
) {
  if (!groupId) return { error: "Group ID is required." };

  const result = await fetchApi<GroupApiResponse<null>>(
    `/api/groups/${groupId}/events/${eventId}/respond`,
    { method: "POST", body: { status } }
  );

  if (!result.success) return { error: result.error };
  return { success: true };
}