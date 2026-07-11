import Link from "next/link";

export default function Register() {
  return (
    <div style={{ padding:'2rem', display: 'flex', flexDirection: 'column', justifyContent: 'center' }} id="register-panel">
      <p className="sec-label" style={{ padding:0, marginBottom:'1rem' }}>create account</p>
      <div style={{ display: 'grid', gridTemplateColumns:'1fr 1fr', gap:'8px' }}>
        <div className="form-row"><span className="form-label">first name *</span><input className="inp" placeholder="Amir" /></div>
        <div className="form-row"><span className="form-label">last name *</span><input className="inp" placeholder="Kader" /></div>
      </div>
      <div className="form-row"><span className="form-label">email *</span><input className="inp" type="email" placeholder="you@example.com" /></div>
      <div className="form-row"><span className="form-label">password *</span><input className="inp" type="password" placeholder="••••••••" /></div>
      <div className="form-row"><span className="form-label">date of birth *</span><input className="inp" type="date" /></div>
      <div className="form-row"><span className="form-label">nickname <span className="tag tag-gray">optional</span></span><input className="inp" placeholder="@handle" /></div>
      <div className="form-row"><span className="form-label">about me <span className="tag tag-gray">optional</span></span><textarea className="inp" style={{ height:'50px', resize: 'none' }} placeholder="a few words..."></textarea></div>
      <div className="form-row">
        <span className="form-label">avatar <span className="tag tag-gray">optional</span></span>
        <div style={{ border:'0.5px dashed grey', padding:'10px', textAlign: 'center', fontSize:'11px', color:'#e8e4dc', cursor: 'pointer' }}>
          <i className="ti ti-upload" style={{ fontSize:'16px', display: 'block', marginBottom:'4px' }} aria-hidden="true"></i>
          jpeg / png / gif
        </div>
      </div>
      <button className="btn btn-p" style={{ width:'100%', marginTop:'6px' }} /* onClick={ showMainScreen() } */>create account →</button>
      <p style={{ fontSize:'11px', color:'var(--color-text-tertiary)', marginTop:'10px', textAlign: 'center' }}>already have an account? <Link style={{ color:'#D4537E', cursor: 'pointer' }} href="/login">login</Link></p>
    </div>
  );
}