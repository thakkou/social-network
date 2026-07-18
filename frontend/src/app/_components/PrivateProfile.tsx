"use client";

interface PrivateProfileProps {
    profile: {
        first_name?: string;
        last_name?: string;
        nickname?: string;
        avatar?: string;
    };
    onSendRequest?: () => void;
}

export default function PrivateProfile({
    profile,
    onSendRequest,
}: PrivateProfileProps) {
    return (
        <div className="card">
            <div style={{ display: "flex", alignItems: "center", gap: "12px" }}>
                <div
                    className="av"
                    style={{
                        width: "52px",
                        height: "52px",
                        background: "#EEEDFE",
                        color: "#534AB7",
                        fontSize: "16px",
                        overflow: "hidden",
                    }}
                >
                    {profile.avatar ? (
                        <img
                            src={profile.avatar}
                            alt={profile.nickname}
                            style={{
                                width: "100%",
                                height: "100%",
                                objectFit: "cover",
                            }}
                        />
                    ) : (
                        `${profile.first_name?.[0] ?? ""}${profile.last_name?.[0] ?? ""}`
                    )}
                </div>

                <div>
                    <p style={{ fontSize: "14px", fontWeight: 500 }}>
                        {profile.first_name} {profile.last_name}
                    </p>

                    <span className="tag tag-teal">
                        @{profile.nickname}
                    </span>
                </div>
            </div>

            <div className="divider" />

            <div style={{ textAlign: "center", padding: "20px" }}>
                <p style={{ fontSize: "14px", marginBottom: "10px" }}>
                    This is a private profile
                </p>

                <p
                    style={{
                        fontSize: "12px",
                        color: "var(--color-text-secondary)",
                        marginBottom: "16px",
                    }}
                >
                    Send a follow request to see posts and activity.
                </p>

                <button
                    className="btn btn-g"
                    onClick={onSendRequest}
                >
                    Send request
                </button>
            </div>
        </div>
    );
}