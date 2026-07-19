"use client"
import { getNotifications } from "~/app/api/crud/notification";
import { useState, useEffect } from "react";

// 1. Refactored Card to be fully dynamic
interface NotificationItem {
  id: string | number;
  type: "follow_request" | "group_invite" | "join_request" | "event";
  title?: string;
  senderName?: string;
  senderInitials?: string;
  groupName?: string;
  timeAgo: string;
  isRead: boolean;
}

const NotificationCard = ({ data }: { data: NotificationItem }) => {
  // Setup dynamic color styling based on notification type
  const typeStyles = {
    follow_request: { border: '#D4537E', tagClass: 'tag-pink', label: 'follow request' },
    group_invite: { border: '#7F77DD', tagClass: 'tag-purple', label: 'group invite' },
    join_request: { border: '#7F77DD', tagClass: 'tag-purple', label: 'join request' },
    event: { border: '#1D9E75', tagClass: 'tag-teal', label: 'event' },
  }[data.type] || { border: '#ccc', tagClass: 'tag-g', label: 'notification' };

  return (
    <div 
      className="card" 
      style={{ 
        borderLeft: `2px solid ${typeStyles.border}`,
        opacity: data.isRead ? 0.7 : 1 
      }}
    >
      <div style={{ display: 'flex', alignItems: 'center', gap: '8px', marginBottom: '8px' }}>
        <span className={`tag ${typeStyles.tagClass}`}>{typeStyles.label}</span>
        <span style={{ fontSize: '10px', color: 'var(--color-text-tertiary)', marginLeft: 'auto' }}>
          {data.timeAgo}
        </span>
      </div>

      <div style={{ display: 'flex', alignItems: 'center', gap: '8px', marginBottom: '8px' }}>
        {data.senderInitials && (
          <div className="av" style={{ width: '28px', height: '28px', background: '#FBEAF0', color: '#993556', fontSize: '11px' }}>
            {data.senderInitials}
          </div>
        )}
        <p style={{ fontSize: '12px', color: 'var(--color-text-primary)' }}>
          {data.type === 'follow_request' && <>
            <span style={{ fontWeight: 500 }}>{data.senderName}</span> sent you a follow request
          </>}
          {data.type === 'group_invite' && <>
            <span style={{ fontWeight: 500 }}>{data.senderName}</span> invited you to join <span style={{ fontWeight: 500, color: typeStyles.border }}>{data.groupName}</span>
          </>}
          {data.type === 'join_request' && <>
            <span style={{ fontWeight: 500 }}>{data.senderName}</span> requested to join your group <span style={{ fontWeight: 500, color: '#D4537E' }}>{data.groupName}</span>
          </>}
          {data.type === 'event' && <>
            New event <span style={{ fontWeight: 500 }}>{data.title}</span> was created in <span style={{ fontWeight: 500, color: '#D4537E' }}>{data.groupName}</span>
          </>}
        </p>
      </div>

      {/* Dynamic Action Buttons */}
      <div style={{ display: 'flex', gap: '6px' }}>
        {data.type === 'event' ? (
          <>
            <button className="btn btn-t" style={{ fontSize: '11px', display: 'flex', alignItems: 'center', gap: '3px' }}>
              <i className="ti ti-check" style={{ fontSize: '12px' }} aria-hidden="true"></i> going
            </button>
            <button className="btn btn-g" style={{ fontSize: '11px' }}>not going</button>
          </>
        ) : (
          <>
            <button className="btn btn-t" style={{ fontSize: '11px', display: 'flex', alignItems: 'center', gap: '3px' }}>
              <i className="ti ti-check" style={{ fontSize: '12px' }} aria-hidden="true"></i> accept
            </button>
            <button className="btn btn-red" style={{ fontSize: '11px', display: 'flex', alignItems: 'center', gap: '3px' }}>
              <i className="ti ti-x" style={{ fontSize: '12px' }} aria-hidden="true"></i> decline
            </button>
          </>
        )}
      </div>
    </div>
  );
};

export default function Notifications() {
  const [filter, setFilter] = useState<"all" | "unread">("unread");
  const [loading, setLoading] = useState(true);
  const [notification, setNotifications] = useState<NotificationItem[]>([]);

  const loadNotifications = async (type: "all" | "unread") => {
    setLoading(true);
    const res = await getNotifications(type);

    if (!res.success || !res.data) {
      setNotifications([]);
      setLoading(false);
      return;
    }

    setNotifications(res.data);
    setLoading(false);
  };

  // Switch filter type cleanly
  const handleFilterChange = (type: "all" | "unread") => {
    setFilter(type);
    loadNotifications(type);
  };

  useEffect(() => {
    loadNotifications(filter);
  }, []);

  if (loading) {
    return (
      <main className="main">
        <p style={{ fontSize: 13 }}>Loading...</p>
      </main>
    );
  }

  return (
    <main className="main">
      {/* Dynamic Filters Bar - Accessible Always */}
      <div style={{ display: "flex", gap: "6px", marginBottom: "16px" }}>
        <button
          className={`btn ${filter === "all" ? "btn-t" : "btn-g"}`}
          style={{ fontSize: "10px" }}
          onClick={() => handleFilterChange("all")}
        >
          get all
        </button>

        <button
          className={`btn ${filter === "unread" ? "btn-t" : "btn-g"}`}
          style={{ fontSize: "10px" }}
          onClick={() => handleFilterChange("unread")}
        >
          unread
        </button>

        <button className="btn btn-red" style={{ fontSize: "10px", marginLeft: "auto" }}>
          delete all
        </button>

        <button className="btn btn-g" style={{ fontSize: "10px" }}>
          mark all read
        </button>
      </div>

      <div
        style={{
          display: "flex",
          alignItems: "center",
          justifyContent: "space-between",
          marginBottom: "12px",
        }}
      >
        <p style={{ fontSize: "13px", fontWeight: 500, color: "var(--color-text-primary)" }}>
          notifications ({notification?.length || 0})
        </p>
      </div>

      {/* Conditionally Render List vs Empty State */}
      {notification.length === 0 ? (
        <div
          className="card"
          style={{
            textAlign: "center",
            padding: "24px",
            color: "var(--color-text-secondary)",
          }}
        >
          No notifications at the moment.
        </div>
      ) : (
        <div style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
          {notification.map((item) => (
            <NotificationCard key={item.id} data={item} />
          ))}
        </div>
      )}
    </main>
  );
}