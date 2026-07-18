"use server";
import { auth } from "~/server/auth";
import { ur } from "zod/v4/locales";

type SearchResult<T = unknown> =
  | { success: true; data: T }
  | { error: string };

export async function search(text: string): Promise<SearchResult> {

const session = await auth();

  const query = text.trim();

  if (!query) {
    return { error: "Search query is required." };
  }

  const url = new URL(
    "/api/search",
    process.env.GO_BACKEND_URL,
  );

  url.searchParams.set("text", query);

  console.log(url)

const res = await fetch(url.toString(), {
  method: "GET",
  headers: {
    Cookie: `session_id=${session?.user?.session_id}`,
  },
  credentials: "include",
});

  console.log("the result of search is",res)
  if (!res.ok) {
    const data = await res.json().catch(() => null);

    return {
      error: data?.error ?? "Search failed",
    };
  }

  const data = await res.json();

  return {
    success: true,
    data,
  };
}