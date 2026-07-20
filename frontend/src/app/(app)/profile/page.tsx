"use client";

import React, { useState, useEffect } from "react";
import { useSession } from "next-auth/react";
import Link from "next/link";
import { getProfileData } from "~/app/api/crud/getProfile";
import { getUserGroups, updateProfilePrivacy } from "~/app/api/crud/groups";
import ProfilePosts from "~/app/_components/ProfilePosts";

interface GroupSummary {
  id: number;
  title: string;
  description: string;
  created_at: string;
}

export default function Profile() {
  const { data: session } = useSession();
  const userId = session?.user?.id;

  const [profile, setProfile] = useState<any>(null);
  const [isPrivate, setIsPrivate] = useState(false);
  const [loading, setLoading] = useState(true);
  const [groups, setGroups] = useState<GroupSummary[]>([]);
  const [groupsLoading, setGroupsLoading] = useState(true);
  const [activeTab, setActiveTab] = useState<"posts" | "groups">("posts");

  useEffect(() => {
    if (!userId) return;

    const load = async () => {
      setLoading(true);
      const res = await getProfileData(userId);
      if (res.success) {
        setProfile(res.data);
        setIsPrivate(res.data.is_private === 1);
      }
      setLoading(false);
    };

    void load();
  }, [userId]);

  useEffect(() => {
    if (!userId) return;
    const loadGroups = async () => {
      setGroupsLoading(true);
      const res = await getUserGroups();
      if (res.success) {
        setGroups(res.data);
      }
      setGroupsLoading(false);
    };
    void loadGroups();
  }, [userId]);

  const handleTogglePrivacy = async () => {
    const newVal = !isPrivate;
    // Optimistic update
    setIsPrivate(newVal);
    const res = await updateProfilePrivacy(newVal);
    if (!res.success) {
      setIsPrivate(!newVal); // revert on failure
    }
  };

  const initials =
    profile?.firstname?.[0]?.toUpperCase() +
      profile?.lastname?.[0]?.toUpperCase() || "—";

  const name =
    `${profile?.firstname || ""} ${profile?.lastname || ""}`.trim() || "User";

  if (loading) {
    return (
      <main className="main">
        <div className="card" style={{ textAlign: "center", padding: "32px" }}>
          <p style={{ fontSize: 12, color: "#a09c94" }}>Loading profile...</p>
        </div>
      </main>
    );
  }

  return (
    <main className="main">
      {/* Profile Card */}
      <div className="card">
        <div style={{ display: "flex", alignItems: "flex-start", gap: "12px", marginBottom: "12px" }}>
          <div
            className="av"
            style={{
              width: "52px",
              height: "52px",
              background: profile?.avatar
                ? `url(${profile.avatar}) center/cover`
                : "#EEEDFE",
              color: "#534AB7",
              fontSize: "16px",
              overflow: "hidden",
            }}
          >
            {!profile?.avatar && initials}
          </div>

          <div style={{ flex: 1 }}>
            <div
              style={{
                display: "flex",
                alignItems: "center",
                gap: "8px",
                flexWrap: "wrap",
                marginBottom: "4px",
              }}
            >
              <p style={{ fontSize: "14px", fontWeight: 500 }}>{name}</p>
              {profile?.nickname && (
                <span className="tag tag-teal">@{profile.nickname}</span>
              )}
              <span
                className={`tag ${isPrivate ? "tag-gray" : "tag-purple"}`}
                id="profile-visibility-tag"
              >
                {isPrivate ? "private" : "public"}
              </span>
              <button
                className="btn btn-g"
                style={{
                  fontSize: "10px",
                  display: "flex",
                  alignItems: "center",
                  gap: "3px",
                  marginLeft: "auto",
                }}
                id="visibility-btn"
                onClick={() => void handleTogglePrivacy()}
              >
                <i
                  className={`ti ${isPrivate ? "ti-lock-open" : "ti-lock"}`}
                  style={{ fontSize: "12px" }}
                  aria-hidden="true"
                />{" "}
                make {isPrivate ? "public" : "private"}
              </button>
            </div>

            {profile?.birthdate && (
              <p style={{ fontSize: "11px", color: "#a09c94", marginBottom: "4px" }}>
                Born {profile.birthdate} · {profile?.email || ""}
              </p>
            )}

            {profile?.aboutme && (
              <p
                style={{
                  fontSize: "12px",
                  color: "#e8e4dc",
                  lineHeight: 1.5,
                }}
              >
                {profile.aboutme}
              </p>
            )}
          </div>
        </div>

        <div className="divider" />

        <div
          style={{
            display: "grid",
            gridTemplateColumns: "repeat(3, 1fr)",
            gap: "8px",
            textAlign: "center",
          }}
        >
          <div style={{ background: "#2e2b27", padding: "8px" }}>
            <p style={{ fontSize: "18px", fontWeight: 500, color: "#D4537E" }}>
              {profile?.posts?.length || "0"}
            </p>
            <p style={{ fontSize: "10px", color: "#6b6760" }}>posts</p>
          </div>
          <div style={{ background: "#2e2b27", padding: "8px" }}>
            <p style={{ fontSize: "18px", fontWeight: 500, color: "#534AB7" }}>
              {profile?.followers?.length || "0"}
            </p>
            <p style={{ fontSize: "10px", color: "#6b6760" }}>followers</p>
          </div>
          <div style={{ background: "#2e2b27", padding: "8px" }}>
            <p style={{ fontSize: "18px", fontWeight: 500, color: "#0F6E56" }}>
              {profile?.following?.length || "0"}
            </p>
            <p style={{ fontSize: "10px", color: "#6b6760" }}>following</p>
          </div>
        </div>
      </div>

      {/* Tabs */}
      <div
        style={{
          display: "flex",
          gap: 0,
          border: "0.5px solid #3a3733",
          background: "#272420",
        }}
      >
        {(["posts", "groups"] as const).map((tab) => (
          <div
            key={tab}
            style={{
              padding: "8px 16px",
              fontSize: "11px",
              cursor: "pointer",
              borderRight: "0.5px solid #3a3733",
              color: activeTab === tab ? "#D4537E" : "#a09c94",
              borderBottom: activeTab === tab ? "2px solid #D4537E" : "none",
            }}
            onClick={() => setActiveTab(tab)}
          >
            {tab}
          </div>
        ))}
      </div>

      {/* Posts Tab */}
      {activeTab === "posts" && <ProfilePosts posts={profile?.posts ?? []} />}

      {/* Groups Tab */}
      {activeTab === "groups" && (
        <div style={{ display: "flex", flexDirection: "column", gap: "8px" }}>
          {groupsLoading ? (
            <div className="card" style={{ textAlign: "center", padding: "16px" }}>
              <p style={{ fontSize: 12, color: "#a09c94" }}>Loading groups...</p>
            </div>
          ) : groups.length === 0 ? (
            <div className="card" style={{ textAlign: "center", padding: "16px" }}>
              <p style={{ fontSize: 12, color: "#6b6760" }}>Not a member of any group yet.</p>
            </div>
          ) : (
            groups.map((g) => (
              <div key={g.id} className="card">
                <div style={{ display: "flex", alignItems: "center", gap: "8px" }}>
                  <span
                    style={{
                      width: "8px",
                      height: "8px",
                      background: "#D4537E",
                      flexShrink: 0,
                    }}
                  />
                  <div style={{ flex: 1 }}>
                    <Link
                      href={`/groups/${g.id}`}
                      style={{
                        fontSize: "13px",
                        fontWeight: 500,
                        color: "#e8e4dc",
                        textDecoration: "none",
                      }}
                    >
                      {g.title}
                    </Link>
                    {g.description && (
                      <p style={{ fontSize: "11px", color: "#6b6760", marginTop: 2 }}>
                        {g.description}
                      </p>
                    )}
                  </div>
                  <Link
                    href={`/groups/${g.id}`}
                    className="btn btn-g"
                    style={{ fontSize: "10px" }}
                  >
                    view
                  </Link>
                </div>
              </div>
            ))
          )}
        </div>
      )}
    </main>
  );
}
