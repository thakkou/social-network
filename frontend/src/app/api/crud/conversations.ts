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
  if (!id) {
    return { error: "Conversation ID is required." };
  }

  const params = new URLSearchParams({
    offset: String(offset),
    limit: String(limit),
  });

  const result = await fetchApi<GetConversationByIdResponse>(
    `/api/conversations/${type}/${id}?${params.toString()}`,
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