export default function ChatPage() {
  return (
    <div className="screen active" id="screen-chat">
      <div className="layout">
        <aside className="sidebar">
          <p className="sec-label">navigate</p>
          {/* <div className="navlink" onClick={ showScreen('feed') }><i className="ti ti-home" style={{ fontSize:'14px' }} aria-hidden="true"></i> feed</div>
          <div className="navlink" onClick={ showScreen('profile') }><i className="ti ti-user" style={{ fontSize:'14px' }} aria-hidden="true"></i> profile</div>
          <div className="navlink" onClick={ showScreen('groups') }><i className="ti ti-users" style={{ fontSize:'14px' }} aria-hidden="true"></i> groups</div>
          <div className="navlink active" onClick={ showScreen('chat') }><i className="ti ti-message" style={{ fontSize:'14px' }} aria-hidden="true"></i> messages</div>
          <div className="navlink" onClick={ showScreen('notifications') }><i className="ti ti-bell" style={{ fontSize:'14px' }} aria-hidden="true"></i> notifications</div> */}
          <div className="divider" style={{ margin: '0.5rem 0.75rem' }}></div>
          <p className="sec-label">direct</p>
          {/* <div className="navlink active-chat" style={{ background:'#FBEAF0', borderLeft:'2px solid #D4537E', color:'#D4537E', fontSize:'11px' }} onClick={ switchChat('selin') }>
            <div className="av" style={{ width:'20px', height:'20px', background:'#FBEAF0', color:'#993556', fontSize:'9px' }}>SR</div> Selin R. <span className="notif-dot" style={{ marginLeft: 'auto' }}>2</span>
          </div> */}
          {/* <div className="navlink" style={{ fontSize:'11px' }} onClick={ switchChat('jonas') }>
            <div className="av" style={{ width:'20px', height:'20px', background:'#E1F5EE', color:'#0F6E56', fontSize:'9px' }}>JM</div> Jonas M.
          </div> */}
          <div className="divider" style={{ margin: '0.5rem 0.75rem' }}></div>
          <p className="sec-label">groups</p>
          {/* <div className="navlink" style={{ fontSize:'11px' }} onClick={ switchChat('godevs') }>
            <span style={{ width:'6px', height:'6px', background:'#D4537E', display: 'inline-block', flexShrink:0 }}></span> go devs
          </div> */}
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
            {/* <button className="btn btn-p" style={{ fontSize:'12px', display: 'flex', alignItems: 'center', gap:'3px' }} onClick={ sendMsg() }><i className="ti ti-send" style={{ fontSize:'12px' }} aria-hidden="true"></i></button> */}
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
  );
}