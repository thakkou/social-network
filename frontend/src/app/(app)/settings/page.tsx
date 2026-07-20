"use client";

import { useState, useEffect, useRef } from "react";
import { useSession } from "next-auth/react";
import Link from "next/link";
import { getProfileData } from "~/app/api/crud/getProfile";
import {
  getUserGroups,
  updateProfilePrivacy,
  updateProfileNickname,
} from "~/app/api/crud/groups";

interface GroupSummary {
  id: number;
  title: string;
  description: string;
  created_at: string;
  logo?: string;
  creator_id: number;
}

export default function Settings() {
  const { data: session } = useSession();
  const userId = session?.user?.id;

  const [profile, setProfile] = useState<any>(null);
  const [groups, setGroups] = useState<GroupSummary[]>([]);
  const [loading, setLoading] = useState(true);
  const [groupsLoading, setGroupsLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [saved, setSaved] = useState(false);

  // Form state
  const [nickname, setNickname] = useState("");
  const [aboutme, setAboutme] = useState("");
  const [avatarFile, setAvatarFile] = useState<File | null>(null);
  const [avatarPreview, setAvatarPreview] = useState<string | null>(null);
  const [isPrivate, setIsPrivate] = useState(false);

  const fileInputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    if (!userId) return;
    const load = async () => {
      setLoading(true);
      const res = await getProfileData(userId);
      if (res.success) {
        setProfile(res.data);
        setNickname(res.data.nickname || "");
        setAboutme(res.data.aboutme || "");
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
        setGroups(res.data || []);
      }
      setGroupsLoading(false);
    };
    void loadGroups();
  }, [userId]);

  const handleSave = async () => {
    setSaving(true);
    setSaved(false);

    const fd = new FormData();
    fd.set("nickname", nickname);
    fd.set("aboutme", aboutme);
    if (avatarFile) fd.set("avatar", avatarFile);

    const res = await updateProfileNickname(fd);

    if (res.success) {
      setSaved(true);
      setTimeout(() => setSaved(false), 3000);
    }
    setSaving(false);
  };

  const handleTogglePrivacy = async () => {
    const newVal = !isPrivate;
    setIsPrivate(newVal);
    const res = await updateProfilePrivacy(newVal);
    if (!res.success) setIsPrivate(!newVal);
  };

  const initials = profile
    ? (profile.firstname?.[0]?.toUpperCase() || "") +
      (profile.lastname?.[0]?.toUpperCase() || "")
    : "—";

  // Filter groups where user is the creator (admin)
  const adminGroups = groups.filter((g) => {
    return g.creator_id === Number(userId);
  });

  if (loading || groupsLoading) {
    return (
      <main className="main">
        <div className="card" style={{ textAlign: "center", padding: "32px" }}>
          <p style={{ fontSize: 12, color: "#a09c94" }}>Loading settings...</p>
        </div>
      </main>
    );
  }

  return (
    <main className="main">
      {/* Profile Update */}
      <div className="card">
        <p
          style={{
            fontSize: "12px",
            fontWeight: 600,
            color: "#e8e4dc",
            marginBottom: "16px",
            textTransform: "uppercase",
            letterSpacing: "0.5px",
          }}
        >
          Profile Settings
        </p>

        {/* Avatar */}
        <div
          style={{
            display: "flex",
            alignItems: "center",
            gap: "12px",
            marginBottom: "16px",
          }}
        >
          <div
            className="av"
            style={{
              width: "56px",
              height: "56px",
              background: avatarPreview
                ? `url(${avatarPreview}) center/cover`
                : profile?.avatar
                ? `url(${profile.avatar}) center/cover`
                : "#EEEDFE",
              color: "#534AB7",
              fontSize: "16px",
              overflow: "hidden",
              cursor: "pointer",
            }}
            onClick={() => fileInputRef.current?.click()}
          >
            {!avatarPreview && !profile?.avatar && initials}
          </div>
          <div>
            <p style={{ fontSize: "13px", fontWeight: 500, color: "#e8e4dc" }}>
              {profile?.firstname} {profile?.lastname}
            </p>
            <p style={{ fontSize: "11px", color: "#6b6760" }}>
              {profile?.email}
            </p>
          </div>
          <input
            ref={fileInputRef}
            type="file"
            accept="image/*"
            style={{ display: "none" }}
            onChange={(e) => {
              const file = e.target.files?.[0];
              if (file) {
                setAvatarFile(file);
                const reader = new FileReader();
                reader.onload = () =>
                  setAvatarPreview(reader.result as string);
                reader.readAsDataURL(file);
              }
            }}
          />
        </div>

        {/* Nickname */}
        <div style={{ marginBottom: "12px" }}>
          <label
            style={{
              fontSize: "11px",
              color: "#a09c94",
              display: "block",
              marginBottom: "4px",
            }}
          >
            Nickname
          </label>
          <input
            className="inp"
            value={nickname}
            onChange={(e) => setNickname(e.target.value)}
            placeholder="e.g. cool_user_42"
            style={{ fontSize: "12px" }}
          />
        </div>

        {/* About Me */}
        <div style={{ marginBottom: "16px" }}>
          <label
            style={{
              fontSize: "11px",
              color: "#a09c94",
              display: "block",
              marginBottom: "4px",
            }}
          >
            About Me
          </label>
          <textarea
            className="inp"
            rows={3}
            value={aboutme}
            onChange={(e) => setAboutme(e.target.value)}
            placeholder="Tell the world about yourself..."
            style={{ resize: "none", fontSize: "12px" }}
          />
        </div>

        {/* Privacy toggle */}
        <div
          style={{
            display: "flex",
            alignItems: "center",
            justifyContent: "space-between",
            marginBottom: "16px",
          }}
        >
          <div>
            <p style={{ fontSize: "12px", color: "#e8e4dc" }}>
              Private Profile
            </p>
            <p style={{ fontSize: "10px", color: "#6b6760" }}>
              Only followers can see your posts and info
            </p>
          </div>
          <button
            className={`btn ${isPrivate ? "btn-t" : "btn-g"}`}
            style={{ fontSize: "10px" }}
            onClick={() => void handleTogglePrivacy()}
          >
            <i
              className={`ti ${isPrivate ? "ti-lock" : "ti-lock-open"}`}
              style={{ fontSize: "12px" }}
            />{" "}
            {isPrivate ? "private" : "public"}
          </button>
        </div>

        {/* Save */}
        <div style={{ display: "flex", gap: "8px", alignItems: "center" }}>
          <button
            className="btn btn-p"
            style={{
              fontSize: "11px",
              display: "flex",
              alignItems: "center",
              gap: 4,
              opacity: saving ? 0.6 : 1,
            }}
            disabled={saving}
            onClick={() => void handleSave()}
          >
            <i className="ti ti-device-floppy" style={{ fontSize: "12px" }} />{" "}
            {saving ? "saving..." : "save changes"}
          </button>
          {saved && (
            <span
              style={{
                fontSize: "11px",
                color: "#1D9E75",
                display: "flex",
                alignItems: "center",
                gap: 4,
              }}
            >
              <i className="ti ti-check" /> saved
            </span>
          )}
        </div>
      </div>

      {/* My Groups */}
      <div className="card">
        <p
          style={{
            fontSize: "12px",
            fontWeight: 600,
            color: "#e8e4dc",
            marginBottom: "12px",
            textTransform: "uppercase",
            letterSpacing: "0.5px",
          }}
        >
          My Groups
        </p>

        {groups.length === 0 ? (
          <p style={{ fontSize: "12px", color: "#6b6760" }}>
            You haven&apos;t joined any groups yet.
          </p>
        ) : (
          <div style={{ display: "flex", flexDirection: "column", gap: "8px" }}>
            {groups.map((g) => (
              <div
                key={g.id}
                style={{
                  display: "flex",
                  alignItems: "center",
                  gap: "8px",
                  padding: "8px",
                  border: "0.5px solid #3a3733",
                  background: "#2e2b27",
                }}
              >
                {g.logo ? (
                  <img
                    src={g.logo}
                    alt={g.title}
                    style={{
                      width: "28px",
                      height: "28px",
                      borderRadius: "4px",
                      objectFit: "cover",
                    }}
                  />
                ) : (
                  <div
                    style={{
                      width: "28px",
                      height: "28px",
                      borderRadius: "4px",
                      background: "#D4537E",
                      display: "flex",
                      alignItems: "center",
                      justifyContent: "center",
                      color: "#fff",
                      fontSize: "11px",
                      fontWeight: 600,
                    }}
                  >
                    {g.title?.charAt(0)?.toUpperCase() || "G"}
                  </div>
                )}
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
                    <p style={{ fontSize: "10px", color: "#6b6760", marginTop: 2 }}>
                      {g.description}
                    </p>
                  )}
                </div>
                <Link
                  href={`/groups/${g.id}`}
                  className="btn btn-g"
                  style={{ fontSize: "10px" }}
                >
                  manage
                </Link>
              </div>
            ))}
          </div>
        )}
      </div>
    </main>
  );
}
