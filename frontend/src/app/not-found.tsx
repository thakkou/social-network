import Link from "next/link";

export default function NotFound() {
  return (
    <main className="main" style={{ minHeight: "100vh", display: "flex", alignItems: "center", justifyContent: "center" }}>
      <div className="card" style={{ textAlign: "center", maxWidth: 420, padding: "48px 32px" }}>
        <div
          style={{
            width: 72,
            height: 72,
            borderRadius: "50%",
            background: "#2e1e24",
            color: "#D4537E",
            display: "flex",
            alignItems: "center",
            justifyContent: "center",
            margin: "0 auto 20px",
            fontSize: 32,
          }}
        >
          <i className="ti ti-compass-off" />
        </div>

        <h1 style={{ fontSize: 28, fontWeight: 700, color: "#e8e4dc", marginBottom: 8 }}>
          404
        </h1>

        <p style={{ fontSize: 15, color: "#a09c94", lineHeight: 1.6, marginBottom: 24 }}>
          hi i think you are lost come back
        </p>

        <Link
          href="/"
          className="btn btn-p"
          style={{
            display: "inline-flex",
            alignItems: "center",
            gap: 6,
            fontSize: 13,
            textDecoration: "none",
          }}
        >
          <i className="ti ti-arrow-left" style={{ fontSize: 14 }} />
          go home
        </Link>
      </div>
    </main>
  );
}
