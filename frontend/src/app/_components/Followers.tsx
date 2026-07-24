import Image from "next/image";
import Link from "next/link";

type FollowerUser = {
  id: number;
  nickname?: string;
  firstname?: string;
  lastname?: string;
  avatar?: string;
}

export default function Followers({ users = [] }: { users?: FollowerUser[] }) {
    return (
        <div id="profile-followers">
            <div className="card">
                <div style={{ display: "flex", flexDirection: "column", gap: "10px" }}>
                    {users.map((u) => {
                        const name =
                            u.nickname ||
                            `${u.firstname || ""} ${u.lastname || ""}`.trim();

                        const initials = name
                            .split(" ")
                            .map((n) => n[0])
                            .join("")
                            .slice(0, 2)
                            .toUpperCase();

                        return (
                            <div
                                key={u.id}
                                style={{
                                    display: "flex",
                                    alignItems: "center",
                                    gap: "8px",
                                }}
                            >
                                {u.avatar ? (
                                    <Image
                                        src={u.avatar}
                                        alt={name}
                                        width={28}
                                        height={28}
                                        className="av"
                                        style={{
                                            borderRadius: "50%",
                                            objectFit: "cover",
                                        }}
                                    />
                                ) : (
                                    <div
                                        className="av"
                                        style={{
                                            width: "28px",
                                            height: "28px",
                                            background: "#FBEAF0",
                                            color: "#993556",
                                            fontSize: "11px",
                                            display: "flex",
                                            alignItems: "center",
                                            justifyContent: "center",
                                            borderRadius: "50%",
                                        }}
                                    >
                                        {initials}
                                    </div>
                                )}

                                <span
                                    style={{
                                        fontSize: "12px",
                                        fontWeight: 500,
                                        color: "var(--color-text-primary)",
                                    }}
                                >
                                    {name}
                                </span>

                                <Link
                                    href={`/profile/${u.id}`}
                                    className="btn btn-g"
                                    style={{
                                        marginLeft: "auto",
                                        fontSize: "10px",
                                    }}
                                >
                                    view profile
                                </Link>
                            </div>
                        );
                    })}
                </div>
            </div>
        </div>
    );
}