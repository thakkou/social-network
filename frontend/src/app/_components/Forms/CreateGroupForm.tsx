"use client";

import { useState, useTransition, useEffect, useRef } from "react";
import Image from "next/image";
import { createGroup } from "~/app/_services/crud/groups";
import { search } from "~/app/_services/crud/search";

type InviteUser = {
  id: number;
  nickname: string;
  firstname: string;
  lastname: string;
  avatar: string;
}

export default function CreateGroupForm() {
  const [showCreateGroup, setShowCreateGroup] = useState(false);
  const [isPending, startTransition] = useTransition();

  // Form State
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [logo, setLogo] = useState<File | null>(null);
  const [background, setBackground] = useState<File | null>(null);

  // Image Preview States
  const [logoPreview, setLogoPreview] = useState<string | null>(null);
  const [backgroundPreview, setBackgroundPreview] = useState<string | null>(null);

  const [error, setError] = useState<string | null>(null);

  // ── Invite state ──
  const [inviteQuery, setInviteQuery] = useState("");
  const [searchResults, setSearchResults] = useState<InviteUser[]>([]);
  const [selectedUsers, setSelectedUsers] = useState<InviteUser[]>([]);
  const [searchOpen, setSearchOpen] = useState(false);
  const debounceRef = useRef<NodeJS.Timeout | null>(null);

  // Handle preview memory cleanup & updates
  useEffect(() => {
    if (!logo) {
      setLogoPreview(null);
      return;
    }
    const objectUrl = URL.createObjectURL(logo);
    setLogoPreview(objectUrl);
    return () => URL.revokeObjectURL(objectUrl);
  }, [logo]);

  useEffect(() => {
    if (!background) {
      setBackgroundPreview(null);
      return;
    }
    const objectUrl = URL.createObjectURL(background);
    setBackgroundPreview(objectUrl);
    return () => URL.revokeObjectURL(objectUrl);
  }, [background]);

  // ── Invite search with debounce ──
  useEffect(() => {
    if (debounceRef.current) clearTimeout(debounceRef.current);
    if (!inviteQuery.trim()) {
      setSearchResults([]);
      setSearchOpen(false);
      return;
    }
    debounceRef.current = setTimeout(async () => {
      const res = await search(inviteQuery.trim());
      if (res.success) {
        setSearchResults(
          res.data.profiles
            .filter((u: any) => !selectedUsers.some((s) => s.id === u.id))
            .map((u: any) => ({
              id: u.id,
              nickname: u.nickname || `${u.firstname} ${u.lastname}`.trim(),
              firstname: u.firstname,
              lastname: u.lastname,
              avatar: u.avatar || "",
            }))
        );
        setSearchOpen(true);
      }
    }, 400);
    return () => {
      if (debounceRef.current) clearTimeout(debounceRef.current);
    };
  }, [inviteQuery, selectedUsers]);

  const addUser = (user: InviteUser) => {
    setSelectedUsers((prev) => [...prev, user]);
    setInviteQuery("");
    setSearchResults([]);
    setSearchOpen(false);
  };

  const removeUser = (userId: number) => {
    setSelectedUsers((prev) => prev.filter((u) => u.id !== userId));
  };

  const handleReset = () => {
    setTitle("");
    setDescription("");
    setLogo(null);
    setBackground(null);
    setSelectedUsers([]);
    setInviteQuery("");
    setError(null);
    setShowCreateGroup(false);
  };

  const handleSubmit = () => {
    setError(null);

    if (!title.trim()) {
      setError("Group title is required.");
      return;
    }

    const formData = new FormData();
    formData.append("title", title.trim());
    formData.append("description", description.trim());

    if (logo) formData.append("logo", logo);
    if (background) formData.append("background", background);

    // Pass selected invite user IDs as a JSON string in the FormData
    if (selectedUsers.length > 0) {
      formData.append("invite_ids", JSON.stringify(selectedUsers.map((u) => u.id)));
    }

    startTransition(async () => {
      const res = await createGroup(formData);

      if (res.error) {
        setError(res.error);
      } else {
        handleReset();
        if (res.inviteErrors && res.inviteErrors.length > 0) {
          alert(
            `Group created! ${res.invitesSent?.length ?? 0} invited, ${res.inviteErrors.length} failed.\n\nFailed: ${res.inviteErrors.join(", ")}`
          );
        }
        window.location.reload();
      }
    });
  };

  return (
    <>
      <div
        style={{
          display: "flex",
          alignItems: "center",
          justifyContent: "space-between",
          marginBottom: "16px",
        }}
      >
        <p style={{ fontSize: "13px", fontWeight: 500, color: "var(--color-text-primary)" }}>
          browse all groups
        </p>
        <button
          className="btn btn-p"
          onClick={() => setShowCreateGroup(!showCreateGroup)}
          style={{ display: "flex", alignItems: "center", gap: "4px" }}
        >
          <i className="ti ti-plus" style={{ fontSize: "12px" }} aria-hidden="true" />
          create group
        </button>
      </div>

      {showCreateGroup && (
        <div id="create-group-panel" style={{ marginBottom: "16px" }}>
          <div className="card">
            <p
              style={{
                fontSize: "11px",
                fontWeight: 500,
                color: "var(--color-text-primary)",
                marginBottom: "10px",
              }}
            >
              new group
            </p>

            {error && (
              <div
                style={{
                  fontSize: "11px",
                  color: "var(--color-error, #e53e3e)",
                  marginBottom: "8px",
                }}
              >
                {error}
              </div>
            )}

            {/* Title */}
            <div className="form-row">
              <span className="form-label">title *</span>
              <input
                className="inp"
                placeholder="group name"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                disabled={isPending}
                required
              />
            </div>

            {/* Description */}
            <div className="form-row">
              <span className="form-label">description</span>
              <textarea
                className="inp"
                style={{ height: "56px", resize: "none" }}
                placeholder="what is this group about?"
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                disabled={isPending}
              />
            </div>

            {/* Logo Upload & Preview */}
            <div className="form-row">
              <span className="form-label">logo</span>
              <input
                type="file"
                className="inp"
                accept="image/*"
                onChange={(e) => setLogo(e.target.files?.[0] || null)}
                disabled={isPending}
              />
              {logoPreview && (
                <div style={{ marginTop: "8px", position: "relative", width: "fit-content" }}>
                  <Image
                    src={logoPreview}
                    alt="Logo preview"
                    width={48}
                    height={48}
                    style={{
                      objectFit: "cover",
                      borderRadius: "6px",
                      border: "1px solid var(--color-border, #ccc)",
                    }}
                  />
                  <button
                    type="button"
                    onClick={() => setLogo(null)}
                    style={{
                      position: "absolute",
                      top: "-6px",
                      right: "-6px",
                      background: "#e53e3e",
                      color: "#fff",
                      border: "none",
                      borderRadius: "50%",
                      width: "18px",
                      height: "18px",
                      fontSize: "10px",
                      cursor: "pointer",
                    }}
                  >
                    ×
                  </button>
                </div>
              )}
            </div>

            {/* Banner Upload & Preview */}
            <div className="form-row">
              <span className="form-label">banner</span>
              <input
                type="file"
                className="inp"
                accept="image/*"
                onChange={(e) => setBackground(e.target.files?.[0] || null)}
                disabled={isPending}
              />
              {backgroundPreview && (
                <div style={{ marginTop: "8px", position: "relative", width: "100%" }}>
                  <Image
                    src={backgroundPreview}
                    alt="Banner preview"
                    width={0}
                    height={80}
                    sizes="100vw"
                    style={{
                      width: "100%",
                      height: "80px",
                      objectFit: "cover",
                      borderRadius: "6px",
                      border: "1px solid var(--color-border, #ccc)",
                    }}
                    unoptimized
                  />
                  <button
                    type="button"
                    onClick={() => setBackground(null)}
                    style={{
                      position: "absolute",
                      top: "4px",
                      right: "4px",
                      background: "#e53e3e",
                      color: "#fff",
                      border: "none",
                      borderRadius: "50%",
                      width: "20px",
                      height: "20px",
                      fontSize: "12px",
                      cursor: "pointer",
                    }}
                  >
                    ×
                  </button>
                </div>
              )}
            </div>

            {/* ── Invite Members ── */}
            <div className="form-row">
              <span className="form-label">invite members</span>

              {/* Selected users chips */}
              {selectedUsers.length > 0 && (
                <div
                  style={{
                    display: "flex",
                    flexWrap: "wrap",
                    gap: 6,
                    marginBottom: 8,
                  }}
                >
                  {selectedUsers.map((u) => (
                    <span
                      key={u.id}
                      className="tag tag-pink"
                      style={{ fontSize: 10, cursor: "default", display: "flex", alignItems: "center", gap: 4 }}
                    >
                      {u.nickname || `${u.firstname} ${u.lastname}`.trim()}
                      <span
                        style={{ cursor: "pointer", marginLeft: 2 }}
                        onClick={() => removeUser(u.id)}
                      >
                        ×
                      </span>
                    </span>
                  ))}
                </div>
              )}

              <div style={{ position: "relative" }}>
                <input
                  className="inp"
                  placeholder="search by name..."
                  value={inviteQuery}
                  onChange={(e) => setInviteQuery(e.target.value)}
                  onFocus={() => searchResults.length > 0 && setSearchOpen(true)}
                />

                {searchOpen && searchResults.length > 0 && (
                  <div
                    className="card"
                    style={{
                      position: "absolute",
                      top: "calc(100% + 2px)",
                      left: 0,
                      right: 0,
                      zIndex: 10,
                      padding: "4px 0",
                      maxHeight: 180,
                      overflowY: "auto",
                    }}
                  >
                    {searchResults.map((u) => (
                      <div
                        key={u.id}
                        onClick={() => addUser(u)}
                        style={{
                          display: "flex",
                          alignItems: "center",
                          gap: 8,
                          padding: "6px 10px",
                          cursor: "pointer",
                          fontSize: 11,
                        }}
                        onMouseEnter={(e) => {
                          (e.currentTarget as HTMLElement).style.background = "#2e2b27";
                        }}
                        onMouseLeave={(e) => {
                          (e.currentTarget as HTMLElement).style.background = "transparent";
                        }}
                      >
                        <div
                          className="av"
                          style={{
                            width: 22,
                            height: 22,
                            borderRadius: "50%",
                            background: u.avatar ? `url(${u.avatar}) center/cover` : "#2e1e24",
                            color: "#D4537E",
                            fontSize: 9,
                          }}
                        >
                          {!u.avatar && (u.nickname?.[0]?.toUpperCase() || u.firstname?.[0]?.toUpperCase() || "?")}
                        </div>
                        <span>{u.nickname || `${u.firstname} ${u.lastname}`.trim()}</span>
                        <span style={{ marginLeft: "auto", color: "#D4537E", fontSize: 10 }}>+ add</span>
                      </div>
                    ))}
                  </div>
                )}
              </div>
            </div>

            {/* Submit & Cancel Buttons */}
            <div style={{ display: "flex", gap: "6px", marginTop: "8px" }}>
              <button className="btn btn-p" type="button" onClick={handleSubmit} disabled={isPending}>
                {isPending ? "creating..." : "create →"}
              </button>
              <button
                className="btn btn-g"
                type="button"
                onClick={handleReset}
                disabled={isPending}
              >
                cancel
              </button>
            </div>
          </div>
        </div>
      )}
    </>
  );
}