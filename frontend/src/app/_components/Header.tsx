import Link from "next/link";
import Image from "next/image";
import LogoutBtn from "./LogoutBtn";
import SearchInput from "./SearchInput";
import NotifCount from "./NotifCount";

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

        <NotifCount />
        <Link className="av" href="/profile" style={{ width: 28, height: 28, borderRadius: '50%', background: avatar ? 'transparent' : '#EEEDFE', overflow: 'hidden', textDecoration: 'none' }}>
          {avatar ? (
            <Image src={avatar} width={28} height={28} alt="Avatar" style={{ objectFit: 'cover' }} />
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