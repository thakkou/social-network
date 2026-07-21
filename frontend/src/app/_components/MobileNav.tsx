"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";

const links = [
  { href: "/", label: "feed", icon: "ti-home" },
  { href: "/profile", label: "profile", icon: "ti-user" },
  { href: "/groups", label: "groups", icon: "ti-users" },
  { href: "/messages", label: "messages", icon: "ti-message" },
  { href: "/notifications", label: "notifs", icon: "ti-bell" },
  { href: "/settings", label: "settings", icon: "ti-settings" },
];

export default function MobileNav() {
  const pathname = usePathname();

  const isActive = (href: string) => {
    if (href === "/") return pathname === href;
    return pathname?.startsWith(href);
  };

  return (
    <nav className="mobile-bottom-nav">
      {links.map((link) => (
        <Link
          key={link.href}
          href={link.href}
          className={isActive(link.href) ? "active" : ""}
        >
          <i className={`ti ${link.icon}`} aria-hidden="true" />
          <span>{link.label}</span>
        </Link>
      ))}
    </nav>
  );
}
