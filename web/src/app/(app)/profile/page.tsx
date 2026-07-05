// SIDEBAR :

// {/* <aside className="sidebar">
// <p className="sec-label">navigate</p>
// {/* <div className="navlink" onClick={ showScreen('feed') }><i className="ti ti-home" style={{ fontSize:'14px' }} aria-hidden="true"></i> feed</div>
// <div className="navlink active" onClick={ showScreen('profile') }><i className="ti ti-user" style={{ fontSize:'14px' }} aria-hidden="true"></i> profile</div>
// <div className="navlink" onClick={ showScreen('groups') }><i className="ti ti-users" style={{ fontSize:'14px' }} aria-hidden="true"></i> groups</div>
// <div className="navlink" onClick={ showScreen('chat') }><i className="ti ti-message" style={{ fontSize:'14px' }} aria-hidden="true"></i> messages</div>
// <div className="navlink" onClick={ showScreen('notifications') }><i className="ti ti-bell" style={{ fontSize:'14px' }} aria-hidden="true"></i> notifications</div> */}
// </aside> */}

// right sidebar :

// <aside className="sidebar2">
//     <p className="sec-label" style={{ padding:0, marginBottom:'8px' }}>follow requests</p>
//     <div style={{ background:'#FBEAF0', border:'0.5px solid #ED93B1', padding:'8px', marginBottom:'6px', fontSize:'11px' }}>
//         <div style={{ display: 'flex', alignItems: 'center', gap:'6px', marginBottom:'6px' }}>
//         <div className="av" style={{ width:'22px', height:'22px', background:'#FBEAF0', color:'#993556', fontSize:'10px' }}>TP</div>
//         <span style={{ color:'#72243E', fontWeight:500 }}>Talia P.</span>
//         <span style={{ color:'#993556', fontSize:'10px' }}>wants to follow you</span>
//         </div>
//         <div style={{ display: 'flex', gap:'4px' }}><button className="btn btn-t" style={{ fontSize:'10px', padding:'3px 8px' }}>accept</button><button className="btn btn-red" style={{ fontSize:'10px', padding:'3px 8px' }}>decline</button></div>
//     </div>
//     <div className="divider"></div>
//     <p className="sec-label" style={{ padding:0, marginBottom:'8px' }}>send follow</p>
//     <div style={{ fontSize:'11px', color:'var(--color-text-secondary)', marginBottom:'6px' }}>find a user and send a follow request</div>
//     <input className="inp" style={{ fontSize:'11px', marginBottom:'6px' }} placeholder="search by name or email..." />
//     <div style={{ border:'0.5px solid var(--color-border-tertiary)', padding:'7px', display: 'flex', alignItems: 'center', gap:'6px', background:'var(--color-background-secondary)' }}>
//         <div className="av" style={{ width:'24px', height:'24px', background:'#E1F5EE', color:'#0F6E56', fontSize:'10px' }}>JM</div>
//         <span style={{ fontSize:'11px', color:'var(--color-text-primary)' }}>Jonas M.</span>
//         <button className="btn btn-p" style={{ marginLeft: 'auto', fontSize:'10px', padding:'3px 8px' }}>+ follow</button>
//     </div>
// </aside>


export default function Profile() {
  return (
        <main className="main">
            <div className="card">
                <div style={{ display: 'flex', alignItems:'flex-start', gap:'12px', marginBottom:'12px' }}>
                <div className="av" style={{ width:'52px', height:'52px', background:'#EEEDFE', color:'#534AB7', fontSize:'16px' }}>AK</div>
                <div style={{ flex:1 }}>
                    <div style={{ display: 'flex', alignItems: 'center', gap:'8px', flexWrap: 'wrap', marginBottom:'4px' }}>
                    <p style={{ fontSize:'14px', fontWeight:500, color:'var(--color-text-primary)' }}>Amir Kader</p>
                    <span className="tag tag-teal">@amirkader</span>
                    <span className="tag tag-purple" id="profile-visibility-tag">public</span>
                    {/* <button className="btn btn-g" style={{ fontSize:'10px', display: 'flex', alignItems: 'center', gap:'3px', marginLeft: 'auto' }} id="visibility-btn" onClick={ toggleVisibility() }><i className="ti ti-lock" style={{ fontSize:'12px' }} aria-hidden="true"></i> make private</button> */}
                    </div>
                    <p style={{ fontSize:'11px', color:'var(--color-text-secondary)', marginBottom:'4px' }}>Born 1998-07-14 · amir@example.com</p>
                    <p style={{ fontSize:'12px', color:'var(--color-text-primary)', lineHeight:1.5 }}>Full-stack dev working on distributed systems. Passionate about Go, open source, and clean architecture.</p>
                </div>
                </div>
                <div className="divider"></div>
                <div style={{ display: 'grid', gridTemplateColumns: 'repeat(3,1fr)', gap:'8px', textAlign: 'center' }}>
                <div style={{ background:'var(--color-background-secondary)', padding:'8px' }}>
                    <p style={{ fontSize:'18px', fontWeight:500, color:'#D4537E' }}>12</p>
                    <p style={{ fontSize:'10px', color:'var(--color-text-tertiary)' }}>posts</p>
                </div>
                {/* <div style={{ background:'var(--color-background-secondary)', padding:'8px', cursor: 'pointer' }} onClick={ showTab('followers') }>
                    <p style={{ fontSize:'18px', fontWeight:500, color:'#534AB7' }}>48</p>
                    <p style={{ fontSize:'10px', color:'var(--color-text-tertiary)' }}>followers</p>
                </div>
                <div style={{ background:'var(--color-background-secondary)', padding:'8px', cursor: 'pointer' }} onClick={ showTab('following') }>
                    <p style={{ fontSize:'18px', fontWeight:500, color:'#0F6E56' }}>31</p>
                    <p style={{ fontSize:'10px', color:'var(--color-text-tertiary)' }}>following</p>
                </div> */}
                </div>
            </div>

            <div style={{ display: 'flex', gap:0, border:'0.5px solid var(--color-border-tertiary)', background:'var(--color-background-primary)' }}>
                {/* <div className="profile-tab active-tab" id="tab-posts" style={{ padding:'8px 16px', fontSize:'11px', cursor: 'pointer', borderRight:'0.5px solid var(--color-border-tertiary)', color:'#D4537E', borderBottom:'2px solid #D4537E' }} onClick={ showTab('posts') }>posts</div>
                <div className="profile-tab" id="tab-followers" style={{ padding:'8px 16px', fontSize:'11px', cursor: 'pointer', borderRight:'0.5px solid var(--color-border-tertiary)', color:'var(--color-text-secondary)' }} onClick={ showTab('followers') }>followers</div> */}
                {/* <div className="profile-tab" id="tab-following" style={{ padding:'8px 16px', fontSize:'11px', cursor: 'pointer', color:'var(--color-text-secondary)' }} onClick={ showTab('following') }>following</div> */}
            </div>

            <div id="profile-posts">
                <div className="card" style={{ marginBottom:'6px' }}>
                <p style={{ fontSize:'10px', color:'var(--color-text-tertiary)', marginBottom:'4px' }}>2 days ago · <span className="tag tag-teal">public</span></p>
                <p style={{ fontSize:'12px', color:'var(--color-text-primary)', lineHeight:1.5 }}>Just deployed the migration system — golang-migrate running smooth with our SQLite schema. All 6 tables up.</p>
                </div>
                <div className="card">
                <p style={{ fontSize:'10px', color:'var(--color-text-tertiary)', marginBottom:'4px' }}>5 days ago · <span className="tag tag-gray">followers</span></p>
                <p style={{ fontSize:'12px', color:'var(--color-text-primary)', lineHeight:1.5 }}>Bcrypt session cookie auth is done. Logout clears the cookie and invalidates the server-side session.</p>
                </div>
            </div>
            <div id="profile-followers" style={{ display: 'none' }}>
                <div className="card">
                <div style={{ display: 'flex', flexDirection: 'column', gap:'10px' }}>
                    <div style={{ display: 'flex', alignItems: 'center', gap:'8px' }}><div className="av" style={{ width:'28px', height:'28px', background:'#FBEAF0', color:'#993556', fontSize:'11px' }}>SR</div><span style={{ fontSize:'12px', fontWeight:500, color:'var(--color-text-primary)' }}>Selin Rauf</span><button className="btn btn-g" style={{ marginLeft: 'auto', fontSize:'10px' }}>view profile</button></div>
                    <div style={{ display: 'flex', alignItems: 'center', gap:'8px' }}><div className="av" style={{ width:'28px', height:'28px', background:'#E1F5EE', color:'#0F6E56', fontSize:'11px' }}>JM</div><span style={{ fontSize:'12px', fontWeight:500, color:'var(--color-text-primary)' }}>Jonas M.</span><button className="btn btn-g" style={{ marginLeft: 'auto', fontSize:'10px' }}>view profile</button></div>
                    <div style={{ display: 'flex', alignItems: 'center', gap:'8px' }}><div className="av" style={{ width:'28px', height:'28px', background:'#EEEDFE', color:'#534AB7', fontSize:'11px' }}>OS</div><span style={{ fontSize:'12px', fontWeight:500, color:'var(--color-text-primary)' }}>Omar S.</span><button className="btn btn-g" style={{ marginLeft: 'auto', fontSize:'10px' }}>view profile</button></div>
                </div>
                </div>
            </div>
            <div id="profile-following" style={{ display: 'none' }}>
                <div className="card">
                <div style={{ display: 'flex', flexDirection: 'column', gap:'10px' }}>
                    <div style={{ display: 'flex', alignItems: 'center', gap:'8px' }}><div className="av" style={{ width:'28px', height:'28px', background:'#FBEAF0', color:'#993556', fontSize:'11px' }}>SR</div><span style={{ fontSize:'12px', fontWeight:500, color:'var(--color-text-primary)' }}>Selin Rauf</span><button className="btn btn-red" style={{ marginLeft: 'auto', fontSize:'10px' }}>unfollow</button></div>
                    <div style={{ display: 'flex', alignItems: 'center', gap:'8px' }}><div className="av" style={{ width:'28px', height:'28px', background:'#FAEEDA', color:'#633806', fontSize:'11px' }}>LK</div><span style={{ fontSize:'12px', fontWeight:500, color:'var(--color-text-primary)' }}>Lena K.</span><button className="btn btn-red" style={{ marginLeft: 'auto', fontSize:'10px' }}>unfollow</button></div>
                </div>
                </div>
            </div>
        </main>
    );
}