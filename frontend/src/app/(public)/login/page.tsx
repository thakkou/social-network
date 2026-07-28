"use client";

import { signIn } from "next-auth/react";
import { useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";

const devMode = true;

const testUsers = [
  { label: "alice", identifier: "alice@example.com", password: "password" },
  { label: "bob", identifier: "bob@example.com", password: "password" },
  { label: "chloe", identifier: "chloe@example.com", password: "password" },
  { label: "farid", identifier: "farid@example.com", password: "password" },
  { label: "isabella", identifier: "isabella@example.com", password: "password" },
  { label: "jack", identifier: "jack@example.com", password: "password" },
];

export default function Login() {
  const [identifier, setIdentifier] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const router = useRouter();

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const res = await signIn("credentials", {
      identifier,
      password,
      redirect: false,
    });

    if (res?.error) {
      setError("Invalid email/nickname or password");
    } else {
      router.push("/");
    }
  };

  const fillTestUser = (user: (typeof testUsers)[number]) => {
    setIdentifier(user.identifier);
    setPassword(user.password);
    setError("");
  };

  return (
    <form
      onSubmit={handleSubmit}
      style={{
        background: "#211f1c",
        borderRight: "0.5px solid #3a3733",
        padding: "2.5rem",
        display: "flex",
        flexDirection: "column",
        justifyContent: "center",
      }}
      className="auth-form"
    >
      {/* Brand */}
      <div style={{ marginBottom: "1.5rem" }}>
        <p className="logo" style={{ fontSize: 18, marginBottom: 4 }}>
          social<span>network</span>
        </p>
        <p style={{ fontSize: 11, color: "#6b6760" }}>sign in to your account</p>
      </div>

      {error && (
        <div
          style={{
            fontSize: 11,
            color: "#e07070",
            background: "#2a1818",
            border: "0.5px solid #7a2c2c",
            padding: "8px 10px",
            marginBottom: 12,
          }}
        >
          {error}
        </div>
      )}

      <div className="form-row">
        <span className="form-label">email / nickname</span>
        <input
          className="inp"
          type="text"
          value={identifier}
          onChange={(e) => setIdentifier(e.target.value)}
          placeholder="you@example.com"
        />
      </div>

      <div className="form-row">
        <span className="form-label">password</span>
        <input
          className="inp"
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          placeholder="••••••••"
        />
      </div>

      <button
        className="btn btn-p"
        style={{ width: "100%", marginTop: "8px", padding: "8px 12px", fontSize: 12 }}
        type="submit"
      >
        sign in →
      </button>

      <p
        style={{
          fontSize: 11,
          color: "#6b6760",
          marginTop: 12,
          textAlign: "center",
        }}
      >
        no account?{" "}
        <Link
          style={{ color: "#D4537E", textDecoration: "none" }}
          href="/register"
        >
          create one
        </Link>
      </p>

      {devMode && (
        <div
          style={{
            marginTop: "20px",
            paddingTop: "16px",
            borderTop: "0.5px dashed #3a3733",
          }}
        >
          <p
            style={{
              fontSize: 10,
              color: "#6b6760",
              textAlign: "center",
              marginBottom: 8,
              textTransform: "uppercase",
              letterSpacing: "0.05em",
            }}
          >
            quick login
          </p>
          <div
            style={{
              display: "flex",
              flexWrap: "wrap",
              gap: 6,
              justifyContent: "center",
            }}
          >
            {testUsers.map((user, _i) => (
              <button
                key={user.identifier}
                type="button"
                className="btn btn-g"
                style={{ fontSize: 10, padding: "4px 10px" }}
                onClick={() => fillTestUser(user)}
              >
                {user.label}
              </button>
            ))}
          </div>
        </div>
      )}
    </form>
  );
}