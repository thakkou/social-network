"use server";

import { fetchApi } from "../helper/fetch";

type FollowStatus = "none" | "pending" | "accepted";



interface FollowActionResponse {
  status_code: number;
  message: string;
  data: {
    status: FollowStatus;
  };
}


// 2. Send a follow request
async function followUser(userId: string | number) {
  const result = await fetchApi<FollowActionResponse>(
    `/api/follow/follow/${userId}`,
    {
      method: "PUT",
    }
  );

  if (!result.success) {
    return { error: result.error };
  }

  return {
    success: true,
    status: result.data.data.status,
  };
}
// 3. Remove a follow / cancel a follow request
async function unfollowUser(userId: string | number) {
  const result = await fetchApi<FollowActionResponse>(
    `/api/follow/unfollow/${userId}`,
    {
      method: "PUT",
    }
  );

  if (!result.success) {
    return { error: result.error };
  }

  return {
    success: true,
    status: result.data.data.status,
  };
}

// 4. Toggle: checks the current status, then follows if "none", otherwise unfollows
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
    return { error: result.error };
  }

  return {
    success: true,
    status: result.status,
  };
}