"use server";

import { fetchApi } from "../helper/fetch";

// 1. Define the structural types matching your Go API response
interface ProfileUserMinimal {
  id: number;
  firstname: string;
  nickname: string;
  avatar: string;
}

interface ProfilePost {
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
  is_liked: number; // 1 for true, 0 for false
  categories: string[] | null;
  comments: unknown[] | null;
}

interface ProfileDetails {
  id: number;
  firstname: string;
  lastname: string;
  nickname: string;
  avatar: string;
  aboutme: string;
  created_at: string;
  is_private: number; // 1 for true, 0 for false
  following_status: "none" | "pending" | "accepted";
  followers: ProfileUserMinimal[];
  following: ProfileUserMinimal[];
  posts: ProfilePost[];
}

// The top-level envelope returned by your Go backend
interface ProfileApiResponse {
  status_code: number;
  message: string;
  data: ProfileDetails;
}

// 2. The Server Action Function
export async function getProfileData(userId: string | number) {
  if (!userId) {
    return { error: "User ID is required to fetch profile data." };
  }

  // Uses clean template literals to pass dynamic ID to the route: /profile/2
  const result = await fetchApi<ProfileApiResponse>(`/api/profile/${userId}`, {
    method: "GET",
  });
  console.log("geting profile",result)

  if (!result.success) {
    return { error: result.error };
  }

  return {
    success: true,
    data: result.data.data, // Extracts the actual profile body directly for easier usage
  };
}