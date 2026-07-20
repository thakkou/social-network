"use client";

import { useState } from "react";
import SearchInput from "~/app/_components/SearchInput";

export default function CreateGroupForm() {
  const [showCreateGroup, setShowCreateGroup] = useState(false);

  return (
    <>
      <div style={{ display: "flex", alignItems: "center", justifyContent: "space-between", marginBottom: "16px" }}>
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
            <p style={{ fontSize: "11px", fontWeight: 500, color: "var(--color-text-primary)", marginBottom: "10px" }}>
              new group
            </p>
            <div className="form-row">
              <span className="form-label">title *</span>
              <input className="inp" placeholder="group name" />
            </div>
            <div className="form-row">
              <span className="form-label">description *</span>
              <textarea
                className="inp"
                style={{ height: "56px", resize: "none" }}
                placeholder="what is this group about?"
              />
            </div>
            <div className="form-row">
              <span className="form-label">invite members</span>
              <SearchInput className="inp" placeholder="search by name..." typeSearch="users"/>
            </div>
            <div style={{ display: "flex", gap: "6px", marginTop: "4px" }}>
              <button className="btn btn-p">create →</button>
              <button className="btn btn-g" onClick={() => setShowCreateGroup(false)}>
                cancel
              </button>
            </div>
          </div>
        </div>
      )}
    </>
  );
}