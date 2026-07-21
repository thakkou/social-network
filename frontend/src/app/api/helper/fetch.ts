import { auth } from "~/server/auth";

type ApiResult<T> =
  | { success: true; data: T; error?: never }
  | { success: false; error: string; data?: never };

interface FetchApiOptions extends Omit<RequestInit, "body"> {
  body?: unknown; // Allow passing objects or FormData directly
  searchParams?: Record<string, string | number | boolean | undefined>;
}

export async function fetchApi<T>(
  endpoint: string,
  options: FetchApiOptions = {}
): Promise<ApiResult<T>> {
  try {
    const session = await auth();
    const backendUrl = process.env.GO_BACKEND_URL;

    if (!backendUrl) {
      return { success: false, error: "Backend URL configuration missing" };
    }

    // 1. Build URL & Append Query Parameters dynamically
    const url = new URL(endpoint, backendUrl);
    if (options.searchParams) {
      Object.entries(options.searchParams).forEach(([key, value]) => {
        if (value !== undefined) {
          url.searchParams.set(key, String(value));
        }
      });
    }

    // 2. Setup Headers
    const headers = new Headers(options.headers);
    if (session?.user?.session_id) {
      headers.set("Cookie", `session_id=${session.user.session_id}`);
    }

    const isFormData = options.body instanceof FormData;

    // Automatically set Content-Type to JSON ONLY if body is a plain object (not FormData)
    if (
      options.body &&
      typeof options.body === "object" &&
      !isFormData &&
      !headers.has("Content-Type")
    ) {
      headers.set("Content-Type", "application/json");
    }

    // 3. Prepare Request Body
    let body: BodyInit | null | undefined;
    if (isFormData) {
      body = options.body as FormData;
    } else if (options.body && typeof options.body === "object") {
      body = JSON.stringify(options.body);
    } else {
      body = options.body as BodyInit;
    }

    // 4. Prepare Request Configuration
    const config: RequestInit = {
      ...options,
      headers,
      credentials: options.credentials ?? "include",
      body,
    };

    // 5. Fire Request
    const res = await fetch(url.toString(), config);

    // 6. Handle Network/Backend Failure
    if (!res.ok) {
      const errorData = await res.json().catch(() => null);
      return {
        success: false,
        error: errorData?.message ?? errorData?.error ?? `Request failed with status ${res.status}`,
      };
    }

    // 7. Return Typed Data
    const data = (await res.json()) as T;
    return { success: true, data };

  } catch (err) {
    return {
      success: false,
      error: err instanceof Error ? err.message : "An unexpected error occurred",
    };
  }
}