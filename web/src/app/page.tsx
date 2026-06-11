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
    <div className="screen active" id="screen-feed">
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
