"use server";

import { fetchApi } from "../helper/fetch";

interface CreateGroupResponse {
  status_code: number;
  message: string;
  data: {
    group_id: number;
    title: string;
  };
}

// Create a new group using FormData (supports file uploads for logo and background)
export async function createGroup(formData: FormData) {
  const title = formData.get("title")?.toString().trim();
  const description = formData.get("description")?.toString().trim();

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

  return {
    success: true,
    message: result.data.message,
    group: result.data.data,
  };
}