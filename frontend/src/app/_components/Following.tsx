export default function Following() {
    return (
        <div id="profile-following">
            <div className="card">
                <div style={{ display: 'flex', flexDirection: 'column', gap: '10px' }}>
                    <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}><div className="av" style={{ width: '28px', height: '28px', background: '#FBEAF0', color: '#993556', fontSize: '11px' }}>SR</div><span style={{ fontSize: '12px', fontWeight: 500, color: 'var(--color-text-primary)' }}>Selin Rauf</span><button className="btn btn-red" style={{ marginLeft: 'auto', fontSize: '10px' }}>unfollow</button></div>
                    <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}><div className="av" style={{ width: '28px', height: '28px', background: '#FAEEDA', color: '#633806', fontSize: '11px' }}>LK</div><span style={{ fontSize: '12px', fontWeight: 500, color: 'var(--color-text-primary)' }}>Lena K.</span><button className="btn btn-red" style={{ marginLeft: 'auto', fontSize: '10px' }}>unfollow</button></div>
                </div>
            </div>
        </div>
    );
}