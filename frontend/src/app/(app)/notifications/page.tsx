
"use client"
import { getNotifications, markNotificationAsRead, markAllNotificationsRead, deleteAllNotifications } from "~/app/api/crud/notification";
import { useState, useEffect } from "react"
import { acceptFollowRequest,rejectFollowRequest } from "~/app/api/crud/follow";

// 1. Exact Interface mapping to your JSON response
interface NotificationActor {
  user_id: number;
  nickname: string;
  avatar: string;
  firstname: string;
  lastname: string;
}

interface NotificationPayload {
  post_id?: number;
  post_title?: string;
  reaction?: string;
  group_id?: number;
  group_name?: string;
  group_avatar?: string;
  invitation_id?: number;
  follow_status?: string;
  follow_request_id?: number;
}

interface NotificationItem {
  id: string | number;
  type: "follow_request" | "follow_accepted" | "group_invite" | "group_join_request" | "post_reaction" | "comment" | "group_event";
  object_type: string;
  is_read: boolean;
  created_at: string;
  actor: NotificationActor | null;
  payload: NotificationPayload | null;
}
const NotificationCard = ({
  data,
  onMarkRead,
  onAcceptFollow,
  onRejectFollow,
}: {
  data: NotificationItem;
  onMarkRead: (id: string | number) => void;
  onAcceptFollow: (userId: number, notificationId: string | number) => void;
  onRejectFollow: (userId: number, notificationId: string | number) => void;
}) => {
  // Setup dynamic color styling and configurations based on notification type
  const typeStyles = {
    follow_request: { border: '#D4537E', tagClass: 'tag-pink', label: 'follow request' },
    follow_accepted: { border: '#1D9E75', tagClass: 'tag-teal', label: 'follow accepted' },
    group_invite: { border: '#7F77DD', tagClass: 'tag-purple', label: 'group invite' },
    group_join_request: { border: '#7F77DD', tagClass: 'tag-purple', label: 'join request' },    
    post_reaction: { border: '#E28743', tagClass: 'tag-orange', label: 'reaction' },
    comment: { border: '#3A86FF', tagClass: 'tag-blue', label: 'comment' },
    group_event: { border: '#1D9E75', tagClass: 'tag-teal', label: 'group event' }
  }[data.type] || { border: '#ccc', tagClass: 'tag-g', label: 'notification' };

  const getActorName = () => {
    if (!data.actor) return "Someone";
    if (data.actor.firstname || data.actor.lastname) {
      return `${data.actor.firstname} ${data.actor.lastname}`.trim();
    }
    return data.actor.nickname || "Someone";
  };

  const getInitials = () => {
    if (!data.actor) return "??";
    if (data.actor.firstname && data.actor.lastname) {
      return (data.actor.firstname.charAt(0) + data.actor.lastname.charAt(0)).toUpperCase();
    }
    return (data.actor.nickname ? data.actor.nickname.slice(0, 2) : "UN").toUpperCase();
  };

  const hasAccept = ["follow_request", "group_invite", "group_join_request"].includes(data.type);
  const displayTime = data.created_at ? new Date(data.created_at).toLocaleDateString() : "just now";

  return (
    <div 
      className="card" 
      style={{ 
        borderLeft: `2px solid ${typeStyles.border}`,
        opacity: data.is_read ? 0.6 : 1 
      }}
    >
      <div style={{ display: 'flex', alignItems: 'center', gap: '8px', marginBottom: '8px' }}>
        <span className={`tag ${typeStyles.tagClass}`}>{typeStyles.label}</span>
        {data.is_read && <span style={{ fontSize: '10px', color: 'var(--color-text-tertiary)', fontStyle: 'italic' }}>read</span>}
        <span style={{ fontSize: '10px', color: 'var(--color-text-tertiary)', marginLeft: 'auto' }}>
          {displayTime}
        </span>
      </div>

      <div style={{ display: 'flex', alignItems: 'center', gap: '8px', marginBottom: '8px' }}>
        {data.actor && (
          <div className="av" style={{ width: '28px', height: '28px', background: '#FBEAF0', color: '#993556', fontSize: '11px', flexShrink: 0, overflow: 'hidden' }}>
            {data.actor.avatar ? (
              <img src={data.actor.avatar} alt="avatar" style={{ width: '100%', height: '100%', objectFit: 'cover' }} />
            ) : (
              getInitials()
            )}
          </div>
        )}
        <p style={{ fontSize: '12px', color: 'var(--color-text-primary)' }}>
          {data.type === 'follow_request' && (
            <>
              <span style={{ fontWeight: 500 }}>{getActorName()}</span> sent you a follow request
            </>
          )}
          {data.type === 'follow_accepted' && (
            <>
              <span style={{ fontWeight: 500 }}>{getActorName()}</span> accepted your follow request
            </>
          )}
          {data.type === 'group_invite' && (
            <>
              <span style={{ fontWeight: 500 }}>{getActorName()}</span> invited you to join the group <span style={{ fontWeight: 500, color: typeStyles.border }}>{data.payload?.group_name || "a group"}</span>
            </>
          )}
          {data.type === 'group_join_request' && (
            <>
              <span style={{ fontWeight: 500 }}>{getActorName()}</span> requested to join your group
            </>
          )}
          {data.type === 'post_reaction' && (
            <>
              <span style={{ fontWeight: 500 }}>{getActorName()}</span> liked your post {data.payload?.post_title && <span style={{ fontStyle: 'italic' }}>"{data.payload.post_title}"</span>}
            </>
          )}
          {data.type === 'comment' && (
            <>
              <span style={{ fontWeight: 500 }}>{getActorName()}</span> commented on your post
            </>
          )}
          {data.type === 'group_event' && (
            <>
              A new event was created inside your group. Do you want to join?
            </>
          )}
        </p>
      </div>

      <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}>
        <div style={{ display: "flex", gap: "6px" }}>
    {data.type === "group_event" ? (
  <>
    <button className="btn btn-t" style={{ fontSize: "11px" }}>
      going
    </button>
    <button className="btn btn-g" style={{ fontSize: "11px" }}>
      not going
    </button>
  </>
) : data.type === "follow_request" ? (
  <>
    <button
      className="btn btn-t"
      style={{ fontSize: "11px" }}
      onClick={() => {
        if (data.actor) {
          onAcceptFollow(data.actor.user_id, data.id);
        }
      }}
    >
      accept
    </button>

    <button
      className="btn btn-red"
      style={{ fontSize: "11px" }}
      onClick={() => {
        if (data.actor) {
          onRejectFollow(data.actor.user_id, data.id);
        }
      }}
    >
      decline
    </button>
  </>
) : data.type === "group_invite" || data.type === "group_join_request" ? (
  <>
    <button
      className="btn btn-t"
      style={{ fontSize: "11px" }}
    >
      accept
    </button>

    <button
      className="btn btn-red"
      style={{ fontSize: "11px" }}
    >
      decline
    </button>
  </>
) : (
            ["post_reaction", "comment"].includes(data.type) && (
              <button className="btn btn-g" style={{ fontSize: "11px" }}>view post</button>
            )
          )}
        </div>

        {/* Updated individual button to handle 'mark as read' state instead of flat out hard-deleting */}
        {!data.is_read && (
          <button
            className="btn btn-g"
            style={{ fontSize: "11px" }}
            onClick={() => onMarkRead(data.id)}
          >
            <i className="ti ti-check" /> mark read
          </button>
        )}
      </div>
    </div>
  );
};

export default function Notifications() {
  const [filter, setFilter] = useState<"all" | "unread">("unread");
  const [loading, setLoading] = useState(true);
  const [notification, setNotifications] = useState<NotificationItem[]>([]);
  const [actionError, setActionError] = useState<string | null>(null);

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

  const handleMarkRead = async (id: string | number) => {
    const res = await markNotificationAsRead(id);

    if (res.success) {
      if (filter === "unread") {
        // If viewing only unread, filter it completely out of sight
        setNotifications(prev => prev.filter(n => n.id !== id));
      } else {
        // Otherwise, visually change its inline read-status values
        setNotifications(prev => prev.map(n => n.id === id ? { ...n, is_read: true } : n));
      }
    }
  };

  const handleMarkAllRead = async () => {
    const res = await markAllNotificationsRead();

    if (res.success) {
      if (filter === "unread") {
        setNotifications([]);
      } else {
        setNotifications(prev => prev.map(n => ({ ...n, is_read: true })));
      }
    }
  };

  const handleAcceptFollow = async (
  userId: number,
  notificationId: string | number
) => {
  setActionError(null);
  const res = await acceptFollowRequest(userId);

  if ("error" in res) {
    setActionError(res.error ?? "Failed to accept follow request");
  } else {
    // remove the handled request
    setNotifications((prev) =>
      prev.filter((n) => n.id !== notificationId)
    );
  }
};

const handleRejectFollow = async (
  userId: number,
  notificationId: string | number
) => {
  setActionError(null);
  const res = await rejectFollowRequest(userId);

  if ("error" in res) {
    setActionError(res.error ?? "Failed to reject follow request");
  } else {
    // remove the handled request
    setNotifications((prev) =>
      prev.filter((n) => n.id !== notificationId)
    );
  }
};
  const handleDeleteAll = async () => {
    if (!window.confirm("Are you sure you want to clear all your notifications? This cannot be undone.")) return;
    
    const res = await deleteAllNotifications();
    if (res.success) {
      setNotifications([]);
    }
  };

  const handleFilterChange = (type: "all" | "unread") => {
    setFilter(type);
    loadNotifications(type);
  };

  useEffect(() => {
    void loadNotifications(filter);
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

        {/* Wired up to trigger the backend API delete-all route */}
        <button 
          className="btn btn-red" 
          style={{ fontSize: "10px", marginLeft: "auto" }}
          onClick={handleDeleteAll}
        >
          delete all
        </button>

        <button
          className="btn btn-g"
          style={{ fontSize: "10px" }}
          onClick={handleMarkAllRead}
        >
          mark all read
        </button>
      </div>

      <div style={{ display: "flex", alignItems: "center", justifyContent: "space-between", marginBottom: "12px" }}>
        <p style={{ fontSize: "13px", fontWeight: 500, color: "var(--color-text-primary)" }}>
          notifications ({notification?.length || 0})
        </p>
      </div>

      {actionError && (
        <div
          className="card"
          style={{
            borderLeft: "2px solid #e74c3c",
            padding: "10px 14px",
            fontSize: "12px",
            color: "#e74c3c",
            marginBottom: "12px",
          }}
        >
          {actionError}
        </div>
      )}

      {notification.length === 0 ? (
        <div className="card" style={{ textAlign: "center", padding: "24px", color: "var(--color-text-secondary)" }}>
          No notifications at the moment.
        </div>
      ) : (
        <div style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
          {notification.map((item) => (
        <NotificationCard
  key={item.id}
  data={item}
  onMarkRead={handleMarkRead}
  onAcceptFollow={handleAcceptFollow}
  onRejectFollow={handleRejectFollow}
/>
          ))}
        </div>
      )}
    </main>
  );
}
