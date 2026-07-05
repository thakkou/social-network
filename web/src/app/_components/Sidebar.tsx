'use client'; // used only for navigation

import Link from "next/link";
import { usePathname } from 'next/navigation';

const links = [
  { href: '/', label: 'feed', icon: 'ti-home' },
  { href: '/profile', label: 'profile', icon: 'ti-user' },
  { href: '/groups', label: 'groups', icon: 'ti-users' },
  { href: '/messages', label: 'messages', icon: 'ti-message' },
  { href: '/notifications', label: 'notifications', icon: 'ti-bell' },
];

export default function Sidebar() {
  const pathname = usePathname();
  return (
    <aside className="sidebar" id="main-sidebar">
        <p className="sec-label">navigate</p>
        {links.map(link => (
          <Link
            key={link.label}
            className={`navlink ${
              pathname === link.href ||
              (link.href !== '/' && pathname.startsWith(link.href))
                ? 'active'
                : ''
            }`}
            href={link.href}
          >
            <i className={`ti ${link.icon}`} style={{ fontSize:'14px' }} aria-hidden="true"></i> {link.label}
            {link.href === '/messages' || link.href === '/notifications' ? <span className="notif-dot" style={{ marginLeft: 'auto' }}>2</span> : <></> }
          </Link>
        ))}
        <div className="divider" style={{ margin:'0.5rem 0.75rem' }}></div>
        <p className="sec-label">my groups</p>
        <div className="navlink" style={{ fontSize:'11px' }}><span style={{ width:'6px', height:'6px', background:'#D4537E', display: 'inline-block', flexShrink:0 }}></span> go devs</div>
        <div className="navlink" style={{ fontSize:'11px' }}><span style={{ width:'6px', height:'6px', background:'#534AB7', display: 'inline-block', flexShrink:0 }}></span> design sys</div>
        <div className="navlink" style={{ fontSize:'11px' }}><span style={{ width:'6px', height:'6px', background:'#1D9E75', display: 'inline-block', flexShrink:0 }}></span> open src</div>
        <div style={{ padding:'6px 12px', marginTop:'4px' }}><button className="btn btn-g" style={{ width:'100%', fontSize:'11px', display: 'flex', alignItems: 'center', justifyContent: 'center', gap:'4px' }} /* onClick={ showScreen('groups') }*/><i className="ti ti-plus" style={{ fontSize:'12px' }} aria-hidden="true"></i> new group</button></div>
    </aside>
  );
}