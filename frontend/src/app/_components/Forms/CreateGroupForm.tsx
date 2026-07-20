"use client";

import { useState, useTransition, useEffect } from "react";
import SearchInput from "~/app/_components/SearchInput";
import { createGroup } from "~/app/api/crud/groups"; // Adjust import path as needed

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

  const handleReset = () => {
    setTitle("");
    setDescription("");
    setLogo(null);
    setBackground(null);
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

    if (logo) {
      formData.append("logo", logo);
    }
    if (background) {
      formData.append("background", background);
    }

    startTransition(async () => {
      const res = await createGroup(formData);

      if (res.error) {
        setError(res.error);
      } else {
        handleReset();
        // Optional: Trigger a refresh or redirect to the new group
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
          {/* Changed from <form> to <div> to avoid nesting <form> elements from SearchInput */}
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
                  {/* eslint-disable-next-line @next/next/no-img-element */}
                  <img
                    src={logoPreview}
                    alt="Logo preview"
                    style={{
                      width: "48px",
                      height: "48px",
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
                  {/* eslint-disable-next-line @next/next/no-img-element */}
                  <img
                    src={backgroundPreview}
                    alt="Banner preview"
                    style={{
                      width: "100%",
                      height: "80px",
                      objectFit: "cover",
                      borderRadius: "6px",
                      border: "1px solid var(--color-border, #ccc)",
                    }}
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

            {/* Invite Members */}
            <div className="form-row">
              <span className="form-label">invite members</span>
              <SearchInput className="inp" placeholder="search by name..." typeSearch="users" />
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