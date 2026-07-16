import Link from "next/link";
import LogoutBtn from "./LogoutBtn";

import { auth } from "~/server/auth";

export default async function Header() {
  const session = await auth();
  console.log("the session",session)
  const avatar = session?.user?.avatar;
  // console.log(avatar)
  return (
    <div className="nav">
      <span className="logo">
        social <span>network</span>
        <span style={{fontSize: '10px', marginLeft: '10px'}} className="tag tag-purple">v1.0</span>
      </span>
      <div style={{ display: 'flex', gap:'10px', alignItems: 'center' }} id="nav-top">
        <input className="inp" style={{ width:'180px',  padding:'4px 8px', fontSize:'11px' }} placeholder="Search users, groups..." />
        <Link className="btn btn-g" style={{ display: 'flex', alignItems: 'center', gap:'4px' }} href="/notifications">
          <i className="ti ti-bell" style={{ fontSize:'14px' }} aria-hidden="true"></i>
          <span className="notif-dot">4</span>
        </Link>
        <Link className="av" /* style={{ width:'28px', height:'28px', background:'#EEEDFE', color:'#534AB7', fontSize:'11px', cursor: 'pointer' }} */ href="/profile">
          {/* AK */}
          <img src={avatar} width={28} height={28} style={{background: '#EEEDFE' }} alt="Avatar" />
        </Link>
        <LogoutBtn />
      </div>
    </div>
  );
}