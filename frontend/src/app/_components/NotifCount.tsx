"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { getNotifications } from "~/app/api/crud/notification";

export default function NotifCount() {
  const [unread, setUnread] = useState(0);

  useEffect(() => {
    const fetchCount = async () => {
      const res = await getNotifications("unread");
      if (res.success && res.data) {
        setUnread(res.data.length);
      }
    };
    void fetchCount();

    // Poll every 30 seconds
    const interval = setInterval(fetchCount, 30000);
    return () => clearInterval(interval);
  }, []);

  return (
    <Link
      className="btn btn-g"
      style={{
        display: "flex",
        alignItems: "center",
        gap: "4px",
      }}
      href="/notifications"
    >
      <i className="ti ti-bell" style={{ fontSize: "14px" }} aria-hidden="true" />
      {unread > 0 && <span className="notif-dot">{unread}</span>}
    </Link>
  );
}
