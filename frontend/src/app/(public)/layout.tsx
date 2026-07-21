export default function PublicLayout({ children }: Readonly<{ children: React.ReactNode }>) {
  return (
    <div className="screen active">
      <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', minHeight: '100vh' }}>
        {children}

        {/* Decorative right panel */}
        <div
          style={{
            background:
              "linear-gradient(135deg, #2e1e24 0%, #1a1917 50%, #162820 100%)",
            display: "flex",
            flexDirection: "column",
            alignItems: "center",
            justifyContent: "center",
            padding: "2rem",
            position: "relative",
            overflow: "hidden",
          }}
        >
          {/* Decorative circles */}
          <div
            style={{
              position: "absolute",
              top: "-60px",
              right: "-60px",
              width: "200px",
              height: "200px",
              borderRadius: "50%",
              background: "radial-gradient(circle, rgba(212,83,126,0.08) 0%, transparent 70%)",
            }}
          />
          <div
            style={{
              position: "absolute",
              bottom: "-40px",
              left: "-40px",
              width: "160px",
              height: "160px",
              borderRadius: "50%",
              background: "radial-gradient(circle, rgba(29,158,117,0.06) 0%, transparent 70%)",
            }}
          />
          <div
            style={{
              position: "absolute",
              top: "40%",
              left: "20%",
              width: "80px",
              height: "80px",
              borderRadius: "50%",
              background: "radial-gradient(circle, rgba(212,83,126,0.04) 0%, transparent 60%)",
            }}
          />

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