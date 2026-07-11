export default function ProfilePosts() {
    return (
        <div id="profile-posts">
            <div className="card" style={{ marginBottom: '6px' }}>
                <p style={{ fontSize: '10px', color: 'var(--color-text-tertiary)', marginBottom: '4px' }}>2 days ago · <span className="tag tag-teal">public</span></p>
                <p style={{ fontSize: '12px', color: 'var(--color-text-primary)', lineHeight: 1.5 }}>Just deployed the migration system — golang-migrate running smooth with our SQLite schema. All 6 tables up.</p>
            </div>
            <div className="card">
                <p style={{ fontSize: '10px', color: 'var(--color-text-tertiary)', marginBottom: '4px' }}>5 days ago · <span className="tag tag-gray">followers</span></p>
                <p style={{ fontSize: '12px', color: 'var(--color-text-primary)', lineHeight: 1.5 }}>Bcrypt session cookie auth is done. Logout clears the cookie and invalidates the server-side session.</p>
            </div>
        </div>
    );
}