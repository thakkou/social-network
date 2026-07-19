import React from "react";

export const MessagesSidebar: React.ComponentType<any> = ({
  users = [
    { id: "u1", name: "Selin Rauf", initials: "SR", color: "#FBEAF0" },
    { id: "u2", name: "Adam Smith", initials: "AS", color: "#EAF3FB" },
    { id: "u3", name: "Maya Ali", initials: "MA", color: "#EAFBEF" },
    { id: "u4", name: "John Doe", initials: "JD", color: "#FFF3E8" },
  ],
  groups = [
    { id: "g1", name: "go devs", color: "#D4537E" },
    { id: "g2", name: "open src", color: "#1D9E75" },
  ],
  activeId,
  onSelect,
}) => (
  <aside className="sidebar2">
    {/* Section 1: Users */}
    <p className="sec-label" style={{ padding: 0, marginBottom: "8px" }}>
      users
    </p>
    <div style={{ fontSize: "11px", display: "flex", flexDirection: "column", gap: "6px" }}>
      {users.map((user: any) => {
        const isSelected = activeId === user.id;
        return (
          <div
            key={user.id}
            onClick={() => onSelect?.(user)}
            style={{
              display: "flex",
              alignItems: "center",
              gap: "8px",
              padding: "6px",
              cursor: "pointer",
              background: isSelected ? "var(--color-background-tertiary)" : "var(--color-background-secondary)",
              border: "0.5px solid var(--color-border-tertiary)",
              fontWeight: isSelected ? 500 : "normal",
            }}
          >
            <div
              style={{
                width: "16px",
                height: "16px",
                borderRadius: "50%",
                background: user.color,
                color: "#993556",
                fontSize: "8px",
                display: "flex",
                alignItems: "center",
                justifyContent: "center",
                fontWeight: 600,
                flexShrink: 0,
              }}
            >
              {user.initials}
            </div>
            <span style={{ color: "var(--color-text-primary)" }}>{user.name}</span>
          </div>
        );
      })}
    </div>

    <div className="divider" style={{ margin: "12px 0" }}></div>

    {/* Section 2: Shared Groups */}
    <p className="sec-label" style={{ padding: 0, marginBottom: "8px" }}>
      shared groups
    </p>
    <div style={{ fontSize: "11px", display: "flex", flexDirection: "column", gap: "6px" }}>
      {groups.map((group: any) => {
        const isSelected = activeId === group.id;
        return (
          <div
            key={group.id}
            onClick={() => onSelect?.(group)}
            style={{
              display: "flex",
              alignItems: "center",
              gap: "6px",
              padding: "6px",
              cursor: "pointer",
              background: isSelected ? "var(--color-background-tertiary)" : "var(--color-background-secondary)",
              border: "0.5px solid var(--color-border-tertiary)",
              fontWeight: isSelected ? 500 : "normal",
            }}
          >
            <span
              style={{
                width: "6px",
                height: "6px",
                background: group.color,
                flexShrink: 0,
              }}
            ></span>
            <span style={{ color: "var(--color-text-primary)" }}>{group.name}</span>
          </div>
        );
      })}
    </div>
  </aside>
);