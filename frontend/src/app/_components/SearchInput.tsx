"use client";

import { useState, useEffect, useRef } from "react";
import { search } from "../api/search/search";
import Image from "next/image";
import Link from "next/link";

function UserItem({ user }: { user: any }) {


  return (
    <Link
      href={`/profile/${user.id}`}     
       style={{
        display: "flex",
        alignItems: "center",
        gap: "10px",
        padding: "8px 12px",
        cursor: "pointer",
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

function GroupItem({ group }: { group: any }) {

  return (
    <Link
      href={`/groups/${group.id}`}    
        style={{
        display: "flex",
        alignItems: "center",
        gap: "10px",
        padding: "8px 12px",
        cursor: "pointer",
      }}
    >
      <div
        className="av"
        style={{
          width: "28px",
          height: "28px",
          borderRadius: "6px",
          background: "#162820",
          color: "#4dbf95",
        }}
      >
        <i className="ti ti-users" />
      </div>

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


export default function SearchInput() {
  const [query, setQuery] = useState("");

  const debounceTimer = useRef<NodeJS.Timeout | null>(null);

  const [results, setResults] = useState<{
    profiles: any[];
    groups: any[];
  }>({
    profiles: [],
    groups: [],
  });


  const executeSearch = async (searchTerm: string) => {
    if (!searchTerm.trim()) {
      setResults({
        profiles: [],
        groups: [],
      });
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
      }
    } catch (error) {
      console.error("Search failed:", error);

      setResults({
        profiles: [],
        groups: [],
      });
    }
  };


  useEffect(() => {
    if (debounceTimer.current) {
      clearTimeout(debounceTimer.current);
    }

    if (!query.trim()) {
      setResults({
        profiles: [],
        groups: [],
      });
      return;
    }

    debounceTimer.current = setTimeout(() => {
      executeSearch(query);
    }, 1500);


    return () => {
      if (debounceTimer.current) {
        clearTimeout(debounceTimer.current);
      }
    };

  }, [query]);


  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    if (debounceTimer.current) {
      clearTimeout(debounceTimer.current);
    }

    executeSearch(query);
  };


  return (
    <form onSubmit={handleSubmit}>
      <div style={{ position: "relative" }}>

        <input
          className="inp"
          style={{
            width: "180px",
            padding: "4px 8px",
            fontSize: "11px",
          }}
          placeholder="Search users, groups..."
          value={query}
          onChange={(e) => setQuery(e.target.value)}
        />


        {(results.profiles.length > 0 ||
          results.groups.length > 0) && (

          <div
            className="card"
            style={{
              position: "absolute",
              top: "35px",
              left: 0,
              width: "280px",
              zIndex: 100,
              padding: "8px 0",
            }}
          >

            {results.profiles.length > 0 && (
              <>
                <div className="sec-label">
                  USERS
                </div>

                {results.profiles.map((user) => (
                  <UserItem
                    key={user.id}
                    user={user}
                  />
                ))}
              </>
            )}


            {results.groups.length > 0 && (
              <>
                <hr className="divider" />

                <div className="sec-label">
                  GROUPS
                </div>

                {results.groups.map((group) => (
                  <GroupItem
                    key={group.id}
                    group={group}
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