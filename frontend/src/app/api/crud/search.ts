"use server";

import { fetchApi } from "../helper/fetch";

// Define what your Go backend actually returns for search
interface SearchData {
  results: Array<{ id: string; name: string }>; 
}

export async function search(text: string) {
  const query = text.trim();

  if (!query) {
    return { error: "Search query is required." };
  }

  // Pass your expected generic type <SearchData> to get full type safety
  const result = await fetchApi<SearchData>("/api/search", {
    method: "GET",
    searchParams: { text: query },
  });

  if (!result.success) {
    return { error: result.error };
  }

  return {
    success: true,
    data: result.data,
  };
}