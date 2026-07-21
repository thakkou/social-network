"use client";

import { useState, useEffect, useRef } from "react";
import { useSession } from "next-auth/react";
import Link from "next/link";
import { getProfileData } from "~/app/api/crud/getProfile";
import {
  getUserGroups,
  updateProfilePrivacy,
  updateProfileNickname,
  updateGroup,
  inviteUserToGroup,
  getGroupMembers,
} from "~/app/api/crud/groups";
import { search } from "~/app/api/crud/search";

interface GroupSummary {
  id: number;
  title: string;
  description: string;
  created_at: string;
  logo?: string;
  background?: string;
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

  // ── Group edit state ──
  const [editingGroupId, setEditingGroupId] = useState<number | null>(null);
  const [editTitle, setEditTitle] = useState("");
  const [editDescription, setEditDescription] = useState("");
  const [editLogoFile, setEditLogoFile] = useState<File | null>(null);
  const [editBgFile, setEditBgFile] = useState<File | null>(null);
  const [editLogoPreview, setEditLogoPreview] = useState<string | null>(null);
  const [editBgPreview, setEditBgPreview] = useState<string | null>(null);
  const [savingGroup, setSavingGroup] = useState(false);
  const [savedGroupId, setSavedGroupId] = useState<number | null>(null);

  const logoInputRef = useRef<HTMLInputElement>(null);
  const bgInputRef = useRef<HTMLInputElement>(null);

  const startEditing = (g: GroupSummary) => {
    setEditingGroupId(g.id);
    setEditTitle(g.title);
    setEditDescription(g.description);
    setEditLogoFile(null);
    setEditBgFile(null);
    setEditLogoPreview(null);
    setEditBgPreview(null);
  };

  const handleSaveGroup = async (g: GroupSummary) => {
    setSavingGroup(true);
    const fd = new FormData();
    fd.set("title", editTitle);
    fd.set("description", editDescription);
    if (editLogoFile) fd.set("logo", editLogoFile);
    if (editBgFile) fd.set("background", editBgFile);

    const res = await updateGroup(String(g.id), fd);
    if (res.success) {
      setEditingGroupId(null);
      setSavedGroupId(g.id);
      const r = await getUserGroups();
      if (r.success) setGroups(r.data || []);
      setTimeout(() => setSavedGroupId(null), 2000);
    }
    setSavingGroup(false);
  };

  // ── Invite state ──
  const [inviteGroupId, setInviteGroupId] = useState<number | null>(null);
  const [inviteQuery, setInviteQuery] = useState("");
  const [inviteResults, setInviteResults] = useState<{ id: number; nickname: string; firstname: string; lastname: string; avatar: string }[]>([]);
  const [inviteSending, setInviteSending] = useState<Record<number, boolean>>({});
  const [inviteMsg, setInviteMsg] = useState<string | null>(null);
  const [memberIDs, setMemberIDs] = useState<number[]>([]);

  const openInviteModal = async (groupId: number) => {
    setInviteGroupId(groupId);
    setInviteQuery("");
    setInviteResults([]);
    setInviteMsg(null);
    // Fetch existing member profiles
    const res = await getGroupMembers(String(groupId));
    if (res.success) {
      setMemberIDs(res.data.map((m) => m.id));
    }
  };

  const handleInviteSearch = async (query: string) => {
    setInviteQuery(query);
    if (!query.trim()) {
      setInviteResults([]);
      return;
    }
    const res = await search(query);
    if (res.success) {
      setInviteResults(
        res.data.profiles
          .filter((u) => u.id !== Number(userId))
          .map((u) => ({
            id: u.id,
            nickname: u.nickname || `${u.firstname} ${u.lastname}`.trim(),
            firstname: u.firstname,
            lastname: u.lastname,
            avatar: u.avatar || "",
          }))
      );
    }
  };

  const handleSendInvite = async (targetUserId: number) => {
    if (!inviteGroupId) return;
    setInviteSending((prev) => ({ ...prev, [targetUserId]: true }));
    setInviteMsg(null);
    try {
      const res = await inviteUserToGroup(String(inviteGroupId), targetUserId);
      if (res.success) {
        setInviteMsg("invite sent!");
        setInviteResults([]);
        setInviteQuery("");
      } else {
        setInviteMsg(res.error ?? "failed");
      }
    } catch {
      setInviteMsg("something went wrong");
    } finally {
      setInviteSending((prev) => ({ ...prev, [targetUserId]: false }));
    }
  };

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

      {/* My Groups (only groups I created) */}
      <div className="card" style={{ padding: 0 }}>
        <p
          style={{
            fontSize: "12px",
            fontWeight: 600,
            color: "#e8e4dc",
            marginBottom: "12px",
            textTransform: "uppercase",
            letterSpacing: "0.5px",
            padding: "0.875rem 1rem 0",
          }}
        >
          My Groups
        </p>

        {adminGroups.length === 0 ? (
          <p style={{ fontSize: "12px", color: "#6b6760", padding: "0 1rem 0.875rem" }}>
            You haven&apos;t created any groups yet.
          </p>
        ) : (
          <div style={{ display: "flex", flexDirection: "column", gap: 0 }}>
            {adminGroups.map((g, idx) => (
              <div
                key={g.id}
                style={{
                  borderTop: idx > 0 ? "0.5px solid #3a3733" : "none",
                }}
              >
                {/* ── Group card hero ── */}
                <div
                  style={{
                    padding: 0,
                    overflow: "hidden",
                  }}
                >
                  {/* Background banner */}
                  <div
                    style={{
                      height: 64,
                      background: g.background
                        ? `url(${g.background}) center/cover`
                        : "linear-gradient(135deg,#2e2b27 0%, #272420 80%, #2e1e24 100%)",
                      borderBottom: "0.5px solid #3a3733",
                    }}
                  />

                  <div style={{ padding: "0 12px 10px" }}>
                    {/* Logo + title row */}
                    <div
                      style={{
                        display: "flex",
                        alignItems: "flex-end",
                        gap: 10,
                        marginTop: -28,
                      }}
                    >
                      <div
                        style={{
                          width: 48,
                          height: 48,
                          borderRadius: "50%",
                          background: g.logo
                            ? `url(${g.logo}) center/cover`
                            : "#D4537E",
                          fontSize: g.logo ? 0 : 18,
                          color: "#fff",
                          border: "3px solid #272420",
                          display: "flex",
                          alignItems: "center",
                          justifyContent: "center",
                          flexShrink: 0,
                          fontWeight: 600,
                        }}
                      >
                        {!g.logo && (g.title?.charAt(0)?.toUpperCase() || "G")}
                      </div>
                      <div style={{ flex: 1, minWidth: 0, paddingTop: 10 }}>
                        <Link
                          href={`/groups/${g.id}`}
                          style={{
                            fontSize: 14,
                            fontWeight: 600,
                            color: "#e8e4dc",
                            textDecoration: "none",
                          }}
                        >
                          {g.title}
                        </Link>
                        {g.description && (
                          <p
                            style={{
                              fontSize: 10,
                              color: "#6b6760",
                              marginTop: 1,
                              overflow: "hidden",
                              textOverflow: "ellipsis",
                              whiteSpace: "nowrap",
                            }}
                          >
                            {g.description}
                          </p>
                        )}
                      </div>
                    </div>

                    {/* Action buttons row */}
                    <div
                      style={{
                        display: "flex",
                        gap: 6,
                        marginTop: 10,
                      }}
                    >
                      <button
                        className="btn btn-g"
                        style={{ fontSize: 10 }}
                        onClick={() => void openInviteModal(g.id)}
                      >
                        <i className="ti ti-user-plus" style={{ fontSize: 11 }} /> invite
                      </button>
                      {editingGroupId === g.id ? (
                        <button
                          className="btn btn-t"
                          style={{ fontSize: 10 }}
                          disabled={savingGroup}
                          onClick={() => void handleSaveGroup(g)}
                        >
                          {savingGroup ? "saving..." : "save"}
                        </button>
                      ) : (
                        <button
                          className="btn btn-g"
                          style={{ fontSize: 10 }}
                          onClick={() => startEditing(g)}
                        >
                          <i className="ti ti-edit" style={{ fontSize: 11 }} /> edit
                        </button>
                      )}
                      {savedGroupId === g.id && (
                        <span
                          style={{
                            fontSize: 10,
                            color: "#1D9E75",
                            whiteSpace: "nowrap",
                            display: "flex",
                            alignItems: "center",
                            gap: 2,
                          }}
                        >
                          <i className="ti ti-check" /> saved
                        </span>
                      )}
                    </div>
                  </div>
                </div>

                {/* Edit form (expandable) */}
                {editingGroupId === g.id && (
                  <div
                    style={{
                      padding: "10px 12px",
                      borderTop: "0.5px solid #3a3733",
                      background: "#1e1c1a",
                    }}
                  >
                    <div style={{ marginBottom: "6px" }}>
                      <label
                        style={{
                          fontSize: "9px",
                          color: "#a09c94",
                          display: "block",
                          marginBottom: "2px",
                        }}
                      >
                        Title
                      </label>
                      <input
                        className="inp"
                        value={editTitle}
                        onChange={(e) => setEditTitle(e.target.value)}
                        style={{ fontSize: "10px" }}
                      />
                    </div>

                    <div style={{ marginBottom: "6px" }}>
                      <label
                        style={{
                          fontSize: "9px",
                          color: "#a09c94",
                          display: "block",
                          marginBottom: "2px",
                        }}
                      >
                        Description
                      </label>
                      <textarea
                        className="inp"
                        rows={2}
                        value={editDescription}
                        onChange={(e) => setEditDescription(e.target.value)}
                        style={{ resize: "none", fontSize: "10px" }}
                      />
                    </div>

                    <div style={{ marginBottom: "6px", display: "flex", gap: 8 }}>
                      <div style={{ flex: 1 }}>
                        <label
                          style={{
                            fontSize: "9px",
                            color: "#a09c94",
                            display: "block",
                            marginBottom: "2px",
                          }}
                        >
                          Logo
                        </label>
                        {(editLogoPreview || g.logo) && (
                          <div
                            style={{
                              width: "40px",
                              height: "40px",
                              borderRadius: "50%",
                              overflow: "hidden",
                              marginBottom: "4px",
                              border: "0.5px solid #3a3733",
                              background: editLogoPreview
                                ? `url(${editLogoPreview}) center/cover`
                                : `url(${g.logo}) center/cover`,
                            }}
                          />
                        )}
                        <input
                          ref={logoInputRef}
                          type="file"
                          accept="image/*"
                          style={{ display: "none" }}
                          onChange={(e) => {
                            const file = e.target.files?.[0];
                            if (file) {
                              setEditLogoFile(file);
                              setEditLogoPreview(URL.createObjectURL(file));
                            }
                          }}
                        />
                        <button
                          className="btn btn-g"
                          style={{ fontSize: "9px", padding: "3px 8px" }}
                          onClick={() => logoInputRef.current?.click()}
                        >
                          {editLogoPreview || g.logo ? "change" : "upload"}
                        </button>
                      </div>
                      <div style={{ flex: 1 }}>
                        <label
                          style={{
                            fontSize: "9px",
                            color: "#a09c94",
                            display: "block",
                            marginBottom: "2px",
                          }}
                        >
                          Background
                        </label>
                        {(editBgPreview || g.background) && (
                          <div
                            style={{
                              width: "100%",
                              height: "30px",
                              marginBottom: "4px",
                              border: "0.5px solid #3a3733",
                              background: editBgPreview
                                ? `url(${editBgPreview}) center/cover`
                                : `url(${g.background}) center/cover`,
                            }}
                          />
                        )}
                        <input
                          ref={bgInputRef}
                          type="file"
                          accept="image/*"
                          style={{ display: "none" }}
                          onChange={(e) => {
                            const file = e.target.files?.[0];
                            if (file) {
                              setEditBgFile(file);
                              setEditBgPreview(URL.createObjectURL(file));
                            }
                          }}
                        />
                        <button
                          className="btn btn-g"
                          style={{ fontSize: "9px", padding: "3px 8px" }}
                          onClick={() => bgInputRef.current?.click()}
                        >
                          {g.background ? "change" : "upload"}
                        </button>
                      </div>
                    </div>

                    <button
                      className="btn btn-red"
                      style={{ fontSize: "10px" }}
                      onClick={() => setEditingGroupId(null)}
                    >
                      cancel
                    </button>
                  </div>
                )}
              </div>
            ))}
          </div>
        )}
      </div>

      {/* ── INVITE MODAL ── */}
      {inviteGroupId !== null && (
        <div
          style={{
            position: "fixed",
            inset: 0,
            background: "rgba(0,0,0,0.6)",
            display: "flex",
            alignItems: "center",
            justifyContent: "center",
            zIndex: 100,
          }}
          onClick={() => { setInviteGroupId(null); setInviteResults([]); setInviteQuery(""); setInviteMsg(null); }}
        >
          <div
            className="card"
            style={{
              maxWidth: 380,
              width: "90%",
              maxHeight: "70vh",
              display: "flex",
              flexDirection: "column",
            }}
            onClick={(e) => e.stopPropagation()}
          >
            <div
              style={{
                display: "flex",
                justifyContent: "space-between",
                alignItems: "center",
                marginBottom: 10,
              }}
            >
              <p style={{ fontSize: 13, fontWeight: 600, color: "#e8e4dc" }}>
                Invite to group
              </p>
              <button
                style={{
                  background: "none",
                  border: "none",
                  color: "#6b6760",
                  cursor: "pointer",
                  fontSize: 16,
                }}
                onClick={() => { setInviteGroupId(null); setInviteResults([]); setInviteQuery(""); setInviteMsg(null); }}
              >
                <i className="ti ti-x" />
              </button>
            </div>

            <input
              className="inp"
              value={inviteQuery}
              onChange={(e) => void handleInviteSearch(e.target.value)}
              placeholder="Search users by name..."
              autoFocus
              style={{ marginBottom: 8, fontSize: 11 }}
            />

            {inviteMsg && (
              <p style={{ fontSize: 11, color: "#5cd4a0", marginBottom: 6 }}>{inviteMsg}</p>
            )}

            <div style={{ flex: 1, overflowY: "auto" }}>
              {inviteResults.length === 0 && inviteQuery.trim() && (
                <p style={{ fontSize: 11, color: "#6b6760", textAlign: "center", padding: 12 }}>
                  No users found
                </p>
              )}
              {inviteResults.map((u) => (
                <div
                  key={u.id}
                  style={{
                    display: "flex",
                    alignItems: "center",
                    gap: 8,
                    padding: "6px 0",
                    borderBottom: "0.5px solid #3a3733",
                  }}
                >
                  <div
                    className="av"
                    style={{
                      width: 26,
                      height: 26,
                      borderRadius: "50%",
                      background: u.avatar ? `url(${u.avatar}) center/cover` : "#2e1e24",
                      color: "#D4537E",
                      fontSize: 10,
                    }}
                  >
                    {!u.avatar && (u.nickname?.[0]?.toUpperCase() || u.firstname?.[0]?.toUpperCase() || "?")}
                  </div>
                  <div style={{ flex: 1, fontSize: 11 }}>
                    {u.nickname || `${u.firstname} ${u.lastname}`.trim()}
                  </div>
                  {memberIDs.includes(u.id) ? (
                    <span className="tag tag-gray" style={{ fontSize: 9 }}>member</span>
                  ) : (
                    <button
                      className="btn btn-t"
                      style={{ fontSize: 9 }}
                      disabled={inviteSending[u.id]}
                      onClick={() => void handleSendInvite(u.id)}
                    >
                      {inviteSending[u.id] ? "..." : "invite"}
                    </button>
                  )}
                </div>
              ))}
            </div>
          </div>
        </div>
      )}
    </main>
  );
}
