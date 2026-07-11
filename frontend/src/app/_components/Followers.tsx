export default function Followers() {
    return (
        <div id="profile-followers">
            <div className="card">
                <div style={{ display: 'flex', flexDirection: 'column', gap: '10px' }}>
                    <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}><div className="av" style={{ width: '28px', height: '28px', background: '#FBEAF0', color: '#993556', fontSize: '11px' }}>SR</div><span style={{ fontSize: '12px', fontWeight: 500, color: 'var(--color-text-primary)' }}>Selin Rauf</span><button className="btn btn-g" style={{ marginLeft: 'auto', fontSize: '10px' }}>view profile</button></div>
                    <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}><div className="av" style={{ width: '28px', height: '28px', background: '#E1F5EE', color: '#0F6E56', fontSize: '11px' }}>JM</div><span style={{ fontSize: '12px', fontWeight: 500, color: 'var(--color-text-primary)' }}>Jonas M.</span><button className="btn btn-g" style={{ marginLeft: 'auto', fontSize: '10px' }}>view profile</button></div>
                    <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}><div className="av" style={{ width: '28px', height: '28px', background: '#EEEDFE', color: '#534AB7', fontSize: '11px' }}>OS</div><span style={{ fontSize: '12px', fontWeight: 500, color: 'var(--color-text-primary)' }}>Omar S.</span><button className="btn btn-g" style={{ marginLeft: 'auto', fontSize: '10px' }}>view profile</button></div>
                </div>
            </div>
        </div>
    );
}