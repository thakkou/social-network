"use server";

import { fetchApi } from "../helper/fetch";

interface SearchProfile {
  id: number;
  firstname: string;
  lastname: string;
  nickname?: string;
  avatar?: string;
  is_private?: number;
}

interface SearchGroup {
  id: number;
  title: string;
  description?: string;
  logo?: string;
}

interface SearchResponse {
  status_code: number;
  message: string;
  data: {
    profiles: SearchProfile[];
    groups: SearchGroup[];
  };
}

export async function search(text: string) {
  const query = text.trim();

  if (!query) {
    return { error: "Search query is required." };
  }

  const result = await fetchApi<SearchResponse>("/api/search", {
    method: "GET",
    searchParams: { text: query },
  });

  if (!result.success) {
    return { error: result.error };
  }

  return {
    success: true,
    data: result.data.data,
  };
}