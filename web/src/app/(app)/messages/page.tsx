export default function Chat() {
  return (
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
        <button className="btn btn-p" style={{ fontSize: '12px', display: 'flex', alignItems: 'center', gap: '3px' }}><i className="ti ti-send" style={{ fontSize:'12px' }} aria-hidden="true"></i></button>
      </div>
    </main>
  );
}

// send message button :
// onClick={sendMsg()}

// function sendMsg() {
//   const inp = document.getElementById('chat-input');
//   const val = inp.value.trim();
//   if (!val) return;
//   const msgs = document.getElementById('chat-messages');
//   const b = document.createElement('div');
//   b.className = 'msg-bubble-me';
//   b.textContent = val;
//   msgs.appendChild(b);
//   inp.value = '';
//   msgs.scrollTop = msgs.scrollHeight;
//   return undefined
// }