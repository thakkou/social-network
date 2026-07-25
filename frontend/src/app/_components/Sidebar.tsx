'use client';

import Link from "next/link";
import { usePathname } from 'next/navigation';
import { useState, useEffect } from "react";
import { useSession } from "next-auth/react";
import { getUserGroups } from "~/app/_services/crud/groups";
import { getProfileData } from "~/app/_services/crud/getProfile";

const links = [
  { href: '/', label: 'feed', icon: 'ti-home' },
  { href: '/profile', label: 'profile', icon: 'ti-user' },
  { href: '/groups', label: 'groups', icon: 'ti-users' },
  { href: '/messages', label: 'messages', icon: 'ti-message' },
  { href: '/notifications', label: 'notifications', icon: 'ti-bell' },
  { href: '/settings', label: 'settings', icon: 'ti-settings' },
];

const GROUP_COLORS = ["#D4537E", "#534AB7", "#1D9E75", "#E28743", "#7F77DD", "#3B82F6"];

function colorFor(id: number) {
  return GROUP_COLORS[Math.abs(id) % GROUP_COLORS.length];
}

type MyGroup = { id: number; title: string };
type FollowingUser = { id: number; nickname: string; firstname: string; lastname: string; avatar: string };

export default function Sidebar() {
  const pathname = usePathname();
  const { data: session } = useSession();
  const userId = session?.user?.id;

  const [myGroups, setMyGroups] = useState<MyGroup[]>([]);
  const [following, setFollowing] = useState<FollowingUser[]>([]);
  const [followers, setFollowers] = useState<FollowingUser[]>([]);
  const [loadingGroups, setLoadingGroups] = useState(true);
  const [loadingDirect, setLoadingDirect] = useState(true);

  useEffect(() => {
    if (!userId) return;
    const load = async () => {
      setLoadingGroups(true);
      const res = await getUserGroups();
      if (res.success) setMyGroups(res.data ?? []);
      setLoadingGroups(false);
    };
    void load();
  }, [userId]);

  useEffect(() => {
    if (!userId) return;
    const load = async () => {
      setLoadingDirect(true);
      const res = await getProfileData(userId);
      if (res.success) {
        if (res.data?.following) setFollowing(res.data.following ?? []);
        if (res.data?.followers) setFollowers(res.data.followers ?? []);
      }
      setLoadingDirect(false);
    };
    void load();
  }, [userId]);

  return (
    <aside className="sidebar" id="main-sidebar">
      {/* NAVIGATION */}
      <p className="sec-label">navigate</p>
      {links.map(link => (
        <Link
          key={link.label}
          className={`navlink ${
            pathname === link.href ||
            (link.href !== '/' && pathname?.startsWith(link.href))
            ? 'active'
            : ''
          }`}
          href={link.href}
        >
          <i className={`ti ${link.icon}`} style={{ fontSize:'14px' }} aria-hidden="true"></i> {link.label}
        </Link>
      ))}

      <div className="divider" style={{ margin:'0.5rem 0.75rem' }}></div>
      
      {/* MY GROUPS */}
      <p className="sec-label">my groups</p>
      {loadingGroups ? (
        <div className="navlink" style={{ fontSize: '10px', color: '#6b6760' }}>loading...</div>
      ) : (myGroups?.length ?? 0) === 0 ? (
        <div className="navlink" style={{ fontSize: '10px', color: '#6b6760' }}>no groups yet</div>
      ) : (
        myGroups?.map(g => (
          <Link
            key={g?.id ?? 0}
            href={`/groups/${g?.id ?? 0}`}
            className="navlink"
            style={{ fontSize: '11px', textDecoration: 'none' }}
          >
            <span style={{ width:'6px', height:'6px', background: colorFor(g?.id ?? 0), display: 'inline-block', flexShrink:0 }}></span>
            {g?.title ?? ''}
          </Link>
        ))
      )}
      <div style={{ padding:'6px 12px', marginTop:'4px' }}>
        <Link className="btn btn-g" style={{ width:'100%', fontSize:'11px', display: 'flex', alignItems: 'center', justifyContent: 'center', gap:'4px' }} href="/groups">
          <i className="ti ti-plus" style={{ fontSize:'12px' }} aria-hidden="true"></i> new group
        </Link>
      </div>

      <div className="divider" style={{ margin: '0.5rem 0.75rem' }}></div>

      {/* DIRECT - who I follow */}
      <p className="sec-label">following ({following?.length || 0})</p>
      {loadingDirect ? (
        <div className="navlink" style={{ fontSize: '10px', color: '#6b6760' }}>loading...</div>
      ) : (following?.length ?? 0) === 0 ? (
        <div className="navlink" style={{ fontSize: '10px', color: '#6b6760' }}>not following anyone yet</div>
      ) : (
        following?.map(u => (
          <Link
            key={u?.id ?? 0}
            href={`/profile/${u?.id ?? 0}`}
            className="navlink"
            style={{ fontSize: '11px', textDecoration: 'none' }}
          >
            <div className="av" style={{ width:'20px', height:'20px', background: colorFor(u?.id ?? 0), color: '#fff', fontSize:'9px', display:'flex', alignItems:'center', justifyContent:'center', flexShrink: 0 }}>
              {u?.nickname
                ? (u.nickname[0]?.toUpperCase() ?? '?')
                : (u?.firstname?.[0]?.toUpperCase() ?? '?')
              }
            </div>
            {u?.nickname || `${u?.firstname ?? ''} ${u?.lastname ?? ''}`.trim()}
          </Link>
        ))
      )}

      <div className="divider" style={{ margin: '0.5rem 0.75rem' }}></div>

      {/* DIRECT - followers */}
      <p className="sec-label">followers ({followers?.length || 0})</p>
      {loadingDirect ? (
        <div className="navlink" style={{ fontSize: '10px', color: '#6b6760' }}>loading...</div>
      ) : (followers?.length ?? 0) === 0 ? (
        <div className="navlink" style={{ fontSize: '10px', color: '#6b6760' }}>no followers yet</div>
      ) : (
        followers?.map(u => (
          <Link
            key={u?.id ?? 0}
            href={`/profile/${u?.id ?? 0}`}
            className="navlink"
            style={{ fontSize: '11px', textDecoration: 'none' }}
          >
            <div className="av" style={{ width:'20px', height:'20px', background: colorFor(u?.id ?? 0), color: '#fff', fontSize:'9px', display:'flex', alignItems:'center', justifyContent:'center', flexShrink: 0 }}>
              {u?.nickname
                ? (u.nickname[0]?.toUpperCase() ?? '?')
                : (u?.firstname?.[0]?.toUpperCase() ?? '?')
              }
            </div>
            {u?.nickname || `${u?.firstname ?? ''} ${u?.lastname ?? ''}`.trim()}
          </Link>
        ))
      )}
    </aside>
  );
}
