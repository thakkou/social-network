export default function Home() {
  return (
    <main className="main">
      <div className="card">
        <div style={{ display: 'flex', gap:'8px', alignItems:'flex-start' }}>
          <div className="av" style={{ width:'32px', height:'32px', background:'#EEEDFE', color:'#534AB7', fontSize:'12px' }}>AK</div>
          <div style={{ flex:1 }}>
            <textarea className="inp" style={{ resize: 'none', height:'56px' }} placeholder="what's on your mind?"></textarea>
            <div className="post-privacy-row">
              <button className="btn btn-g" style={{ display: 'flex', alignItems: 'center', gap:'4px', fontSize:'11px' }}><i className="ti ti-photo" style={{ fontSize:'13px' }} aria-hidden="true"></i> image/gif</button>
              <select className="inp" style={{ width: 'auto', padding:'4px 8px', fontSize:'11px' }}>
                <option>public — everyone</option>
                <option>almost private — followers</option>
                <option>private — select followers</option>
              </select>
              <button className="btn btn-p" style={{ marginLeft: 'auto' }}>post →</button>
            </div>
          </div>
        </div>
      </div>

      <div className="card">
        <div style={{ display: 'flex', alignItems: 'center', gap:'8px', marginBottom:'8px' }}>
          <div className="av" style={{ width:'30px', height:'30px', background:'#FBEAF0', color:'#993556', fontSize:'11px' }}>SR</div>
          <div style={{ flex:1 }}>
            <p style={{ fontSize:'12px', fontWeight:500, color:'var(--color-text-primary)' }}>Selin Rauf <span className="tag tag-teal">public</span></p>
            <p style={{ fontSize:'10px', color:'var(--color-text-tertiary)' }}>2h ago</p>
          </div>
        </div>
        <p style={{ fontSize:'12px', color:'var(--color-text-primary)', lineHeight:1.6, marginBottom:'8px' }}>Just merged the websocket handler for real-time group broadcasts. Notifications now push instantly — test it out and drop feedback below.</p>
        <div  className="code-card">
        {/* <div style={{background: 'var(--color-background-secondary)', border: '0.5px solid var(--color-border-tertiary)', padding: '10px 14px', marginBottom: '10px', fontSize: '12px', color: 'var(--color-text-secondary)', fontFamily: 'var(--font-mono, monospace)'}}> */}
          <span style={{ color:'#D4537E' }}>// feat:</span> websocket broadcast + group chat rooms<br/>
          <span style={{ color:'#1D9E75' }}>+247</span> <span style={{ color:'#D85A30' }}>-12</span>
        </div>
        <div className="divider"></div>
        <div style={{ display: 'flex', gap:'6px', alignItems: 'center' }}>
          <button className="btn btn-g" style={{ display: 'flex', alignItems: 'center', gap:'3px' }}><i className="ti ti-heart" style={{ fontSize:'12px' }} aria-hidden="true"></i> 24</button>
          <button className="btn btn-g" style={{ display: 'flex', alignItems: 'center', gap:'3px' }}><i className="ti ti-message-circle" style={{ fontSize:'12px' }} aria-hidden="true"></i> comment</button>
          <span style={{ fontSize:'10px', color:'var(--color-text-tertiary)', marginLeft: 'auto' }}>seen by 14 followers</span>
        </div>
        <div style={{ marginTop:'8px', paddingTop:'8px', borderTop:'0.5px solid var(--color-border-tertiary)' }}>
          <div style={{ display: 'flex', gap:'8px', alignItems:'flex-start', marginBottom:'6px' }}>
            <div className="av" style={{ width:'22px', height:'22px', background:'#E1F5EE', color:'#0F6E56', fontSize:'10px' }}>JM</div>
            <div style={{ background:'var(--color-background-secondary)', border:'0.5px solid var(--color-border-tertiary)', padding:'5px 8px', fontSize:'11px', flex:1, color:'var(--color-text-primary)' }}>great work, tested it and works perfectly!</div>
          </div>
          <div style={{ display: 'flex', gap:'6px', marginTop:'4px' }}>
            <input className="inp" style={{ fontSize:'11px', padding:'4px 8px' }} placeholder="add a comment..." />
            <button className="btn btn-p" style={{ fontSize:'11px' }}>→</button>
          </div>
        </div>
      </div>

      <div className="card">
        <div style={{ display: 'flex', alignItems: 'center', gap:'8px', marginBottom:'8px' }}>
          <div className="av" style={{ width:'30px', height:'30px', background:'#E1F5EE', color:'#0F6E56', fontSize:'11px' }}>JM</div>
          <div>
            <p style={{ fontSize:'12px', fontWeight:500, color:'var(--color-text-primary)' }}>Jonas M. <span className="tag tag-gray">followers</span></p>
            <p style={{ fontSize:'10px', color:'var(--color-text-tertiary)' }}>5h ago · go devs group</p>
          </div>
        </div>
        <p style={{ fontSize:'12px', color:'var(--color-text-primary)', lineHeight:1.6, marginBottom:'8px' }}>New group event coming up — Docker deep dive this Friday. Everyone welcome.</p>
        <div className="event-card">
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems:'flex-start', marginBottom:'6px' }}>
            <div>
              <p style={{ fontSize:'12px', fontWeight:500, color:'var(--color-text-primary)' }}>Docker deep dive</p>
              <p style={{ fontSize:'11px', color:'var(--color-text-secondary)' }}>go devs · Friday 13 Jun · 18:00</p>
            </div>
            <span className="tag tag-purple">event</span>
          </div>
          <p style={{ fontSize:'11px', color:'var(--color-text-secondary)', marginBottom:'8px' }}>Hands-on session covering multi-stage builds, compose setups and deployment pipelines for the social network project.</p>
          <div style={{ display: 'flex', gap:'6px' }}>
            <button className="btn btn-t" style={{ display: 'flex', alignItems: 'center', gap:'3px' }}><i className="ti ti-check" style={{ fontSize:'12px' }} aria-hidden="true"></i> going (11)</button>
            <button className="btn btn-red" style={{ display: 'flex', alignItems: 'center', gap:'3px' }}><i className="ti ti-x" style={{ fontSize:'12px' }} aria-hidden="true"></i> not going (3)</button>
          </div>
        </div>
      </div>
    </main>
  );
}