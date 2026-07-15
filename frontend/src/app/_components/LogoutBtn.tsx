"use client";

import { signOut } from "next-auth/react";

export default function LogoutBtn() {
  return (
    <button className="btn btn-p" onClick={() => signOut({ callbackUrl: "/login" })}>
        Log out
    </button>
  );
}