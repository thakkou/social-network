import Link from "next/link";
import type { MouseEventHandler } from "react";

function showScreen(name: string): MouseEventHandler<HTMLDivElement | HTMLButtonElement> | undefined {
  // document.querySelectorAll('.screen').forEach(s => s.classList.remove('active'));
  // const t = document.getElementById('screen-' + name);
  // if (t) t.classList.add('active');
  // const isAuth = name === 'auth';
  // const navTop = document.getElementById('nav-top');
  // if (navTop) {
  //   navTop.style.display = isAuth ? 'none' : 'flex';
  // }
  return undefined
}

function showMainScreen(): MouseEventHandler<HTMLDivElement | HTMLButtonElement> | undefined {
  showScreen('feed');
  return undefined
}

function showRegister(): MouseEventHandler<HTMLDivElement> | undefined {
  return undefined
}

function toggleVisibility(): MouseEventHandler<HTMLDivElement | HTMLButtonElement> | undefined {
  // const tag = document.getElementById('profile-visibility-tag');
  // const btn = document.getElementById('visibility-btn');
  // if (tag && btn) {
  //   if (tag.textContent === 'public') {
  //     tag.textContent = 'private';
  //     tag.className = 'tag tag-gray';
  //     btn.innerHTML = '<i class="ti ti-lock-open" style={{ fontSize: \'12px\' }} aria-hidden="true"></i> make public';
  //   } else {
  //     tag.textContent = 'public';
  //     tag.className = 'tag tag-purple';
  //     btn.innerHTML = '<i class="ti ti-lock" style={{ fontSize: \'12px\' }} aria-hidden="true"></i> make private';
  //   }
  // }
  return undefined
}

function showTab(tab: string): MouseEventHandler<HTMLDivElement> | undefined {
  // ['posts','followers','following'].forEach(t => {
  //   const el = document.getElementById('profile-' + t);
  //   const tab_el = document.getElementById('tab-' + t);
  //   if (el) el.style.display = t === tab ? 'block' : 'none';
  //   if (tab_el) {
  //     tab_el.style.color = t === tab ? '#D4537E' : 'var(--color-text-secondary)';
  //     tab_el.style.borderBottom = t === tab ? '2px solid #D4537E' : 'none';
  //   }
  // });
  return undefined
}

function toggleCreateGroup(): MouseEventHandler<HTMLDivElement | HTMLButtonElement> | undefined {
  // const p = document.getElementById('create-group-panel');
  // if (p) p.style.display = p.style.display === 'none' ? 'block' : 'none';
  return undefined
}

function sendMsg(): MouseEventHandler<HTMLDivElement | HTMLButtonElement> | undefined {
  // const inp = document.getElementById('chat-input');
  // const val = inp.value.trim();
  // if (!val) return;
  // const msgs = document.getElementById('chat-messages');
  // const b = document.createElement('div');
  // b.className = 'msg-bubble-me';
  // b.textContent = val;
  // msgs.appendChild(b);
  // inp.value = '';
  // msgs.scrollTop = msgs.scrollHeight;
  return undefined
}

function switchChat(who:string): MouseEventHandler<HTMLDivElement> | undefined {
  return undefined
}

export default function HomePage() {
  return (
    <main>
    <h2 className="sr-only">Social network full UI — all screens including auth, feed, profile, groups, chat, and notifications</h2>

    <div id="app">

    <div className="nav">
      <span className="logo">social<span>net</span></span>
      <div style={{ display: 'flex', gap:'4px', alignItems: 'center' }} id="nav-top">
        <input className="inp" style={{ width:'180px',  padding:'4px 8px', fontSize:'11px' }} placeholder="search..." />
        <button className="btn btn-g" style={{ display: 'flex', alignItems: 'center', gap:'4px' }} onClick={ showScreen('notifications') }>
          <i className="ti ti-bell" style={{ fontSize:'14px' }} aria-hidden="true"></i>
          <span className="notif-dot">4</span>
        </button>
        <div className="av" style={{ width:'28px', height:'28px', background:'#EEEDFE', color:'#534AB7', fontSize:'11px', cursor: 'pointer' }} onClick={ showScreen('profile') }>AK</div>
      </div>
    </div>

    <div className="screen active" id="screen-auth">
      <div style={{ display: 'grid', gridTemplateColumns:'1fr 1fr', minHeight:'560px' }}>
        <div style={{ background:'var(--color-background-primary)', borderRight:'0.5px solid var(--color-border-tertiary)', padding:'2rem', display: 'flex', flexDirection: 'column', justifyContent: 'center' }}>
          <p className="sec-label" style={{ padding:0, marginBottom:'1rem' }}>sign in</p>
          <div className="form-row"><span className="form-label">email</span><input className="inp" type="email" placeholder="you@example.com" /></div>
          <div className="form-row"><span className="form-label">password</span><input className="inp" type="password" placeholder="••••••••" /></div>
          <button className="btn btn-p" style={{ width:'100%', marginTop:'8px' }} onClick={ showMainScreen() }>sign in →</button>
          <p style={{ fontSize:'11px', color:'var(--color-text-tertiary)', marginTop:'10px', textAlign: 'center' }}>no account? <span style={{ color:'#D4537E', cursor: 'pointer' }} onClick={ showRegister() }>register</span></p>
        </div>
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
            <div style={{ border:'0.5px dashed var(--color-border-tertiary)', padding:'10px', textAlign: 'center', fontSize:'11px', color:'var(--color-text-tertiary)', cursor: 'pointer' }}>
              <i className="ti ti-upload" style={{ fontSize:'16px', display: 'block', marginBottom:'4px' }} aria-hidden="true"></i>
              jpeg / png / gif
            </div>
          </div>
          <button className="btn btn-p" style={{ width:'100%', marginTop:'6px' }} onClick={ showMainScreen() }>create account →</button>
        </div>
      </div>
    </div>

    <div className="screen" id="screen-feed">
      <div className="layout">
        <aside className="sidebar" id="main-sidebar">
          <p className="sec-label">navigate</p>
          <div className="navlink active" onClick={ showScreen('feed') }><i className="ti ti-home" style={{ fontSize:'14px' }} aria-hidden="true"></i> feed</div>
          <div className="navlink" onClick={ showScreen('profile') }><i className="ti ti-user" style={{ fontSize:'14px' }} aria-hidden="true"></i> profile</div>
          <div className="navlink" onClick={ showScreen('groups') }><i className="ti ti-users" style={{ fontSize:'14px' }} aria-hidden="true"></i> groups</div>
          <div className="navlink" onClick={ showScreen('chat') }><i className="ti ti-message" style={{ fontSize:'14px' }} aria-hidden="true"></i> messages <span className="notif-dot" style={{ marginLeft: 'auto' }}>2</span></div>
          <div className="navlink" onClick={ showScreen('notifications') }><i className="ti ti-bell" style={{ fontSize:'14px' }} aria-hidden="true"></i> notifications <span className="notif-dot" style={{ marginLeft: 'auto' }}>4</span></div>
          <div className="divider" style={{ margin:'0.5rem 0.75rem' }}></div>
          <p className="sec-label">my groups</p>
          <div className="navlink" style={{ fontSize:'11px' }}><span style={{ width:'6px', height:'6px', background:'#D4537E', display: 'inline-block', flexShrink:0 }}></span> go devs</div>
          <div className="navlink" style={{ fontSize:'11px' }}><span style={{ width:'6px', height:'6px', background:'#534AB7', display: 'inline-block', flexShrink:0 }}></span> design sys</div>
          <div className="navlink" style={{ fontSize:'11px' }}><span style={{ width:'6px', height:'6px', background:'#1D9E75', display: 'inline-block', flexShrink:0 }}></span> open src</div>
          <div style={{ padding:'6px 12px', marginTop:'4px' }}><button className="btn btn-g" style={{ width:'100%', fontSize:'11px', display: 'flex', alignItems: 'center', justifyContent: 'center', gap:'4px' }} onClick={ showScreen('groups') }><i className="ti ti-plus" style={{ fontSize:'12px' }} aria-hidden="true"></i> new group</button></div>
        </aside>
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
            <div style={{ background:'var(--color-background-secondary)', border:'0.5px solid var(--color-border-tertiary)', padding:'8px 12px', fontSize:'11px', marginBottom:'8px', color:'var(--color-text-secondary)' }}>
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
        <aside className="sidebar2">
          <p className="sec-label" style={{ padding:0, marginBottom:'8px' }}>notifications</p>
          <div className="notif-item" style={{ background:'#FBEAF0', borderLeftColor:'#D4537E' }}><span style={{ fontWeight:500, color:'#72243E' }}>Lena K.</span> <span style={{ color:'#993556' }}>sent a follow request</span>
            <div style={{ display: 'flex', gap:'4px', marginTop:'5px' }}><button className="btn btn-t" style={{ fontSize:'10px', padding:'3px 8px' }}>accept</button><button className="btn btn-red" style={{ fontSize:'10px', padding:'3px 8px' }}>decline</button></div>
          </div>
          <div className="notif-item" style={{ background:'#EEEDFE', borderLeftColor:'#7F77DD' }}><span style={{ fontWeight:500, color:'#3C3489' }}>design sys</span> <span style={{ color:'#534AB7' }}>invited you</span>
            <div style={{ display: 'flex', gap:'4px', marginTop:'5px' }}><button className="btn btn-t" style={{ fontSize:'10px', padding:'3px 8px' }}>accept</button><button className="btn btn-red" style={{ fontSize:'10px', padding:'3px 8px' }}>decline</button></div>
          </div>
          <div className="notif-item" style={{ background:'#E1F5EE', borderLeftColor:'#1D9E75' }}><span style={{ color:'#085041' }}>event created in <span style={{ fontWeight:500 }}>go devs</span></span></div>
          <div className="divider"></div>
          <p className="sec-label" style={{ padding:0, marginBottom:'8px' }}>online now</p>
          <div style={{ display: 'flex', flexDirection: 'column', gap:'7px' }}>
            <div style={{ display: 'flex', alignItems: 'center', gap:'7px' }}>
              <div className="av" style={{ width:'26px', height:'26px', background:'#FBEAF0', color:'#993556', fontSize:'10px' }}>SR</div>
              <span style={{ fontSize:'11px', color:'var(--color-text-primary)' }}>Selin R.</span>
              <span className="online-dot" style={{ marginLeft: 'auto' }}></span>
            </div>
            <div style={{ display: 'flex', alignItems: 'center', gap:'7px' }}>
              <div className="av" style={{ width:'26px', height:'26px', background:'#E1F5EE', color:'#0F6E56', fontSize:'10px' }}>JM</div>
              <span style={{ fontSize:'11px', color:'var(--color-text-primary)' }}>Jonas M.</span>
              <span className="online-dot" style={{ marginLeft: 'auto' }}></span>
            </div>
            <div style={{ display: 'flex', alignItems: 'center', gap:'7px' }}>
              <div className="av" style={{ width:'26px', height:'26px', background:'#EEEDFE', color:'#534AB7', fontSize:'10px' }}>OS</div>
              <span style={{ fontSize:'11px', color:'var(--color-text-primary)' }}>Omar S.</span>
              <span className="offline-dot" style={{ marginLeft: 'auto' }}></span>
            </div>
          </div>
        </aside>
      </div>
    </div>

    <div className="screen" id="screen-profile">
      <div className="layout">
        <aside className="sidebar">
          <p className="sec-label">navigate</p>
          <div className="navlink" onClick={ showScreen('feed') }><i className="ti ti-home" style={{ fontSize:'14px' }} aria-hidden="true"></i> feed</div>
          <div className="navlink active" onClick={ showScreen('profile') }><i className="ti ti-user" style={{ fontSize:'14px' }} aria-hidden="true"></i> profile</div>
          <div className="navlink" onClick={ showScreen('groups') }><i className="ti ti-users" style={{ fontSize:'14px' }} aria-hidden="true"></i> groups</div>
          <div className="navlink" onClick={ showScreen('chat') }><i className="ti ti-message" style={{ fontSize:'14px' }} aria-hidden="true"></i> messages</div>
          <div className="navlink" onClick={ showScreen('notifications') }><i className="ti ti-bell" style={{ fontSize:'14px' }} aria-hidden="true"></i> notifications</div>
        </aside>
        <main className="main">
          <div className="card">
            <div style={{ display: 'flex', alignItems:'flex-start', gap:'12px', marginBottom:'12px' }}>
              <div className="av" style={{ width:'52px', height:'52px', background:'#EEEDFE', color:'#534AB7', fontSize:'16px' }}>AK</div>
              <div style={{ flex:1 }}>
                <div style={{ display: 'flex', alignItems: 'center', gap:'8px', flexWrap: 'wrap', marginBottom:'4px' }}>
                  <p style={{ fontSize:'14px', fontWeight:500, color:'var(--color-text-primary)' }}>Amir Kader</p>
                  <span className="tag tag-teal">@amirkader</span>
                  <span className="tag tag-purple" id="profile-visibility-tag">public</span>
                  <button className="btn btn-g" style={{ fontSize:'10px', display: 'flex', alignItems: 'center', gap:'3px', marginLeft: 'auto' }} id="visibility-btn" onClick={ toggleVisibility() }><i className="ti ti-lock" style={{ fontSize:'12px' }} aria-hidden="true"></i> make private</button>
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
              <div style={{ background:'var(--color-background-secondary)', padding:'8px', cursor: 'pointer' }} onClick={ showTab('followers') }>
                <p style={{ fontSize:'18px', fontWeight:500, color:'#534AB7' }}>48</p>
                <p style={{ fontSize:'10px', color:'var(--color-text-tertiary)' }}>followers</p>
              </div>
              <div style={{ background:'var(--color-background-secondary)', padding:'8px', cursor: 'pointer' }} onClick={ showTab('following') }>
                <p style={{ fontSize:'18px', fontWeight:500, color:'#0F6E56' }}>31</p>
                <p style={{ fontSize:'10px', color:'var(--color-text-tertiary)' }}>following</p>
              </div>
            </div>
          </div>

          <div style={{ display: 'flex', gap:0, border:'0.5px solid var(--color-border-tertiary)', background:'var(--color-background-primary)' }}>
            <div className="profile-tab active-tab" id="tab-posts" style={{ padding:'8px 16px', fontSize:'11px', cursor: 'pointer', borderRight:'0.5px solid var(--color-border-tertiary)', color:'#D4537E', borderBottom:'2px solid #D4537E' }} onClick={ showTab('posts') }>posts</div>
            <div className="profile-tab" id="tab-followers" style={{ padding:'8px 16px', fontSize:'11px', cursor: 'pointer', borderRight:'0.5px solid var(--color-border-tertiary)', color:'var(--color-text-secondary)' }} onClick={ showTab('followers') }>followers</div>
            <div className="profile-tab" id="tab-following" style={{ padding:'8px 16px', fontSize:'11px', cursor: 'pointer', color:'var(--color-text-secondary)' }} onClick={ showTab('following') }>following</div>
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
        <aside className="sidebar2">
          <p className="sec-label" style={{ padding:0, marginBottom:'8px' }}>follow requests</p>
          <div style={{ background:'#FBEAF0', border:'0.5px solid #ED93B1', padding:'8px', marginBottom:'6px', fontSize:'11px' }}>
            <div style={{ display: 'flex', alignItems: 'center', gap:'6px', marginBottom:'6px' }}>
              <div className="av" style={{ width:'22px', height:'22px', background:'#FBEAF0', color:'#993556', fontSize:'10px' }}>TP</div>
              <span style={{ color:'#72243E', fontWeight:500 }}>Talia P.</span>
              <span style={{ color:'#993556', fontSize:'10px' }}>wants to follow you</span>
            </div>
            <div style={{ display: 'flex', gap:'4px' }}><button className="btn btn-t" style={{ fontSize:'10px', padding:'3px 8px' }}>accept</button><button className="btn btn-red" style={{ fontSize:'10px', padding:'3px 8px' }}>decline</button></div>
          </div>
          <div className="divider"></div>
          <p className="sec-label" style={{ padding:0, marginBottom:'8px' }}>send follow</p>
          <div style={{ fontSize:'11px', color:'var(--color-text-secondary)', marginBottom:'6px' }}>find a user and send a follow request</div>
          <input className="inp" style={{ fontSize:'11px', marginBottom:'6px' }} placeholder="search by name or email..." />
          <div style={{ border:'0.5px solid var(--color-border-tertiary)', padding:'7px', display: 'flex', alignItems: 'center', gap:'6px', background:'var(--color-background-secondary)' }}>
            <div className="av" style={{ width:'24px', height:'24px', background:'#E1F5EE', color:'#0F6E56', fontSize:'10px' }}>JM</div>
            <span style={{ fontSize:'11px', color:'var(--color-text-primary)' }}>Jonas M.</span>
            <button className="btn btn-p" style={{ marginLeft: 'auto', fontSize:'10px', padding:'3px 8px' }}>+ follow</button>
          </div>
        </aside>
      </div>
    </div>

    <div className="screen" id="screen-groups">
      <div className="layout">
        <aside className="sidebar">
          <p className="sec-label">navigate</p>
          <div className="navlink" onClick={ showScreen('feed') }><i className="ti ti-home" style={{ fontSize:'14px' }} aria-hidden="true"></i> feed</div>
          <div className="navlink" onClick={ showScreen('profile') }><i className="ti ti-user" style={{ fontSize:'14px' }} aria-hidden="true"></i> profile</div>
          <div className="navlink active" onClick={ showScreen('groups') }><i className="ti ti-users" style={{ fontSize:'14px' }} aria-hidden="true"></i> groups</div>
          <div className="navlink" onClick={ showScreen('chat') }><i className="ti ti-message" style={{ fontSize:'14px' }} aria-hidden="true"></i> messages</div>
          <div className="navlink" onClick={ showScreen('notifications') }><i className="ti ti-bell" style={{ fontSize:'14px' }} aria-hidden="true"></i> notifications</div>
        </aside>
        <main className="main">
          <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
            <p style={{ fontSize:'13px', fontWeight:500, color:'var(--color-text-primary)' }}>browse all groups</p>
            <button className="btn btn-p" onClick={ toggleCreateGroup() } style={{ display: 'flex', alignItems: 'center', gap:'4px' }}><i className="ti ti-plus" style={{ fontSize:'12px' }} aria-hidden="true"></i> create group</button>
          </div>

          <div id="create-group-panel" style={{ display: 'none' }}>
            <div className="card">
              <p style={{ fontSize:'11px', fontWeight:500, color:'var(--color-text-primary)', marginBottom:'10px' }}>new group</p>
              <div className="form-row"><span className="form-label">title *</span><input className="inp" placeholder="group name" /></div>
              <div className="form-row"><span className="form-label">description *</span><textarea className="inp" style={{ height:'56px', resize: 'none' }} placeholder="what is this group about?"></textarea></div>
              <div className="form-row"><span className="form-label">invite members</span><input className="inp" placeholder="search by name..." /></div>
              <div style={{ display: 'flex', gap:'6px', marginTop:'4px' }}>
                <button className="btn btn-p">create →</button>
                <button className="btn btn-g" onClick={ toggleCreateGroup() }>cancel</button>
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
              <button className="btn btn-p" onClick={ showScreen('chat') } style={{ fontSize:'11px', display: 'flex', alignItems: 'center', gap:'3px' }}><i className="ti ti-message" style={{ fontSize:'12px' }} aria-hidden="true"></i> group chat</button>
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
        <aside className="sidebar2">
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
        </aside>
      </div>
    </div>

    <div className="screen" id="screen-chat">
      <div className="layout">
        <aside className="sidebar">
          <p className="sec-label">navigate</p>
          <div className="navlink" onClick={ showScreen('feed') }><i className="ti ti-home" style={{ fontSize:'14px' }} aria-hidden="true"></i> feed</div>
          <div className="navlink" onClick={ showScreen('profile') }><i className="ti ti-user" style={{ fontSize:'14px' }} aria-hidden="true"></i> profile</div>
          <div className="navlink" onClick={ showScreen('groups') }><i className="ti ti-users" style={{ fontSize:'14px' }} aria-hidden="true"></i> groups</div>
          <div className="navlink active" onClick={ showScreen('chat') }><i className="ti ti-message" style={{ fontSize:'14px' }} aria-hidden="true"></i> messages</div>
          <div className="navlink" onClick={ showScreen('notifications') }><i className="ti ti-bell" style={{ fontSize:'14px' }} aria-hidden="true"></i> notifications</div>
          <div className="divider" style={{ margin: '0.5rem 0.75rem' }}></div>
          <p className="sec-label">direct</p>
          <div className="navlink active-chat" style={{ background:'#FBEAF0', borderLeft:'2px solid #D4537E', color:'#D4537E', fontSize:'11px' }} onClick={ switchChat('selin') }>
            <div className="av" style={{ width:'20px', height:'20px', background:'#FBEAF0', color:'#993556', fontSize:'9px' }}>SR</div> Selin R. <span className="notif-dot" style={{ marginLeft: 'auto' }}>2</span>
          </div>
          <div className="navlink" style={{ fontSize:'11px' }} onClick={ switchChat('jonas') }>
            <div className="av" style={{ width:'20px', height:'20px', background:'#E1F5EE', color:'#0F6E56', fontSize:'9px' }}>JM</div> Jonas M.
          </div>
          <div className="divider" style={{ margin: '0.5rem 0.75rem' }}></div>
          <p className="sec-label">groups</p>
          <div className="navlink" style={{ fontSize:'11px' }} onClick={ switchChat('godevs') }>
            <span style={{ width:'6px', height:'6px', background:'#D4537E', display: 'inline-block', flexShrink:0 }}></span> go devs
          </div>
        </aside>
        <main className="main" style={{ padding:0, gap:0 }}>
          <div style={{ background:'var(--color-background-primary)', borderBottom:'0.5px solid var(--color-border-tertiary)', padding:'10px 14px', display: 'flex', alignItems: 'center', gap:'8px' }}>
            <div className="av" id="chat-av" style={{ width:'28px', height:'28px', background:'#FBEAF0', color:'#993556', fontSize:'11px' }}>SR</div>
            <div>
              <p id="chat-name" style={{ fontSize:'12px', fontWeight:500, color:'var(--color-text-primary)' }}>Selin Rauf</p>
              <p style={{ fontSize:'10px', color:'var(--color-text-tertiary)' }}><span className="online-dot" style={{ display: 'inline-block', marginRight:'3px', verticalAlign: 'middle' }}></span>online now · websocket</p>
            </div>
          </div>
          <div id="chat-messages" style={{ flex:1, padding:'12px', display: 'flex', flexDirection: 'column', gap:'8px', background:'var(--color-background-tertiary)', minHeight:'340px' }}>
            <div style={{ textAlign: 'center', fontSize:'10px', color:'var(--color-text-tertiary)', padding: '4px 0' }}>today · 14:22</div>
            <div className="msg-bubble-them">hey, tested the websocket handler — looks solid</div>
            <div className="msg-bubble-me">thanks! just added group broadcast support too</div>
            <div className="msg-bubble-them">nice. any issue with the gorilla lib version?</div>
            <div className="msg-bubble-me">nope, 1.5.1 works fine with our go.mod</div>
            <div className="msg-bubble-them">just merged it, take a look 👀</div>
          </div>
          <div style={{ background:'var(--color-background-primary)', borderTop:'0.5px solid var(--color-border-tertiary)', padding:'8px 12px', display: 'flex', gap:'6px', alignItems: 'center' }}>
            <button className="btn btn-g" style={{ fontSize:'14px', padding:'5px 8px', border: 'none' }} title="emoji"><i className="ti ti-mood-smile" aria-hidden="true"></i></button>
            <input className="inp" id="chat-input" style={{ flex:1, fontSize:'12px' }} placeholder="message Selin..." />
            <button className="btn btn-p" style={{ fontSize:'12px', display: 'flex', alignItems: 'center', gap:'3px' }} onClick={ sendMsg() }><i className="ti ti-send" style={{ fontSize:'12px' }} aria-hidden="true"></i></button>
          </div>
        </main>
        <aside className="sidebar2">
          <p className="sec-label" style={{ padding:0, marginBottom:'8px' }}>shared groups</p>
          <div style={{ fontSize:'11px', display: 'flex', flexDirection: 'column', gap:'6px' }}>
            <div style={{ display: 'flex', alignItems: 'center', gap:'6px', padding:'6px', background:'var(--color-background-secondary)', border:'0.5px solid var(--color-border-tertiary)' }}>
              <span style={{ width:'6px', height:'6px', background:'#D4537E', flexShrink:0 }}></span>
              <span style={{ color:'var(--color-text-primary)' }}>go devs</span>
            </div>
            <div style={{ display: 'flex', alignItems: 'center', gap:'6px', padding:'6px', background:'var(--color-background-secondary)', border:'0.5px solid var(--color-border-tertiary)' }}>
              <span style={{ width:'6px', height:'6px', background:'#1D9E75', flexShrink:0 }}></span>
              <span style={{ color:'var(--color-text-primary)' }}>open src</span>
            </div>
          </div>
          <div className="divider"></div>
          <p className="sec-label" style={{ padding:0, marginBottom:'8px' }}>quick emoji</p>
          <div style={{ display: 'flex', flexWrap: 'wrap', gap:'4px' }}>
            <span style={{ fontSize:'16px', cursor: 'pointer' }} title="thumbs up">&'#128077', </span>
            <span style={{ fontSize:'16px', cursor: 'pointer' }}>&'#128078', </span>
            <span style={{ fontSize:'16px', cursor: 'pointer' }}>&'#128514', </span>
            <span style={{ fontSize:'16px', cursor: 'pointer' }}>&#10084, </span>
            <span style={{ fontSize:'16px', cursor: 'pointer' }}>&'#128293', </span>
            <span style={{ fontSize:'16px', cursor: 'pointer' }}>&'#128591', </span>
            <span style={{ fontSize:'16px', cursor: 'pointer' }}>&'#127881', </span>
            <span style={{ fontSize:'16px', cursor: 'pointer' }}>&'#128522', </span>
          </div>
        </aside>
      </div>
    </div>

    <div className="screen" id="screen-notifications">
      <div className="layout">
        <aside className="sidebar">
          <p className="sec-label">navigate</p>
          <div className="navlink" onClick={ showScreen('feed') }><i className="ti ti-home" style={{ fontSize:'14px' }} aria-hidden="true"></i> feed</div>
          <div className="navlink" onClick={ showScreen('profile') }><i className="ti ti-user" style={{ fontSize:'14px' }} aria-hidden="true"></i> profile</div>
          <div className="navlink" onClick={ showScreen('groups') }><i className="ti ti-users" style={{ fontSize:'14px' }} aria-hidden="true"></i> groups</div>
          <div className="navlink" onClick={ showScreen('chat') }><i className="ti ti-message" style={{ fontSize:'14px' }} aria-hidden="true"></i> messages</div>
          <div className="navlink active" onClick={ showScreen('notifications') }><i className="ti ti-bell" style={{ fontSize:'14px' }} aria-hidden="true"></i> notifications</div>
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
          <div style={{ border:'0.5px solid var(--color-border-tertiary)', padding:'8px', background:'var(--color-background-secondary)', marginBottom:'6px', cursor: 'pointer' }} onClick={ showScreen('chat') }>
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
          </div>
          <div style={{ textAlign: 'center', marginTop:'8px' }}><button className="btn btn-p" style={{ fontSize:'10px', width: '100%' }} onClick={ showScreen('chat') }>open messages →</button></div>
        </aside>
      </div>
    </div>

    </div>
    </main>

  );
}

// <main className="flex min-h-screen flex-col items-center justify-center bg-gradient-to-b from-['#2e026d'] to-['#15162c'] text-white">
//   <div className="container flex flex-col items-center justify-center gap-12 px-4 py-16">
//     <h1 className="text-5xl font-extrabold tracking-tight text-white sm:text-[5rem]">
//       Create <span className="text-[hsl(280,100%,70%)]">T3</span> App
//     </h1>
//     <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 md:gap-8">
//       <Link
//         className="flex max-w-xs flex-col gap-4 rounded-xl bg-white/10 p-4 text-white hover:bg-white/20"
//         href="https://create.t3.gg/en/usage/first-steps"
//         target="_blank"
//       >
//         <h3 className="text-2xl font-bold">First Steps →</h3>
//         <div className="text-lg">
//           Just the basics - Everything you need to know to set up your
//           database and authentication.
//         </div>
//       </Link>
//       <Link
//         className="flex max-w-xs flex-col gap-4 rounded-xl bg-white/10 p-4 text-white hover:bg-white/20"
//         href="https://create.t3.gg/en/introduction"
//         target="_blank"
//       >
//         <h3 className="text-2xl font-bold">Documentation →</h3>
//         <div className="text-lg">
//           Learn more about Create T3 App, the libraries it uses, and how to
//           deploy it.
//         </div>
//       </Link>
//     </div>
//   </div>
// </main>
