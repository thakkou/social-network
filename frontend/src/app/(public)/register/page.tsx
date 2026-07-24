"use client";

import { useState, useTransition } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import Image from "next/image";

import { registerUser } from "./actions";

export default function Register() {
  const [firstname, setFirstname] = useState("");
  const [lastname, setLastname] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [birthDate, setBirthDate] = useState("");
  const [nickname, setNickname] = useState("");
  const [aboutme, setAboutme] = useState("");

  const [error, setError] = useState("");
  const [preview, setPreview] = useState<string | null>(null);
  const [isPending, startTransition] = useTransition();
  const router = useRouter();

  const handleAvatarChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return setPreview(null);

    if (file.size > 1 * 1024 * 1024) {
      setError("Image must be under 1MB");
      e.target.value = "";
      return;
    }

    setPreview(URL.createObjectURL(file));
  };

  const handleSubmit = async (formData: FormData) => {
    startTransition(async () => {
      const result = await registerUser(formData);

      if ("error" in result) {
        setError(result.error);
      } else {
        router.push("/");
        router.refresh();
      }
    });
  };

  return (
    <form
      action={handleSubmit}
      style={{
        background: "#211f1c",
        padding: "2.5rem",
        display: "flex",
        flexDirection: "column",
        justifyContent: "center",
        overflowY: "auto",
      }}
    >
      {/* Brand */}
      <div style={{ marginBottom: "1.25rem" }}>
        <p className="logo" style={{ fontSize: 18, marginBottom: 4 }}>
          social<span>network</span>
        </p>
        <p style={{ fontSize: 11, color: "#6b6760" }}>create your account</p>
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

      <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: "8px" }}>
        <div className="form-row" style={{ marginBottom: 0 }}>
          <span className="form-label">first name *</span>
          <input
            className="inp"
            name="firstname"
            type="text"
            value={firstname}
            onChange={(e) => setFirstname(e.target.value)}
            placeholder="Amir"
            style={{ fontSize: 12 }}
          />
        </div>
        <div className="form-row" style={{ marginBottom: 0 }}>
          <span className="form-label">last name *</span>
          <input
            className="inp"
            name="lastname"
            type="text"
            value={lastname}
            onChange={(e) => setLastname(e.target.value)}
            placeholder="Kader"
            style={{ fontSize: 12 }}
          />
        </div>
      </div>

      <div className="form-row">
        <span className="form-label">email *</span>
        <input
          className="inp"
          name="email"
          type="email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          placeholder="you@example.com"
          style={{ fontSize: 12 }}
        />
      </div>

      <div className="form-row">
        <span className="form-label">password *</span>
        <input
          className="inp"
          name="password"
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          placeholder="••••••••"
          style={{ fontSize: 12 }}
        />
      </div>

      <div className="form-row">
        <span className="form-label">date of birth *</span>
        <input
          className="inp"
          name="birthDate"
          type="date"
          value={birthDate}
          onChange={(e) => setBirthDate(e.target.value)}
          style={{ fontSize: 12 }}
        />
      </div>

      <div className="form-row">
        <span className="form-label">
          nickname{" "}
          <span
            className="tag tag-gray"
            style={{ fontSize: 9, verticalAlign: "middle" }}
          >
            optional
          </span>
        </span>
        <input
          className="inp"
          name="nickname"
          type="text"
          value={nickname}
          onChange={(e) => setNickname(e.target.value)}
          placeholder="@handle"
          style={{ fontSize: 12 }}
        />
      </div>

      <div className="form-row">
        <span className="form-label">
          about me{" "}
          <span
            className="tag tag-gray"
            style={{ fontSize: 9, verticalAlign: "middle" }}
          >
            optional
          </span>
        </span>
        <textarea
          className="inp"
          name="aboutme"
          value={aboutme}
          onChange={(e) => setAboutme(e.target.value)}
          style={{ height: "44px", resize: "none", fontSize: 12 }}
          placeholder="a few words..."
        />
      </div>

      <div className="form-row">
        <span className="form-label">
          avatar{" "}
          <span
            className="tag tag-gray"
            style={{ fontSize: 9, verticalAlign: "middle" }}
          >
            optional
          </span>
        </span>
        {preview && (
          <div
            style={{
              width: 56,
              height: 56,
              borderRadius: "50%",
              overflow: "hidden",
              marginBottom: 6,
            }}
          >
            <Image
              src={preview}
              alt="Avatar preview"
              width={56}
              height={56}
              style={{ objectFit: "cover" }}
            />
          </div>
        )}
        <input id="avatar" name="avatar" type="file" accept=".jpg,.jpeg,.png,.gif" onChange={handleAvatarChange} hidden />
        <label
          htmlFor="avatar"
          style={{
            border: "0.5px dashed #3a3733",
            padding: "10px",
            textAlign: "center",
            fontSize: 11,
            color: "#a09c94",
            cursor: "pointer",
            display: "block",
          }}
        >
          <i
            className="ti ti-upload"
            style={{ fontSize: 14, display: "block", marginBottom: 2 }}
          />
          jpeg / png / gif
        </label>
      </div>

      <button
        className="btn btn-p"
        style={{ width: "100%", marginTop: "4px", padding: "8px 12px", fontSize: 12 }}
        type="submit"
        disabled={isPending}
      >
        {isPending ? "creating account..." : "create account →"}
      </button>

      <p
        style={{
          fontSize: 11,
          color: "#6b6760",
          marginTop: 12,
          textAlign: "center",
        }}
      >
        already have an account?{" "}
        <Link
          style={{ color: "#D4537E", textDecoration: "none" }}
          href="/login"
        >
          sign in
        </Link>
      </p>
    </form>
  );
}