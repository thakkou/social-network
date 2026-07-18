import { auth } from "~/server/auth";

type ApiResult<T> =
  | { success: true; data: T; error?: never }
  | { success: false; error: string; data?: never };

interface FetchApiOptions extends Omit<RequestInit, "body"> {
  body?: unknown; // Allow passing objects directly instead of stringifying manually
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

    // 2. Setup Headers (Merge default session cookies with custom headers)
    const headers = new Headers(options.headers);
    if (session?.user?.session_id) {
      headers.set("Cookie", `session_id=${session.user.session_id}`);
    }
    
    // Automatically set Content-Type if sending JSON body
    if (options.body && typeof options.body === "object" && !headers.has("Content-Type")) {
      headers.set("Content-Type", "application/json");
    }

    // 3. Prepare Request Configuration
    const config: RequestInit = {
      ...options,
      headers,
      credentials: options.credentials ?? "include",
      body: options.body && typeof options.body === "object" 
        ? JSON.stringify(options.body) 
        : (options.body as BodyInit),
    };

    // 4. Fire Request
    const res = await fetch(url.toString(), config);

    // 5. Handle Network/Backend Failure
    if (!res.ok) {
      const errorData = await res.json().catch(() => null);
      return {
        success: false,
        error: errorData?.error ?? `Request failed with status ${res.status}`,
      };
    }

    // 6. Return Typed Data
    const data = (await res.json()) as T;
    return { success: true, data };

  } catch (err) {
    // Catch-all for network timeouts or unexpected parsing errors
    return {
      success: false,
      error: err instanceof Error ? err.message : "An unexpected error occurred",
    };
  }
}