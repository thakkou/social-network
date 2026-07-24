"use server";

import { fetchApi } from "../helper/fetch";

type FollowStatus = "none" | "pending" | "accepted";

type FollowActionResponse = {
  status_code: number;
  message: string;
  data: {
    status?: FollowStatus;
  };
}

async function followAction(
  action: "follow" | "unfollow" | "accept" | "reject",
  userId: string | number
) {
  const result = await fetchApi<FollowActionResponse>(
    `/api/follow/${action}/${userId}`,
    {
      method: "PUT",
    }
  );

  if (!result.success) {
    return { error: result.error };
  }

  return {
    success: true,
    message: result.data.message,
    status: result.data.data?.status,
  };
}

// Follow
export async function followUser(userId: string | number) {
  return followAction("follow", userId);
}

// Unfollow / Cancel request
export async function unfollowUser(userId: string | number) {
  return followAction("unfollow", userId);
}

// Accept follow request
export async function acceptFollowRequest(userId: string | number) {
  return followAction("accept", userId);
}

// Reject follow request
export async function rejectFollowRequest(userId: string | number) {
  return followAction("reject", userId);
}

// Toggle follow
export async function toggleFollow(
  userId: string | number,
  followingState: FollowStatus
) {
  if (!userId) {
    return { error: "User ID is required to toggle follow." };
  }

  const result =
    followingState === "none"
      ? await followUser(userId)
      : await unfollowUser(userId);

  if ("error" in result) {
    return result;
  }

  return {
    success: true,
    status: result.status,
  };
}