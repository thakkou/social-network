"use client"
import { getNotifications } from "~/app/api/crud/notification";
import { useState, useEffect } from "react";

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

const NotificationCard = ({ data }: { data: NotificationItem }) => {
  console.log("notification of ", data.type, "the data:", data);

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

  // Generate clear user labels safely out of your actor object fields
  const getActorName = () => {
    if (!data.actor) return "Someone";
    if (data.actor.firstname || data.actor.lastname) {
      return `${data.actor.firstname} ${data.actor.lastname}`.trim();
    }
    return data.actor.nickname || "Someone";
  };

  // Safe user avatar abbreviation fallback
  const getInitials = () => {
    if (!data.actor) return "??";
    if (data.actor.firstname && data.actor.lastname) {
      return (data.actor.firstname[0] + data.actor.lastname[0]).toUpperCase();
    }
    return (data.actor.nickname ? data.actor.nickname.slice(0, 2) : "UN").toUpperCase();
  };

  const hasAccept = ["follow_request", "group_invite", "group_join_request"].includes(data.type);

  // Parse ISO date into a simple viewable string format
  const displayTime = data.created_at ? new Date(data.created_at).toLocaleDateString() : "just now";

  return (
    <div 
      className="card" 
      style={{ 
        borderLeft: `2px solid ${typeStyles.border}`,
        opacity: data.is_read ? 0.7 : 1 
      }}
    >
      <div style={{ display: 'flex', alignItems: 'center', gap: '8px', marginBottom: '8px' }}>
        <span className={`tag ${typeStyles.tagClass}`}>{typeStyles.label}</span>
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

      {/* Dynamic Action Buttons based on parsed types */}
      <div style={{ display: 'flex', gap: '6px' }}>
        {data.type === 'group_event' ? (
          <>
            <button className="btn btn-t" style={{ fontSize: '11px', display: 'flex', alignItems: 'center', gap: '3px' }}>
              <i className="ti ti-check" style={{ fontSize: '12px' }} aria-hidden="true"></i> going
            </button>
            <button className="btn btn-g" style={{ fontSize: '11px' }}>not going</button>
          </>
        ) : hasAccept ? (
          <>
            <button className="btn btn-t" style={{ fontSize: '11px', display: 'flex', alignItems: 'center', gap: '3px' }}>
              <i className="ti ti-check" style={{ fontSize: '12px' }} aria-hidden="true"></i> accept
            </button>
            <button className="btn btn-red" style={{ fontSize: '11px', display: 'flex', alignItems: 'center', gap: '3px' }}>
              <i className="ti ti-x" style={{ fontSize: '12px' }} aria-hidden="true"></i> decline
            </button>
          </>
        ) : (
          ["post_reaction", "comment"].includes(data.type) && (
            <button className="btn btn-g" style={{ fontSize: '11px' }}>view post</button>
          )
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

      <div style={{ display: "flex", alignItems: "center", justifyContent: "space-between", marginBottom: "12px" }}>
        <p style={{ fontSize: "13px", fontWeight: 500, color: "var(--color-text-primary)" }}>
          notifications ({notification?.length || 0})
        </p>
      </div>

      {notification.length === 0 ? (
        <div className="card" style={{ textAlign: "center", padding: "24px", color: "var(--color-text-secondary)" }}>
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