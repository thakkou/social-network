// sidebar :

// {/* <aside className="sidebar">
//   <p className="sec-label">navigate</p>
//   {/* <div className="navlink" onClick={ showScreen('feed') }><i className="ti ti-home" style={{ fontSize:'14px' }} aria-hidden="true"></i> feed</div>
//   <div className="navlink" onClick={ showScreen('profile') }><i className="ti ti-user" style={{ fontSize:'14px' }} aria-hidden="true"></i> profile</div>
//   <div className="navlink active" onClick={ showScreen('groups') }><i className="ti ti-users" style={{ fontSize:'14px' }} aria-hidden="true"></i> groups</div>
//   <div className="navlink" onClick={ showScreen('chat') }><i className="ti ti-message" style={{ fontSize:'14px' }} aria-hidden="true"></i> messages</div>
//   <div className="navlink" onClick={ showScreen('notifications') }><i className="ti ti-bell" style={{ fontSize:'14px' }} aria-hidden="true"></i> notifications</div> */}
// </aside> */}

// right sidebar :

{/* <aside className="sidebar2">
  <p className="sec-label" style={{ padding:0, marginBottom:'8px' }}>pending requests</p>
  <div style={{ background:'#EEEDFE', border:'0.5px solid #AFA9EC', padding:'8px', fontSize:'11px', marginBottom:'6px' }}>
    <p style={{ fontWeight:500, color:'#3C3489', marginBottom:'2px' }}>Omar S.</p>
    <p style={{ color:'#534AB7', marginBottom:'5px' }}>wants to join go devs</p>
    <div style={{ display: 'flex', gap:'4px' }}><button className="btn btn-t" style={{ fontSize:'10px', padding:'3px 8px' }}>accept</button><button className="btn btn-red" style={{ fontSize:'10px', padding:'3px 8px' }}>decline</button></div>
  </div>
  <div className="divider"></div>
  <p className="sec-label" style={{ padding:0, marginBottom:'8px' }}>create event</p>
  <div className="form-row"><span className="form-label">title</span><input className="inp" style={{ fontSize:'11px' }} placeholder="event title" /></div>
  <div className="form-row"><span className="form-label">description</span><textarea className="inp" style={{ height:'44px', resize: 'none', fontSize:'11px' }} placeholder="details..."></textarea></div>
  <div style={{ display: 'grid', gridTemplateColumns:'1fr 1fr', gap:'6px', marginBottom:'8px' }}>
    <div className="form-row"><span className="form-label">date</span><input className="inp" type="date" style={{ fontSize:'11px' }} /></div>
    <div className="form-row"><span className="form-label">time</span><input className="inp" type="time" style={{ fontSize:'11px' }} /></div>
  </div>
  <button className="btn btn-p" style={{ width:'100%', fontSize:'11px' }}>create event →</button>
</aside> */}

export default function Groups() {
  return (
    <main className="main">
      <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
        <p style={{ fontSize:'13px', fontWeight:500, color:'var(--color-text-primary)' }}>browse all groups</p>
        {/* <button className="btn btn-p" onClick={ toggleCreateGroup() } style={{ display: 'flex', alignItems: 'center', gap:'4px' }}><i className="ti ti-plus" style={{ fontSize:'12px' }} aria-hidden="true"></i> create group</button> */}
      </div>

      <div id="create-group-panel" style={{ display: 'none' }}>
        <div className="card">
          <p style={{ fontSize:'11px', fontWeight:500, color:'var(--color-text-primary)', marginBottom:'10px' }}>new group</p>
          <div className="form-row"><span className="form-label">title *</span><input className="inp" placeholder="group name" /></div>
          <div className="form-row"><span className="form-label">description *</span><textarea className="inp" style={{ height:'56px', resize: 'none' }} placeholder="what is this group about?"></textarea></div>
          <div className="form-row"><span className="form-label">invite members</span><input className="inp" placeholder="search by name..." /></div>
          <div style={{ display: 'flex', gap:'6px', marginTop:'4px' }}>
            <button className="btn btn-p">create →</button>
            {/* <button className="btn btn-g" onClick={ toggleCreateGroup() }>cancel</button> */}
          </div>
        </div>
      </div>

      <input className="inp" placeholder="search groups..." style={{ fontSize:'12px' }} />

      <div className="card">
        <div style={{ display: 'flex', alignItems: 'center', gap:'8px', marginBottom:'8px' }}>
          <span style={{ width:'8px', height:'8px', background:'#D4537E', flexShrink:0 }}></span>
          <p style={{ fontSize:'13px', fontWeight:500, color:'var(--color-text-primary)' }}>go devs</p>
          <span className="tag tag-pink" style={{ marginLeft: 'auto' }}>member</span>
        </div>
        <p style={{ fontSize:'12px', color:'var(--color-text-secondary)', marginBottom:'8px' }}>A space for Go enthusiasts — sharing projects, tips and co-building the social network backend.</p>
        <div style={{ fontSize:'11px', color:'var(--color-text-tertiary)', marginBottom:'10px' }}>34 members · created by Amir K.</div>
        <div className="event-card" style={{ marginBottom:'8px' }}>
          <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom:'4px' }}>
            <p style={{ fontSize:'12px', fontWeight:500, color:'var(--color-text-primary)' }}>Docker deep dive</p>
            <span className="tag tag-amber">Fri 18:00</span>
          </div>
          <p style={{ fontSize:'11px', color:'var(--color-text-secondary)', marginBottom:'6px' }}>Multi-stage builds, compose, and deployment pipelines.</p>
          <div style={{ display: 'flex', gap:'6px' }}><button className="btn btn-t" style={{ fontSize:'10px', display: 'flex', alignItems: 'center', gap:'3px' }}><i className="ti ti-check" style={{ fontSize:'11px' }} aria-hidden="true"></i> going (11)</button><button className="btn btn-g" style={{ fontSize:'10px' }}>not going (3)</button></div>
        </div>
        <div style={{ display: 'flex', gap:'6px' }}>
          {/* <button className="btn btn-p" onClick={ showScreen('chat') } style={{ fontSize:'11px', display: 'flex', alignItems: 'center', gap:'3px' }}><i className="ti ti-message" style={{ fontSize:'12px' }} aria-hidden="true"></i> group chat</button> */}
          <button className="btn btn-g" style={{ fontSize:'11px', display: 'flex', alignItems: 'center', gap:'3px' }}><i className="ti ti-calendar-plus" style={{ fontSize:'12px' }} aria-hidden="true"></i> add event</button>
          <button className="btn btn-g" style={{ fontSize:'11px', display: 'flex', alignItems: 'center', gap:'3px' }}><i className="ti ti-user-plus" style={{ fontSize:'12px' }} aria-hidden="true"></i> invite</button>
        </div>
      </div>

      <div className="card">
        <div style={{ display: 'flex', alignItems: 'center', gap:'8px', marginBottom:'8px' }}>
          <span style={{ width:'8px', height:'8px', background:'#534AB7', flexShrink:0 }}></span>
          <p style={{ fontSize:'13px', fontWeight:500, color:'var(--color-text-primary)' }}>design systems</p>
          <span className="tag tag-gray" style={{ marginLeft: 'auto' }}>not a member</span>
        </div>
        <p style={{ fontSize:'12px', color:'var(--color-text-secondary)', marginBottom:'8px' }}>Discussing component libraries, tokens, and consistent UI patterns across projects.</p>
        <div style={{ fontSize:'11px', color:'var(--color-text-tertiary)', marginBottom:'10px' }}>21 members · created by Selin R.</div>
        <button className="btn btn-p" style={{ fontSize:'11px' }}>request to join →</button>
      </div>
    </main>
  );
}