import Link from "next/link";
import LogoutBtn from "./LogoutBtn";
import SearchInput from "./SearchInput"; // Import your new component

import { auth } from "~/server/auth";

export default async function Header() {
  const session = await auth();
  console.log("the session", session);
  
  const avatar = session?.user?.avatar;

  return (
    <div className="nav">
      <span className="logo">
        social <span>network</span>
        <span style={{ fontSize: '10px', marginLeft: '10px' }} className="tag tag-purple">v1.0</span>
      </span>
      <div style={{ display: 'flex', gap: '10px', alignItems: 'center' }} id="nav-top">
        
        {/* Swapped out the raw input for your new component */}
        <SearchInput />

        <Link className="btn btn-g" style={{ display: 'flex', alignItems: 'center', gap: '4px' }} href="/notifications">
          <i className="ti ti-bell" style={{ fontSize: '14px' }} aria-hidden="true"></i>
          <span className="notif-dot">4</span>
        </Link>
        <Link className="av" href="/profile" style={{ width: 28, height: 28, borderRadius: '50%', background: avatar ? 'transparent' : '#EEEDFE', overflow: 'hidden', textDecoration: 'none' }}>
          {avatar ? (
            <img src={avatar} width={28} height={28} alt="Avatar" style={{ objectFit: 'cover' }} />
          ) : (
            <span style={{ color: '#534AB7', fontSize: 11, fontWeight: 600 }}>
              {session?.user?.nickname?.[0]?.toUpperCase() || session?.user?.firstname?.[0]?.toUpperCase() || '?'}
            </span>
          )}
        </Link>
        <LogoutBtn />
      </div>
    </div>
  );
}