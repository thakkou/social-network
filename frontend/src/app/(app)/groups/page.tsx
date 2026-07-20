import Link from "next/link";
import SearchInput from "~/app/_components/SearchInput";
import CreateGroupForm from "~/app/_components/Forms/CreateGroupForm"; // Adjust import path as needed

export default function Groups() {
  return (
    <main className="main">
      {/* Client Component handling the interactive toggle & form */}
      <CreateGroupForm />

      <div style={{ marginBottom: "16px" }}>
        <SearchInput placeholder="Search groups..." typeSearch="groups" />
      </div>

      {/* Member Group Card */}
      <div className="card">
        <div style={{ display: "flex", alignItems: "center", gap: "8px", marginBottom: "8px" }}>
          <span style={{ width: "8px", height: "8px", background: "#D4537E", flexShrink: 0 }} />
          <p style={{ fontSize: "13px", fontWeight: 500, color: "var(--color-text-primary)" }}>go devs</p>
          <span className="tag tag-pink" style={{ marginLeft: "auto" }}>member</span>
        </div>
        <p style={{ fontSize: "12px", color: "var(--color-text-secondary)", marginBottom: "8px" }}>
          A space for Go enthusiasts — sharing projects, tips and co-building the social network backend.
        </p>
        <div style={{ fontSize: "11px", color: "var(--color-text-tertiary)", marginBottom: "10px" }}>
          34 members · created by Amir K.
        </div>

        <div className="event-card" style={{ marginBottom: "8px" }}>
          <div style={{ display: "flex", justifyContent: "space-between", marginBottom: "4px" }}>
            <p style={{ fontSize: "12px", fontWeight: 500, color: "var(--color-text-primary)" }}>Docker deep dive</p>
            <span className="tag tag-amber">Fri 18:00</span>
          </div>
          <p style={{ fontSize: "11px", color: "var(--color-text-secondary)", marginBottom: "6px" }}>
            Multi-stage builds, compose, and deployment pipelines.
          </p>
          <div style={{ display: "flex", gap: "6px" }}>
            <button className="btn btn-t" style={{ fontSize: "10px", display: "flex", alignItems: "center", gap: "3px" }}>
              <i className="ti ti-check" style={{ fontSize: "11px" }} aria-hidden="true" /> going (11)
            </button>
            <button className="btn btn-g" style={{ fontSize: "10px" }}>not going (3)</button>
          </div>
        </div>

        <div style={{ display: "flex", gap: "6px" }}>
          <Link className="btn btn-p" href="/messages" style={{ fontSize: "11px", display: "flex", alignItems: "center", gap: "3px" }}>
            <i className="ti ti-message" style={{ fontSize: "12px" }} aria-hidden="true" /> group chat
          </Link>
          <button className="btn btn-g" style={{ fontSize: "11px", display: "flex", alignItems: "center", gap: "3px" }}>
            <i className="ti ti-calendar-plus" style={{ fontSize: "12px" }} aria-hidden="true" /> add event
          </button>
          <button className="btn btn-g" style={{ fontSize: "11px", display: "flex", alignItems: "center", gap: "3px" }}>
            <i className="ti ti-user-plus" style={{ fontSize: "12px" }} aria-hidden="true" /> invite
          </button>
        </div>
      </div>

      {/* Non-Member Group Card */}
      <div className="card">
        <div style={{ display: "flex", alignItems: "center", gap: "8px", marginBottom: "8px" }}>
          <span style={{ width: "8px", height: "8px", background: "#534AB7", flexShrink: 0 }} />
          <p style={{ fontSize: "13px", fontWeight: 500, color: "var(--color-text-primary)" }}>design systems</p>
          <span className="tag tag-gray" style={{ marginLeft: "auto" }}>not a member</span>
        </div>
        <p style={{ fontSize: "12px", color: "var(--color-text-secondary)", marginBottom: "8px" }}>
          Discussing component libraries, tokens, and consistent UI patterns across projects.
        </p>
        <div style={{ fontSize: "11px", color: "var(--color-text-tertiary)", marginBottom: "10px" }}>
          21 members · created by Selin R.
        </div>
        <button className="btn btn-p" style={{ fontSize: "11px" }}>request to join →</button>
      </div>
    </main>
  );
}