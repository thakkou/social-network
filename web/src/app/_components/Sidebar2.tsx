export default function Sidebar2() {
  return (
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
  );
}