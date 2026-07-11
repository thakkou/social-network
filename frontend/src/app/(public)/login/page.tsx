import Link from "next/link";

export default function Login() {
  return (
    <div style={{ background:'var(--color-background-primary)', borderRight:'0.5px solid var(--color-border-tertiary)', padding:'2rem', display: 'flex', flexDirection: 'column', justifyContent: 'center' }}>
      <p className="sec-label" style={{ padding:0, marginBottom:'1rem' }}>sign in</p>
      <div className="form-row"><span className="form-label">email</span><input className="inp" type="email" placeholder="you@example.com" /></div>
      <div className="form-row"><span className="form-label">password</span><input className="inp" type="password" placeholder="••••••••" /></div>
      <button className="btn btn-p" style={{ width:'100%', marginTop:'8px' }} /* onClick={ showMainScreen() } */>sign in →</button>
      <p style={{ fontSize:'11px', color:'var(--color-text-tertiary)', marginTop:'10px', textAlign: 'center' }}>no account? <Link style={{ color:'#D4537E', cursor: 'pointer' }} href="/register">register</Link></p>
    </div>
  );
}