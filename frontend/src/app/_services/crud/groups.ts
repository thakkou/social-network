"use server";

import { fetchApi } from "../helper/fetch";

type CreateGroupResponse = {
  status_code: number;
  message: string;
  data: {
    group_id: number;
    title: string;
  };
}

type GroupApiResponse<T> = {
  status_code: number;
  message: string;
  data: T;
}

// ─── Types for Group Public Details ───
export type GroupPublic = {
  id: number;
  creator_id: number;
  title: string;
  description: string;
  logo: string;
  background: string;
  created_at: string;
}

type GroupPublicResponse = {
  status_code: number;
  message: string;
  data: {
    group: GroupPublic;
    is_member: boolean;
  };
}

// ─── Types for Feed Items (Posts & Events) ───
export type FeedAuthor = {
  id: number;
  nickname: string;
  firstname: string;
  lastname: string;
  avatar: string;
}

export type EventResponder = {
  user_id: number;
  status: string;
  nickname: string;
  firstname: string;
  lastname: string;
  avatar: string;
}

export type GroupFeedComment = {
  id: number;
  user_id: number;
  text: string;
  created_at: string;
  nickname: string;
  firstname: string;
  lastname: string;
  avatar: string;
  likes_count: number;
  dislikes_count: number;
  is_liked: number;
}

export type GroupFeedItem = {
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
  is_liked: number; // 1 = liked, -1 = disliked, 0 = no reaction
  comments_count: number;
  comments: GroupFeedComment[];
  event_responses: EventResponder[];
}

// Create a new group using FormData (supports file uploads for logo and background)
// If FormData contains "invite_ids" (JSON array of user IDs), invites are sent after group creation.
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

  const groupData = result.data.data;

  // If invite IDs were included in the FormData, send invites from the server
  const inviteIdsRaw = formData.get("invite_ids");
  const invited: number[] = [];
  const failedInvites: string[] = [];

  if (inviteIdsRaw && typeof inviteIdsRaw === "string") {
    try {
      const inviteIds: number[] = JSON.parse(inviteIdsRaw);
      for (const userId of inviteIds) {
        const inviteRes = await fetchApi<GroupApiResponse<null>>(
          `/api/groups/${groupData.group_id}/invite`,
          { method: "POST", body: { user_id: userId } }
        );
        if (!inviteRes.success) {
          failedInvites.push(`user ${userId}: ${inviteRes.error}`);
        } else {
          invited.push(userId);
        }
        // Small delay to avoid hitting the backend rate limiter (500ms on /api/groups/)
        await new Promise((r) => setTimeout(r, 600));
      }
    } catch (e) {
      console.error("[CREATE_GROUP] Failed to parse invite_ids:", e);
    }
  }

  return {
    success: true,
    message: result.data.message,
    group: groupData,
    invitesSent: invited,
    inviteErrors: failedInvites,
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
  payload: { title: string; description?: string; event_time: string; image?: File }
) {
  if (!groupId) return { error: "Group ID is required." };

  let body: FormData | { title: string; description?: string; event_time: string };

  if (payload.image) {
    const formData = new FormData();
    formData.set("title", payload.title);
    if (payload.description) formData.set("description", payload.description);
    formData.set("event_time", payload.event_time);
    formData.set("image", payload.image);
    body = formData;
  } else {
    body = { title: payload.title, description: payload.description, event_time: payload.event_time };
  }

  const result = await fetchApi<GroupApiResponse<{ event_id: number }>>(
    `/api/groups/${groupId}/events`,
    { method: "POST", body }
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
export type PendingRequest = {
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

// ─── Group Post Comment Reaction ───

export async function toggleGroupPostCommentReaction(
  groupId: string,
  postId: number,
  commentId: number,
  isLike: number
) {
  if (!groupId) return { error: "Group ID is required." };

  const result = await fetchApi<GroupApiResponse<null>>(
    `/api/groups/${groupId}/posts/${postId}/comments/${commentId}/reaction`,
    { method: "POST", body: { is_like: isLike } }
  );

  if (!result.success) return { error: result.error };
  return { success: true };
}

// ─── Get Group Members (returns enriched profiles) ───

export async function getGroupMembers(groupId: string) {
  if (!groupId) return { error: "Group ID is required." };

  const result = await fetchApi<GroupApiResponse<FeedAuthor[]>>(
    `/api/groups/members/${groupId}`,
    { method: "GET" }
  );

  if (!result.success) return { error: result.error };
  return { success: true, data: result.data.data };
}

// ─── Invite user to group ───

export async function inviteUserToGroup(groupId: string, userId: number) {
  if (!groupId) return { error: "Group ID is required." };

  const result = await fetchApi<GroupApiResponse<null>>(
    `/api/groups/${groupId}/invite`,
    { method: "POST", body: { user_id: userId } }
  );

  if (!result.success) return { error: result.error };
  return { success: true };
}

// ─── Get Invite Candidates (users the current user follows, not in the group) ───

export async function getInviteCandidates(groupId: string) {
  if (!groupId) return { error: "Group ID is required." };

  const result = await fetchApi<GroupApiResponse<
    { id: number; nickname: string; firstname: string; lastname: string; avatar: string }[]
  >>(`/api/groups/invite-candidates/${groupId}`, { method: "GET" });

  if (!result.success) return { error: result.error };
  return { success: true, data: result.data.data };
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

// ─── Update Group ───

export async function updateGroup(groupId: string, formData: FormData) {
  if (!groupId) return { error: "Group ID is required." };

  const result = await fetchApi<GroupApiResponse<null>>(
    `/api/groups/${groupId}/update`,
    { method: "PUT", body: formData }
  );

  if (!result.success) return { error: result.error };
  return { success: true };
}

// ─── Delete Group Post ───

export async function deleteGroupPost(groupId: string, postId: number) {
  if (!groupId) return { error: "Group ID is required." };

  const result = await fetchApi<GroupApiResponse<null>>(
    `/api/groups/${groupId}/posts/${postId}/delete`,
    { method: "POST" }
  );

  if (!result.success) return { error: result.error };
  return { success: true };
}

// ─── Delete Group Event ───

export async function deleteGroupEvent(groupId: string, eventId: number) {
  if (!groupId) return { error: "Group ID is required." };

  const result = await fetchApi<GroupApiResponse<null>>(
    `/api/groups/${groupId}/events/${eventId}/delete`,
    { method: "POST" }
  );

  if (!result.success) return { error: result.error };
  return { success: true };
}

// ─── Delete Group Post Comment ───

export async function deleteGroupPostComment(groupId: string, postId: number, commentId: number) {
  if (!groupId) return { error: "Group ID is required." };

  const result = await fetchApi<GroupApiResponse<null>>(
    `/api/groups/${groupId}/posts/${postId}/comments/${commentId}/delete`,
    { method: "POST" }
  );

  if (!result.success) return { error: result.error };
  return { success: true };
}

// ─── Kick Member ───

export async function kickMember(groupId: string, userId: number) {
  if (!groupId) return { error: "Group ID is required." };

  const result = await fetchApi<GroupApiResponse<null>>(
    `/api/groups/${groupId}/members/${userId}/kick`,
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