"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { getNotifications } from "~/app/api/crud/notification";
import { useWS } from "~/app/_providers/ws-provider";

const NOTIF_EVENTS = [
  "like_posts",
  "new_comments",
  "new_posts",
  "group_event",
  "group_invite",
  "group_join_request",
  "follow_request",
  "follow_accepted",
];

export default function NotifCount() {
  const [unread, setUnread] = useState(0);
  const { on } = useWS();

  // Initial fetch
  useEffect(() => {
    const fetchCount = async () => {
      const res = await getNotifications("unread");
      if (res.success && res.data) {
        setUnread(res.data.length);
      }
    };
    void fetchCount();
  }, []);

  // Listen for WS notification events to increment count live
  useEffect(() => {
    const unsubs = NOTIF_EVENTS.map((event) =>
      on(event, () => {
        setUnread((prev) => prev + 1);
      })
    );

    return () => unsubs.forEach((fn) => fn());
  }, [on]);

  return (
    <Link
      className="btn btn-g"
      style={{
        display: "flex",
        alignItems: "center",
        gap: "4px",
      }}
      href="/notifications"
      onClick={() => setUnread(0)}
    >
      <i className="ti ti-bell" style={{ fontSize: "14px" }} aria-hidden="true" />
      {unread > 0 && <span className="notif-dot">{unread}</span>}
    </Link>
  );
}
