export default function PublicLayout({ children }: Readonly<{ children: React.ReactNode }>) {
  return (
    <div className="screen active">
      <div className="auth-layout">
        {children}

        {/* Decorative right panel */}
        <div className="auth-decorative">
          {/* Decorative circles */}
          <div className="auth-circle auth-circle-1" />
          <div className="auth-circle auth-circle-2" />
          <div className="auth-circle auth-circle-3" />

          <div style={{ textAlign: "center", position: "relative" }}>
            <p
              className="logo"
              style={{
                fontSize: 28,
                color: "#D4537E",
                marginBottom: 12,
                letterSpacing: "-0.03em",
              }}
            >
              social<span style={{ color: "#e8e4dc" }}>network</span>
            </p>
            <p
              style={{
                fontSize: 12,
                color: "#6b6760",
                lineHeight: 1.7,
                maxWidth: 280,
                margin: "0 auto",
              }}
            >
              Connect with friends, share your thoughts, and discover communities
              that matter to you.
            </p>

            <div
              style={{
                marginTop: 24,
                display: "flex",
                gap: 8,
                justifyContent: "center",
              }}
            >
              <span className="tag tag-pink">posts</span>
              <span className="tag tag-teal">groups</span>
              <span className="tag tag-purple">events</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}