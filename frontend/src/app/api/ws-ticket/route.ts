import { NextResponse } from "next/server";
import { auth } from "~/server/auth";

export async function POST() {
  const session = await auth();

  if (!session?.user?.session_id) {
    return NextResponse.json({ error: "not authenticated" }, { status: 401 });
  }

  const backendUrl = process.env.GO_BACKEND_URL;
  if (!backendUrl) {
    return NextResponse.json({ error: "backend not configured" }, { status: 500 });
  }

  const res = await fetch(`${backendUrl}/api/ws-ticket`, {
    method: "POST",
    headers: {
      Cookie: `session_id=${session.user.session_id}`,
    },
  });

  if (!res.ok) {
    const body = await res.json().catch(() => null);
    return NextResponse.json(
      { error: body?.message ?? "ticket creation failed" },
      { status: res.status }
    );
  }

  const body = await res.json();
  return NextResponse.json({ ticket: body.data.ticket });
}
