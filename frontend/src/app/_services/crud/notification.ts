"use server";

import { fetchApi } from "../helper/fetch";

/* =========================
 * Types
 * ========================= */

export type ActorInfo = {
  user_id: number;
  nickname: string;
  avatar: string;
  firstname: string;
  lastname: string;
}

export type Notification = {
  id: number;
  type:
    | "post_reaction"
    | "comment"
    | "follow_request"
    | "follow_accepted"
    | "group_invite"
    | "group_join_request";
  object_type: string;
  is_read: boolean;
  created_at: string;
  actor: ActorInfo;
  payload: any;
}

type NotificationResponse = {
  status_code: number;
  message: string;
  data: Notification[];
}

/* =========================
 * Get Notifications
 * ========================= */

export async function getNotifications(type: "all" | "unread" = "unread") {
  const result = await fetchApi<NotificationResponse>(
    `/api/notifications?type=${type}`,
    {
      method: "GET",
    }
  );
  console.log(result)

  if (!result.success) {
    return {
      success: false,
      error: result.error,
    };
  }

  return {
    success: true,
    data: result.data.data,
  };
}


/* =========================
* Mark One Notification Read
 * POST /api/notifications/read?id=1
 * ========================= */

export async function markNotificationAsRead(id: number | string) {
  const result = await fetchApi(
    `/api/notifications/read?id=${id}`,
    {
      method: "POST",
    }
  );

  if (!result.success) {
    return {
      success: false,
      error: result.error,
    };
  }

  return {
    success: true,
  };
}

/* =========================
 * Mark All Notifications Read
 * POST /api/notifications/read-all
 * ========================= */

export async function markAllNotificationsRead() {
  const result = await fetchApi(
    "/api/notifications/readAll",
    {
      method: "POST",
    }
  );

  if (!result.success) {
    return {
      success: false,
      error: result.error,
    };
  }

  return {
    success: true,
  };}


/* =========================
 * Delete All Notifications
 * DELETE /api/notifications
 * ========================= */

export async function deleteAllNotifications() {
  const result = await fetchApi(
    "/api/notifications/deletAll",
    {
      method: "DELETE",
    }
  );

  if (!result.success) {
    return {
      success: false,
      error: result.error,
    };
  }

  return {
    success: true,
  };
}
