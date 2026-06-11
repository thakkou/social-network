export default function LoginPage() {
  return (
    <div className="screen active" id="screen-notifications">
      <div className="layout">
        <aside className="sidebar">
          <p className="sec-label">navigate</p>
          {/* <div className="navlink" onClick={ showScreen('feed') }><i className="ti ti-home" style={{ fontSize:'14px' }} aria-hidden="true"></i> feed</div>
          <div className="navlink" onClick={ showScreen('profile') }><i className="ti ti-user" style={{ fontSize:'14px' }} aria-hidden="true"></i> profile</div>
          <div className="navlink" onClick={ showScreen('groups') }><i className="ti ti-users" style={{ fontSize:'14px' }} aria-hidden="true"></i> groups</div>
          <div className="navlink" onClick={ showScreen('chat') }><i className="ti ti-message" style={{ fontSize:'14px' }} aria-hidden="true"></i> messages</div>
          <div className="navlink active" onClick={ showScreen('notifications') }><i className="ti ti-bell" style={{ fontSize:'14px' }} aria-hidden="true"></i> notifications</div> */}
        </aside>
        <main className="main">
          <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom:'4px' }}>
            <p style={{ fontSize:'13px', fontWeight:500, color:'var(--color-text-primary)' }}>notifications</p>
            <button className="btn btn-g" style={{ fontSize:'10px' }}>mark all read</button>
          </div>

          <div className="card" style={{ borderLeft:'2px solid #D4537E' }}>
            <div style={{ display: 'flex', alignItems: 'center', gap:'8px', marginBottom:'8px' }}>
              <span className="tag tag-pink">follow request</span>
              <span style={{ fontSize:'10px', color:'var(--color-text-tertiary)', marginLeft: 'auto' }}>just now</span>
            </div>
            <div style={{ display: 'flex', alignItems: 'center', gap:'8px', marginBottom:'8px' }}>
              <div className="av" style={{ width:'28px', height:'28px', background:'#FBEAF0', color:'#993556', fontSize:'11px' }}>LK</div>
              <p style={{ fontSize:'12px', color:'var(--color-text-primary)' }}><span style={{ fontWeight:500 }}>Lena K.</span> sent you a follow request</p>
            </div>
            <div style={{ display: 'flex', gap:'6px' }}><button className="btn btn-t" style={{ fontSize:'11px', display: 'flex', alignItems: 'center', gap:'3px' }}><i className="ti ti-check" style={{ fontSize:'12px' }} aria-hidden="true"></i> accept</button><button className="btn btn-red" style={{ fontSize:'11px', display: 'flex', alignItems: 'center', gap:'3px' }}><i className="ti ti-x" style={{ fontSize:'12px' }} aria-hidden="true"></i> decline</button></div>
          </div>

          <div className="card" style={{ borderLeft:'2px solid #7F77DD' }}>
            <div style={{ display: 'flex', alignItems: 'center', gap:'8px', marginBottom:'8px' }}>
              <span className="tag tag-purple">group invite</span>
              <span style={{ fontSize:'10px', color:'var(--color-text-tertiary)', marginLeft: 'auto' }}>1h ago</span>
            </div>
            <div style={{ display: 'flex', alignItems: 'center', gap:'8px', marginBottom:'8px' }}>
              <span style={{ width:'8px', height:'8px', background:'#534AB7', flexShrink:0 }}></span>
              <p style={{ fontSize:'12px', color:'var(--color-text-primary)' }}><span style={{ fontWeight:500 }}>Selin R.</span> invited you to join <span style={{ fontWeight:500, color:'#534AB7' }}>design systems</span></p>
            </div>
            <div style={{ display: 'flex', gap:'6px' }}><button className="btn btn-t" style={{ fontSize:'11px' }}>accept invite</button><button className="btn btn-red" style={{ fontSize:'11px' }}>decline</button></div>
          </div>

          <div className="card" style={{ borderLeft:'2px solid #7F77DD', opacity:0.7 }}>
            <div style={{ display: 'flex', alignItems: 'center', gap:'8px', marginBottom:'6px' }}>
              <span className="tag tag-purple">join request</span>
              <span style={{ fontSize:'10px', color:'var(--color-text-tertiary)', marginLeft: 'auto' }}>3h ago · go devs</span>
            </div>
            <div style={{ display: 'flex', alignItems: 'center', gap:'8px', marginBottom:'8px' }}>
              <div className="av" style={{ width:'28px', height:'28px', background:'#EEEDFE', color:'#534AB7', fontSize:'11px' }}>OS</div>
              <p style={{ fontSize:'12px', color:'var(--color-text-primary)' }}><span style={{ fontWeight:500 }}>Omar S.</span> requested to join your group <span style={{ fontWeight:500, color:'#D4537E' }}>go devs</span></p>
            </div>
            <div style={{ display: 'flex', gap:'6px' }}><button className="btn btn-t" style={{ fontSize:'11px' }}>accept</button><button className="btn btn-red" style={{ fontSize:'11px' }}>decline</button></div>
          </div>

          <div className="card" style={{ borderLeft:'2px solid #1D9E75', opacity:0.7 }}>
            <div style={{ display: 'flex', alignItems: 'center', gap:'8px', marginBottom:'6px' }}>
              <span className="tag tag-teal">event</span>
              <span style={{ fontSize:'10px', color:'var(--color-text-tertiary)', marginLeft: 'auto' }}>5h ago · go devs</span>
            </div>
            <p style={{ fontSize:'12px', color:'var(--color-text-primary)' }}>New event <span style={{ fontWeight:500 }}>Docker deep dive</span> was created in <span style={{ fontWeight:500, color:'#D4537E' }}>go devs</span> — Friday 18:00</p>
            <div style={{ display: 'flex', gap:'6px', marginTop:'8px' }}><button className="btn btn-t" style={{ fontSize:'11px', display: 'flex', alignItems: 'center', gap:'3px' }}><i className="ti ti-check" style={{ fontSize:'12px' }} aria-hidden="true"></i> going</button><button className="btn btn-g" style={{ fontSize:'11px' }}>not going</button></div>
          </div>
        </main>
        <aside className="sidebar2">
          <p className="sec-label" style={{ padding:0, marginBottom:'8px' }}>messages</p>
          <div style={{ fontSize:'11px', color:'var(--color-text-secondary)', marginBottom:'8px' }}>private messages are separate from notifications</div>
          {/* <div style={{ border:'0.5px solid var(--color-border-tertiary)', padding:'8px', background:'var(--color-background-secondary)', marginBottom:'6px', cursor: 'pointer' }} onClick={ showScreen('chat') }>
            <div style={{ display: 'flex', alignItems: 'center', gap:'6px', marginBottom:'3px' }}>
              <div className="av" style={{ width:'20px', height:'20px', background:'#FBEAF0', color:'#993556', fontSize:'9px' }}>SR</div>
              <span style={{ fontSize:'11px', fontWeight:500, color:'var(--color-text-primary)' }}>Selin R.</span>
              <span className="notif-dot" style={{ marginLeft: 'auto' }}>2</span>
            </div>
            <p style={{ fontSize:'11px', color:'var(--color-text-secondary)', paddingLeft:'26px' }}>just merged it, take a look</p>
          </div>
          <div style={{ border:'0.5px solid var(--color-border-tertiary)', padding:'8px', background:'var(--color-background-secondary)', cursor: 'pointer' }} onClick={ showScreen('chat') }>
            <div style={{ display: 'flex', alignItems: 'center', gap:'6px', marginBottom:'3px' }}>
              <div className="av" style={{ width:'20px', height:'20px', background:'#E1F5EE', color:'#0F6E56', fontSize:'9px' }}>JM</div>
              <span style={{ fontSize:'11px', fontWeight:500, color:'var(--color-text-primary)' }}>Jonas M.</span>
            </div>
            <p style={{ fontSize:'11px', color:'var(--color-text-secondary)', paddingLeft:'26px' }}>great work on the auth layer</p>
          </div> */}
          {/* <div style={{ textAlign: 'center', marginTop:'8px' }}><button className="btn btn-p" style={{ fontSize:'10px', width: '100%' }} onClick={ showScreen('chat') }>open messages →</button></div> */}
        </aside>
      </div>
    </div>
  );
}