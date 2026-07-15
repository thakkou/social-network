"use client";

import { signIn } from "next-auth/react";
import { useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";

export default function Login() {
  const [identifier, setIdentifier] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const router = useRouter();

  const handleSubmit = async (e: React.SubmitEvent) => {
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

  return (
    <form style={{ background:'var(--color-background-primary)', borderRight:'0.5px solid var(--color-border-tertiary)', padding:'2rem', display: 'flex', flexDirection: 'column', justifyContent: 'center' }} onSubmit={handleSubmit}>
      <p className="sec-label" style={{ padding:0, marginBottom:'1rem' }}>sign in</p>
      {error && <p>{error}</p>}
      <div className="form-row"><span className="form-label">email / nickname</span><input className="inp" type="text" value={identifier} onChange={(e) => setIdentifier(e.target.value)} placeholder="email/username" /></div>
      <div className="form-row"><span className="form-label">password</span><input className="inp" type="password" value={password} onChange={(e) => setPassword(e.target.value)} placeholder="••••••••" /></div>
      <button className="btn btn-p" style={{ width:'100%', marginTop:'8px' }} type="submit">sign in →</button>
      <p style={{ fontSize:'11px', color:'var(--color-text-tertiary)', marginTop:'10px', textAlign: 'center' }}>no account? <Link style={{ color:'#D4537E', cursor: 'pointer' }} href="/register">register</Link></p>
    </form>
  );
}