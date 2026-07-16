"use client";

import { signIn } from "next-auth/react";
import { useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";

const devMode = true;

const testUsers = [
  { label: "user1 (alice)", identifier: "alice@example.com", password: "Password123!" },
  { label: "user2 (bob)", identifier: "bob@example.com", password: "Password123!" },
  { label: "user3 (chloe)", identifier: "chloe@example.com", password: "Password123!" },
  { label: "user4 (farid)", identifier: "farid@example.com", password: "Password123!" },
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
    <form style={{ background:'var(--color-background-primary)', borderRight:'0.5px solid var(--color-border-tertiary)', padding:'2rem', display: 'flex', flexDirection: 'column', justifyContent: 'center' }} onSubmit={handleSubmit}>
      <p className="sec-label" style={{ padding:0, marginBottom:'1rem' }}>sign in</p>
      {error && <p>{error}</p>}
      <div className="form-row"><span className="form-label">email / nickname</span><input className="inp" type="text" value={identifier} onChange={(e) => setIdentifier(e.target.value)} placeholder="email/username" /></div>
      <div className="form-row"><span className="form-label">password</span><input className="inp" type="password" value={password} onChange={(e) => setPassword(e.target.value)} placeholder="••••••••" /></div>
      <button className="btn btn-p" style={{ width:'100%', marginTop:'8px' }} type="submit">sign in →</button>
      <p style={{ fontSize:'11px', color:'var(--color-text-tertiary)', marginTop:'10px', textAlign: 'center' }}>no account? <Link style={{ color:'#D4537E', cursor: 'pointer' }} href="/register">register</Link></p>

      {devMode && (
        <div style={{ marginTop: '16px', paddingTop: '16px', borderTop: '0.5px dashed var(--color-border-tertiary)', display: 'flex', flexDirection: 'column', gap: '6px' }}>
          <span style={{ fontSize: '11px', color: 'var(--color-text-tertiary)', textAlign: 'center' }}>dev quick-fill</span>
          <div style={{ display: 'flex', flexWrap: 'wrap', gap: '6px', justifyContent: 'center' }}>
            {testUsers.map((user, i) => (
              <button
                key={user.identifier}
                type="button"
                className="btn"
                style={{ fontSize: '11px', padding: '4px 10px' }}
                onClick={() => fillTestUser(user)}
              >
                user{i + 1}
              </button>
            ))}
          </div>
        </div>
      )}
    </form>
  );
}