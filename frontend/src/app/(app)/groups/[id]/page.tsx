import Link from "next/link";

export default function GroupDetailPage() {
  return (
    <main className="main">
      <Link
        href="/groups"
        style={{
          display: "flex",
          alignItems: "center",
          gap: 6,
          color: "#a09c94",
          fontSize: 12,
          textDecoration: "none",
          width: "fit-content",
        }}
      >
        <i className="ti ti-arrow-left" />
        back to groups
      </Link>

      {/* HERO */}
      <div className="card" style={{ padding: 0, overflow: "hidden" }}>
        <div
          style={{
            height: 120,
            background:
              "linear-gradient(135deg,#2e2b27 0%, #272420 60%, #2e1e24 100%)",
            borderBottom: "1px solid #3a3733",
          }}
        />

        <div style={{ padding: "0 20px 20px" }}>
          <div
            className="av"
            style={{
              width: 72,
              height: 72,
              marginTop: -36,
              borderRadius: "50%",
              background: "#D4537E",
              fontSize: 26,
              color: "#fff",
              border: "4px solid #272420",
            }}
          >
            GO
          </div>

          <div
            style={{
              display: "flex",
              justifyContent: "space-between",
              alignItems: "flex-start",
              marginTop: 12,
              flexWrap: "wrap",
              gap: 16,
            }}
          >
            <div>
              <h2 style={{ fontSize: 22, marginBottom: 6 }}>go devs</h2>

              <p
                style={{
                  color: "#a09c94",
                  maxWidth: 650,
                  lineHeight: 1.6,
                  fontSize: 13,
                }}
              >
                A community dedicated to Go developers. Share backend projects,
                discuss architecture, learn Docker, databases, APIs and build
                cool things together.
              </p>

              <div
                style={{
                  marginTop: 12,
                  display: "flex",
                  gap: 8,
                  flexWrap: "wrap",
                }}
              >
                <span className="tag tag-pink">34 members</span>
                <span className="tag tag-gray">public</span>
                <span className="tag tag-teal">created Jan 2025</span>
              </div>
            </div>

            <div
              style={{
                display: "flex",
                gap: 8,
                flexWrap: "wrap",
              }}
            >
              <button className="btn btn-p">
                <i className="ti ti-message" /> chat
              </button>

              <button className="btn btn-g">
                <i className="ti ti-user-plus" /> invite
              </button>

              <button className="btn btn-red">
                <i className="ti ti-door-exit" /> leave
              </button>
            </div>
          </div>
        </div>
      </div>

      {/* NAVIGATION */}
      <div
        className="card"
        style={{
          display: "flex",
          gap: 10,
          padding: 8,
        }}
      >
        <button className="btn btn-p">all</button>
        <button className="btn btn-g">posts</button>
        <button className="btn btn-g">events</button>
      </div>

      {/* CREATE POST */}
      <div className="card">
        <p
          style={{
            fontSize: 12,
            color: "#a09c94",
            marginBottom: 12,
          }}
        >
          Create a new group post
        </p>

        <textarea
          className="inp"
          rows={5}
          placeholder="Share something with the group..."
          style={{ resize: "none" }}
        />

        <div
          style={{
            display: "flex",
            justifyContent: "space-between",
            marginTop: 12,
          }}
        >
          <button className="btn btn-g">
            <i className="ti ti-photo" /> image
          </button>

          <button className="btn btn-p">
            <i className="ti ti-send" /> publish
          </button>
        </div>
      </div>

      {/* EVENT */}
      <div className="event-card">
        <div
          style={{
            display: "flex",
            justifyContent: "space-between",
            marginBottom: 8,
          }}
        >
          <strong>Upcoming Event</strong>

          <span className="tag tag-amber">Fri 18:00</span>
        </div>

        <p style={{ marginBottom: 10, color: "#a09c94" }}>
          Docker Deep Dive — Multi-stage builds, compose and deployment
          pipelines.
        </p>

        <div style={{ display: "flex", gap: 8 }}>
          <button className="btn btn-t">
            <i className="ti ti-check" /> Going (11)
          </button>

          <button className="btn btn-g">Not Going (3)</button>
        </div>
      </div>

      {/* MEMBERS */}
      <div className="card">
        <div
          style={{
            display: "flex",
            justifyContent: "space-between",
            marginBottom: 16,
          }}
        >
          <strong>Members</strong>

          <Link
            href="#"
            style={{
              color: "#D4537E",
              fontSize: 12,
              textDecoration: "none",
            }}
          >
            view all →
          </Link>
        </div>

        {[
          "Amir",
          "Sarah",
          "Lucas",
          "Youssef",
        ].map((user) => (
          <div
            key={user}
            style={{
              display: "flex",
              alignItems: "center",
              gap: 12,
              padding: "8px 0",
            }}
          >
            <div
              className="av"
              style={{
                width: 34,
                height: 34,
                borderRadius: "50%",
                background: "#2e1e24",
                color: "#D4537E",
              }}
            >
              {user[0]}
            </div>

            <div style={{ flex: 1 }}>
              <div>{user}</div>

              <div
                style={{
                  color: "#6b6760",
                  fontSize: 11,
                }}
              >
                member
              </div>
            </div>
          </div>
        ))}
      </div>

      {/* POSTS */}
      {[1, 2].map((post) => (
        <div key={post} className="card">
          <div
            style={{
              display: "flex",
              alignItems: "center",
              gap: 10,
            }}
          >
            <div
              className="av"
              style={{
                width: 36,
                height: 36,
                borderRadius: "50%",
                background: "#2e1e24",
                color: "#D4537E",
              }}
            >
              A
            </div>

            <div>
              <strong>Amir</strong>

              <div
                style={{
                  color: "#6b6760",
                  fontSize: 11,
                }}
              >
                2 hours ago
              </div>
            </div>
          </div>

          <div className="divider" />

          <h3
            style={{
              marginBottom: 10,
              fontSize: 15,
            }}
          >
            New authentication package
          </h3>

          <p
            style={{
              color: "#a09c94",
              lineHeight: 1.7,
            }}
          >
            I've finished implementing session management and OAuth support.
            Feedback is welcome before merging into the main branch.
          </p>

          <div
            style={{
              display: "flex",
              gap: 8,
              marginTop: 16,
            }}
          >
            <button className="btn btn-g">
              <i className="ti ti-thumb-up" /> 12
            </button>

            <button className="btn btn-g">
              <i className="ti ti-message-circle" /> 4
            </button>
          </div>
        </div>
      ))}
    </main>
  );
}