"use client";

import { useState, useEffect, useRef } from "react";
import { search } from "../api/crud/search";
import Image from "next/image";
import Link from "next/link";
interface SearchInputProps {
  typeSearch?: "all" | "users" | "groups";
  placeholder?: string;
  style?: React.CSSProperties;
  className?: string;
}

function UserItem({ user, onClick }: { user: any; onClick: () => void }) {
  return (
    <Link
      href={`/profile/${user.id}`}
      onClick={onClick}
      style={{
        display: "flex",
        alignItems: "center",
        gap: "10px",
        padding: "8px 12px",
        cursor: "pointer",
        textDecoration: "none",
        color: "inherit",
      }}
    >
      {/* Avatar */}
      {user.avatar ? (
        <Image
          src={user.avatar}
          alt="avatar"
          width={28}
          height={28}
          className="av"
          style={{
            borderRadius: "50%",
            objectFit: "cover",
          }}
        />
      ) : (
        <div
          className="av"
          style={{
            width: "28px",
            height: "28px",
            borderRadius: "50%",
            background: "#D4537E",
            color: "#fff",
            fontSize: "12px",
            display: "flex",
            alignItems: "center",
            justifyContent: "center",
          }}
        >
          {(user.nickname?.[0] || user.firstname?.[0] || "?").toUpperCase()}
        </div>
      )}

      <div>
        <div
          style={{
            fontSize: "12px",
            display: "flex",
            alignItems: "center",
            gap: "5px",
          }}
        >
          {user.nickname || `${user.firstname} ${user.lastname}`}
          <span title={user.is_private === 1 ? "Private profile" : "Public profile"}>
            {user.is_private === 1 ? "🔒" : ""}
          </span>
        </div>

        <div
          style={{
            fontSize: "10px",
            color: "#6b6760",
          }}
        >
          User
        </div>
      </div>
    </Link>
  );
}

function GroupItem({ group, onClick }: { group: any; onClick: () => void }) {
 

  return (
    <Link
      href={`/groups/${group.id}`}
      onClick={onClick}
      style={{
        display: "flex",
        alignItems: "center",
        gap: "10px",
        padding: "8px 12px",
        cursor: "pointer",
        textDecoration: "none",
        color: "inherit",
      }}
    >
      {group.logo ? (
        <Image
          src={group.logo}
          alt={group.title}
          width={28}
          height={28}
          style={{
            width: "28px",
            height: "28px",
            borderRadius: "6px",
            objectFit: "cover",
          }}
        />
      ) : (
        <div
          className="av"
          style={{
            width: "28px",
            height: "28px",
            borderRadius: "6px",
            background: "#162820",
            color: "#4dbf95",
            display: "flex",
            alignItems: "center",
            justifyContent: "center",
          }}
        >
          <i className="ti ti-users" />
        </div>
      )}

      <div>
        <div style={{ fontSize: "12px" }}>
          {group.title}
        </div>

        <div
          style={{
            fontSize: "10px",
            color: "#6b6760",
          }}
        >
          Group
        </div>
      </div>
    </Link>
  );
}

export default function SearchInput({
  typeSearch = "all",
  placeholder = "Search users, groups...",
  style,
  className,
}: SearchInputProps) {
  const [query, setQuery] = useState("");
  const [isOpen, setIsOpen] = useState(false);
  const containerRef = useRef<HTMLDivElement>(null);
  const debounceTimer = useRef<NodeJS.Timeout | null>(null);

  const [results, setResults] = useState<{
    profiles: any[];
    groups: any[];
  }>({
    profiles: [],
    groups: [],
  });

  const profiles = typeSearch === "groups" ? [] : results.profiles;
  const groups = typeSearch === "users" ? [] : results.groups;

  const executeSearch = async (searchTerm: string) => {
    if (!searchTerm.trim()) {
      setResults({ profiles: [], groups: [] });
      setIsOpen(false);
      return;
    }

    try {
      const response = await search(searchTerm.trim());

      if (response.success) {
        const data = response.data.data;

        setResults({
          profiles: data?.profiles ?? [],
          groups: data?.groups ?? [],
        });
        setIsOpen(true);
      }
    } catch (error) {
      console.error("Search failed:", error);
      setResults({ profiles: [], groups: [] });
      setIsOpen(false);
    }
  };

  useEffect(() => {
    if (debounceTimer.current) {
      clearTimeout(debounceTimer.current);
    }

    if (!query.trim()) {
      setResults({ profiles: [], groups: [] });
      setIsOpen(false);
      return;
    }

    debounceTimer.current = setTimeout(() => {
      executeSearch(query);
    }, 500);

    return () => {
      if (debounceTimer.current) {
        clearTimeout(debounceTimer.current);
      }
    };
  }, [query]);

  // Handle click outside to close dropdown
  useEffect(() => {
    const handleClickOutside = (e: MouseEvent) => {
      if (
        containerRef.current &&
        !containerRef.current.contains(e.target as Node)
      ) {
        setIsOpen(false);
      }
    };

    document.addEventListener("mousedown", handleClickOutside);
    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
    };
  }, []);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    if (debounceTimer.current) {
      clearTimeout(debounceTimer.current);
    }

    executeSearch(query);
  };

  const handleSelect = () => {
    setIsOpen(false);
    setQuery("");
  };

  const hasResults = profiles.length > 0 || groups.length > 0;

  return (
    <form  onSubmit={handleSubmit} style={style}>
      <div  ref={containerRef} style={{ position: "relative" }}>
        <input
          className="inp"
          
          placeholder={placeholder}
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          onFocus={() => {
            if (hasResults) setIsOpen(true);
          }}
        />

        {isOpen && hasResults && (
          <div
            className="card"
            style={{
              position: "absolute",
              top: "calc(100% + 4px)",
              left: 0,
              width: "260px",
              maxHeight: "300px",
              overflowY: "auto",
              zIndex: 9999,
              padding: "8px 0",
            }}
          >
            {profiles.length > 0 && (
              <>
                <div className="sec-label">USERS</div>
                {profiles.map((user) => (
                  <UserItem
                    key={user.id}
                    user={user}
                    onClick={handleSelect}
                  />
                ))}
              </>
            )}

            {profiles.length > 0 && groups.length > 0 && (
              <hr className="divider" />
            )}

            {groups.length > 0 && (
              <>
                <div className="sec-label">GROUPS</div>
                {groups.map((group) => (
                  <GroupItem
                    key={group.id}
                    group={group}
                    onClick={handleSelect}
                  />
                ))}
              </>
            )}
          </div>
        )}
      </div>
    </form>
  );
}