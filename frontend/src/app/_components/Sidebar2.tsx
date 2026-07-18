'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';

// homepage: notifications + online now
// profile: follow requests + send follow
// notifications: messages
// messages: shared groups + quick emoji
// groups: pending requests + create event

import type React from "react";

const HomepageSidebar: React.ComponentType<any> = () =>
  (<aside className="sidebar2">
    {/* NOTIFICATIONS */}
    <p className="sec-label" style={{ padding: 0, marginBottom: '8px' }}>notifications</p>
    <div className="notif-item" style={{ background: '#FBEAF0', borderLeftColor: '#D4537E' }}><span style={{ fontWeight: 500, color: '#72243E' }}>Lena K.</span> <span style={{ color: '#993556' }}>sent a follow request</span>
      <div style={{ display: 'flex', gap: '4px', marginTop: '5px' }}><button className="btn btn-t" style={{ fontSize: '10px', padding: '3px 8px' }}>accept</button><button className="btn btn-red" style={{ fontSize: '10px', padding: '3px 8px' }}>decline</button></div>
    </div>
    <div className="notif-item" style={{ background: '#EEEDFE', borderLeftColor: '#7F77DD' }}><span style={{ fontWeight: 500, color: '#3C3489' }}>design sys</span> <span style={{ color: '#534AB7' }}>invited you</span>
      <div style={{ display: 'flex', gap: '4px', marginTop: '5px' }}><button className="btn btn-t" style={{ fontSize: '10px', padding: '3px 8px' }}>accept</button><button className="btn btn-red" style={{ fontSize: '10px', padding: '3px 8px' }}>decline</button></div>
    </div>
    <div className="notif-item" style={{ background: '#E1F5EE', borderLeftColor: '#1D9E75' }}><span style={{ color: '#085041' }}>event created in <span style={{ fontWeight: 500 }}>go devs</span></span></div>

    <div className="divider"></div>

    {/* ONLINE NOW */}
    <p className="sec-label" style={{ padding: 0, marginBottom: '8px' }}>online now</p>
    <div style={{ display: 'flex', flexDirection: 'column', gap: '7px' }}>
      <div style={{ display: 'flex', alignItems: 'center', gap: '7px' }}>
        <div className="av" style={{ width: '26px', height: '26px', background: '#FBEAF0', color: '#993556', fontSize: '10px' }}>SR</div>
        <span style={{ fontSize: '11px', color: 'var(--color-text-primary)' }}>Selin R.</span>
        <span className="online-dot" style={{ marginLeft: 'auto' }}></span>
      </div>
      <div style={{ display: 'flex', alignItems: 'center', gap: '7px' }}>
        <div className="av" style={{ width: '26px', height: '26px', background: '#E1F5EE', color: '#0F6E56', fontSize: '10px' }}>JM</div>
        <span style={{ fontSize: '11px', color: 'var(--color-text-primary)' }}>Jonas M.</span>
        <span className="online-dot" style={{ marginLeft: 'auto' }}></span>
      </div>
      <div style={{ display: 'flex', alignItems: 'center', gap: '7px' }}>
        <div className="av" style={{ width: '26px', height: '26px', background: '#EEEDFE', color: '#534AB7', fontSize: '10px' }}>OS</div>
        <span style={{ fontSize: '11px', color: 'var(--color-text-primary)' }}>Omar S.</span>
        <span className="offline-dot" style={{ marginLeft: 'auto' }}></span>
      </div>
    </div>
  </aside>);

const ProfileSidebar: React.ComponentType<any> = () =>
  (<aside className="sidebar2">
    <p className="sec-label" style={{ padding: 0, marginBottom: '8px' }}>follow requests</p>
    <div style={{ background: '#FBEAF0', border: '0.5px solid #ED93B1', padding: '8px', marginBottom: '6px', fontSize: '11px' }}>
      <div style={{ display: 'flex', alignItems: 'center', gap: '6px', marginBottom: '6px' }}>
        <div className="av" style={{ width: '22px', height: '22px', background: '#FBEAF0', color: '#993556', fontSize: '10px' }}>TP</div>
        <span style={{ color: '#72243E', fontWeight: 500 }}>Talia P.</span>
        <span style={{ color: '#993556', fontSize: '10px' }}>wants to follow you</span>
      </div>
      <div style={{ display: 'flex', gap: '4px' }}><button className="btn btn-t" style={{ fontSize: '10px', padding: '3px 8px' }}>accept</button><button className="btn btn-red" style={{ fontSize: '10px', padding: '3px 8px' }}>decline</button></div>
    </div>
    <div className="divider"></div>
    <p className="sec-label" style={{ padding: 0, marginBottom: '8px' }}>send follow</p>
    <div style={{ fontSize: '11px', color: 'var(--color-text-secondary)', marginBottom: '6px' }}>find a user and send a follow request</div>
    <input className="inp" style={{ fontSize: '11px', marginBottom: '6px' }} placeholder="search by name or email..." />
    <div style={{ border: '0.5px solid var(--color-border-tertiary)', padding: '7px', display: 'flex', alignItems: 'center', gap: '6px', background: 'var(--color-background-secondary)' }}>
      <div className="av" style={{ width: '24px', height: '24px', background: '#E1F5EE', color: '#0F6E56', fontSize: '10px' }}>JM</div>
      <span style={{ fontSize: '11px', color: 'var(--color-text-primary)' }}>Jonas M.</span>
      <button className="btn btn-p" style={{ marginLeft: 'auto', fontSize: '10px', padding: '3px 8px' }}>+ follow</button>
    </div>
  </aside>);

const NotificationsSidebar: React.ComponentType<any> = () =>
  (<aside className="sidebar2">
    <p className="sec-label" style={{ padding: 0, marginBottom: '8px' }}>messages</p>
    <div style={{ fontSize: '11px', color: 'var(--color-text-secondary)', marginBottom: '8px' }}>private messages are separate from notifications</div>
    <Link style={{ border: '0.5px solid var(--color-border-tertiary)', padding: '8px', background: 'var(--color-background-secondary)', marginBottom: '6px', cursor: 'pointer' }} href="/chat">
      <div style={{ display: 'flex', alignItems: 'center', gap: '6px', marginBottom: '3px' }}>
        <div className="av" style={{ width: '20px', height: '20px', background: '#FBEAF0', color: '#993556', fontSize: '9px' }}>SR</div>
        <span style={{ fontSize: '11px', fontWeight: 500, color: 'var(--color-text-primary)' }}>Selin R.</span>
        <span className="notif-dot" style={{ marginLeft: 'auto' }}>2</span>
      </div>
      <p style={{ fontSize: '11px', color: 'var(--color-text-secondary)', paddingLeft: '26px' }}>just merged it, take a look</p>
    </Link>
    <Link style={{ border: '0.5px solid var(--color-border-tertiary)', padding: '8px', background: 'var(--color-background-secondary)', cursor: 'pointer' }} href="/chat">
      <div style={{ display: 'flex', alignItems: 'center', gap: '6px', marginBottom: '3px' }}>
        <div className="av" style={{ width: '20px', height: '20px', background: '#E1F5EE', color: '#0F6E56', fontSize: '9px' }}>JM</div>
        <span style={{ fontSize: '11px', fontWeight: 500, color: 'var(--color-text-primary)' }}>Jonas M.</span>
      </div>
      <p style={{ fontSize: '11px', color: 'var(--color-text-secondary)', paddingLeft: '26px' }}>great work on the auth layer</p>
    </Link>
    <div style={{ textAlign: 'center', marginTop: '8px' }}><Link className="btn btn-p" style={{ fontSize: '10px', width: '100%' }} href="/chat">open messages →</Link></div>
  </aside>);

const MessagesSidebar: React.ComponentType<any> = () =>
  (<aside className="sidebar2">
    <p className="sec-label" style={{ padding: 0, marginBottom: '8px' }}>shared groups</p>
    <div style={{ fontSize: '11px', display: 'flex', flexDirection: 'column', gap: '6px' }}>
      <div style={{ display: 'flex', alignItems: 'center', gap: '6px', padding: '6px', background: 'var(--color-background-secondary)', border: '0.5px solid var(--color-border-tertiary)' }}>
        <span style={{ width: '6px', height: '6px', background: '#D4537E', flexShrink: 0 }}></span>
        <span style={{ color: 'var(--color-text-primary)' }}>go devs</span>
      </div>
      <div style={{ display: 'flex', alignItems: 'center', gap: '6px', padding: '6px', background: 'var(--color-background-secondary)', border: '0.5px solid var(--color-border-tertiary)' }}>
        <span style={{ width: '6px', height: '6px', background: '#1D9E75', flexShrink: 0 }}></span>
        <span style={{ color: 'var(--color-text-primary)' }}>open src</span>
      </div>
    </div>
    <div className="divider"></div>
    <p className="sec-label" style={{ padding: 0, marginBottom: '8px' }}>quick emoji</p>
    <div style={{ display: 'flex', flexWrap: 'wrap', gap: '4px' }}>
      <span style={{ fontSize: '16px', cursor: 'pointer' }} title="thumbs up">&'#128077', </span>
      <span style={{ fontSize: '16px', cursor: 'pointer' }}>&'#128078', </span>
      <span style={{ fontSize: '16px', cursor: 'pointer' }}>&'#128514', </span>
      <span style={{ fontSize: '16px', cursor: 'pointer' }}>&#10084, </span>
      <span style={{ fontSize: '16px', cursor: 'pointer' }}>&'#128293', </span>
      <span style={{ fontSize: '16px', cursor: 'pointer' }}>&'#128591', </span>
      <span style={{ fontSize: '16px', cursor: 'pointer' }}>&'#127881', </span>
      <span style={{ fontSize: '16px', cursor: 'pointer' }}>&'#128522', </span>
    </div>
  </aside>);

const GroupsSidebar: React.ComponentType<any> = () =>
  (<aside className="sidebar2">
    <p className="sec-label" style={{ padding: 0, marginBottom: '8px' }}>pending requests</p>
    <div style={{ background: '#EEEDFE', border: '0.5px solid #AFA9EC', padding: '8px', fontSize: '11px', marginBottom: '6px' }}>
      <p style={{ fontWeight: 500, color: '#3C3489', marginBottom: '2px' }}>Omar S.</p>
      <p style={{ color: '#534AB7', marginBottom: '5px' }}>wants to join go devs</p>
      <div style={{ display: 'flex', gap: '4px' }}><button className="btn btn-t" style={{ fontSize: '10px', padding: '3px 8px' }}>accept</button><button className="btn btn-red" style={{ fontSize: '10px', padding: '3px 8px' }}>decline</button></div>
    </div>
    <div className="divider"></div>
    <p className="sec-label" style={{ padding: 0, marginBottom: '8px' }}>create event</p>
    <div className="form-row"><span className="form-label">title</span><input className="inp" style={{ fontSize: '11px' }} placeholder="event title" /></div>
    <div className="form-row"><span className="form-label">description</span><textarea className="inp" style={{ height: '44px', resize: 'none', fontSize: '11px' }} placeholder="details..."></textarea></div>
    <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '6px', marginBottom: '8px' }}>
      <div className="form-row"><span className="form-label">date</span><input className="inp" type="date" style={{ fontSize: '11px' }} /></div>
      <div className="form-row"><span className="form-label">time</span><input className="inp" type="time" style={{ fontSize: '11px' }} /></div>
    </div>
    <button className="btn btn-p" style={{ width: '100%', fontSize: '11px' }}>create event →</button>
  </aside>);




const DefaultSidebar: React.ComponentType<any> = () => (
  <aside className="sidebar2">
    <p className="sec-label" style={{ padding: 0, marginBottom: "8px" }}>
      quick access
    </p>

    <div
      style={{
        fontSize: "11px",
        color: "var(--color-text-secondary)",
        marginBottom: "10px",
      }}
    >
      explore the app
    </div>

    <Link
      href="/"
      className="navlink"
    >
      <i className="ti ti-home" />
      home
    </Link>

    <Link
      href="/profile"
      className="navlink"
    >
      <i className="ti ti-user" />
      profile
    </Link>

    <Link
      href="/groups"
      className="navlink"
    >
      <i className="ti ti-users" />
      groups
    </Link>

    <Link
      href="/messages"
      className="navlink"
    >
      <i className="ti ti-message" />
      messages
    </Link>

    <Link
      href="/notifications"
      className="navlink"
    >
      <i className="ti ti-bell" />
      notifications
    </Link>

    <div className="divider"></div>

    <p className="sec-label" style={{ padding: 0, marginBottom: "8px" }}>
      status
    </p>

    <div
      style={{
        display: "flex",
        alignItems: "center",
        gap: "8px",
        fontSize: "11px",
      }}
    >
      <span className="online-dot"></span>
      online
    </div>
  </aside>
);

const navbars = [
  { path: '/', Component: HomepageSidebar },
  { path: '/profile', Component: ProfileSidebar },
  { path: '/groups', Component: GroupsSidebar },
  { path: '/messages', Component: MessagesSidebar },
  { path: '/notifications', Component: NotificationsSidebar },
];

export default function Sidebar2() {
  const pathname = usePathname();

  const CustomNavbar =
    navbars.find((nv) => nv.path === pathname)?.Component ||
    DefaultSidebar;

  return <CustomNavbar />;
}