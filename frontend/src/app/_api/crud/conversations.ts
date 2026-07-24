"use server";

import { fetchApi } from "../helper/fetch";

export type ConversationType = "direct" | "group";

export interface ConversationFeedItem {
  type: ConversationType;
  id: number;
  display_name: string;
  avatar?: string;
  last_message?: string;
  last_message_at?: string;
  unread_count: number;
  rank: number;
  other_user_id?: number; // direct only
  member_count?: number; // group only
}

export interface ConversationMessage {
  id: number;
  sender_id: number;
  text: string;
  created_at: string;
  nickname: string;
}

interface GetConversationsResponse {
  status_code: number;
  message: string;
  data: {
    groups: ConversationFeedItem[];
    direct: ConversationFeedItem[];
  };
}

interface GetConversationByIdResponse {
  status_code: number;
  message: string;
  data: {
    type: ConversationType;
    id: number;
    messages: ConversationMessage[];
  };
}

export interface SendMessageRequest {
  type: ConversationType;
  text: string;
  receiver_id?: number;
  conversation_id?: number;
  group_id?: number;
}

interface SendMessageResponse {
  status_code: number;
  message: string;
  data: {
    type: ConversationType;
    message_id: number;
    conversation_id?: number;
    group_id?: number;
  };
}
// Get conversations (groups + direct), paginated & ranked by last activity
export async function getConversations(offset = 0, limit = 30) {
  const params = new URLSearchParams({
    offset: String(offset),
    limit: String(limit),
  });

  const result = await fetchApi<GetConversationsResponse>(
    `/api/conversations?${params.toString()}`,
    {
      method: "GET",
    }
  );

  if (!result.success) {
    return { error: result.error };
  }

  return {
    success: true,
    groups: result.data.data.groups,
    direct: result.data.data.direct,
  };
}

// Get messages for a single conversation (direct or group), paginated
export async function getConversationById(
  type: ConversationType,
  id: string | number,
  offset = 0,
  limit = 10
) {
  console.log("starrt get the conv By ID",type,id,offset,limit)
  if (!id) {
    return { error: "Conversation ID is required." };
  }

  const params = new URLSearchParams({
    offset: String(offset),
    limit: String(limit),
  });

  const result = await fetchApi<GetConversationByIdResponse>(
    `/api/conversation/${type}/${id}?${params.toString()}`,
    {
      method: "GET",
    }
  );

  if (!result.success) {
    return { error: result.error };
  }

  return {
    success: true,
    type: result.data.data.type,
    id: result.data.data.id,
    messages: result.data.data.messages,
  };
}

// Send a message (direct or group)
export async function sendMessage(payload: SendMessageRequest) {
  console.log("start sending message", payload);

  // Frontend validation to match backend requirements
  if (!payload.text.trim()) {
    return { error: "Message text cannot be empty." };
  }

  if (payload.type === "direct" && !payload.receiver_id) {
    return { error: "Receiver ID is required for direct messages." };
  }

  if (payload.type === "group" && !payload.group_id) {
    return { error: "Group ID is required for group messages." };
  }

  const result = await fetchApi<SendMessageResponse>("/api/messages", {
    method: "POST",
    body: JSON.stringify(payload),
  });

  if (!result.success) {
    return { error: result.error };
  }

  return {
    success: true,
    data: result.data.data,
  };
}