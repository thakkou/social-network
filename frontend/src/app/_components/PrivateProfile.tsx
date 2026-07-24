"use client";

import { useState } from "react";
import Image from "next/image";

interface PrivateProfileProps {
    userId: string | number;
    profile: {
        firstname?: string;
        lastname?: string;
        nickname?: string;
        avatar?: string;
        following_status?: string;
    };
    onSendRequest?: () => void;
}

export default function PrivateProfile({
    userId,
    profile,
    onSendRequest,
}: PrivateProfileProps) {
const isPending = profile.following_status === "pending";

const [isLoading, setIsLoading] = useState(false);
const buttonLabel = isLoading
  ? "..."
  : isPending
  ? "Follow request sent"
  : "Send follow request";




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
                        position: "relative",
                    }}
                >
                    {profile.avatar ? (
                        <Image
                            src={profile.avatar}
                            alt={profile.nickname ?? "avatar"}
                            fill
                            style={{ objectFit: "cover" }}
                        />
                    ) : (
                        `${profile?.firstname?.[0] ?? ""}${profile?.lastname?.[0] ?? ""}`
                    )}
                </div>

                <div>
                    <p style={{ fontSize: "14px", fontWeight: 500 }}>
                        {profile.firstname} {profile.lastname}
                    </p>

                    {profile.nickname && (
                        <span className="tag tag-teal">
                            @{profile.nickname}
                        </span>
                    )}
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
  onClick={() => onSendRequest?.()}
  disabled={isLoading}
>
  {buttonLabel}
</button>
            </div>
        </div>
    );
}