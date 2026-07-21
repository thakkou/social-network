"use client";

import { useState, useRef, useEffect } from "react";

const EMOJIS = [
  "😀", "😃", "😄", "😁", "😅", "😂", "🤣", "😊", "😇", "🙂",
  "😉", "😌", "😍", "🥰", "😘", "😗", "😋", "😛", "😜", "🤪",
  "😝", "🤑", "🤗", "🤭", "🤫", "🤔", "🤐", "🤨", "😐", "😑",
  "😶", "😏", "😒", "🙄", "😬", "😮", "😯", "😲", "😳", "🥺",
  "😢", "😭", "😤", "😠", "😡", "🤬", "🤯", "😳", "🥵", "🥶",
  "😱", "😨", "😰", "😥", "😓", "🤗", "🤔", "🤭", "🤫", "🤥",
  "😶", "😐", "😑", "😬", "🙄", "😯", "😦", "😧", "😮", "😲",
  "😴", "🤤", "😪", "😵", "🤐", "🥴", "🤢", "🤮", "🤧", "😷",
  "🤒", "🤕", "🤑", "🤠", "😈", "👿", "👹", "👺", "💀", "👻",
  "👽", "🤖", "💩", "😺", "😸", "😹", "😻", "😼", "😽", "🙀",
  "😿", "😾", "🙈", "🙉", "🙊", "💋", "💌", "💘", "💝", "💖",
  "💗", "💓", "💞", "💕", "💟", "❣️", "💔", "❤️", "🧡", "💛",
  "💚", "💙", "💜", "🤎", "🖤", "🤍", "💯", "💢", "💥", "💫",
  "💦", "💨", "🕳️", "💣", "💬", "👋", "🤚", "🖐️", "✋", "🖖",
  "👌", "🤌", "🤏", "✌️", "🤞", "🤟", "🤘", "🤙", "👈", "👉",
  "👆", "🖕", "👇", "👍", "👎", "✊", "👊", "🤛", "🤜", "👏",
  "🙌", "👐", "🤲", "🤝", "🙏", "✍️", "💅", "🤳", "💪", "🦵",
  "🦶", "👂", "🦻", "👃", "🧠", "🦷", "🦴", "👀", "👁️", "👅",
  "👄", "🎉", "🎊", "🎈", "🎁", "🎀", "🎃", "🎄", "🎆", "🎇",
  "✨", "🎓", "🏆", "🎸", "🎺", "🎻", "🎮", "🎲", "♟️", "🎯",
];

export function EmojiPicker({ onSelect }: { onSelect: (emoji: string) => void }) {
  const [open, setOpen] = useState(false);
  const ref = useRef<HTMLDivElement>(null);

  useEffect(() => {
    function handleClickOutside(e: MouseEvent) {
      if (ref.current && !ref.current.contains(e.target as Node)) {
        setOpen(false);
      }
    }
    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  return (
    <div ref={ref} style={{ position: "relative" }}>
      <button
        className="btn btn-g"
        type="button"
        onClick={() => setOpen(!open)}
        style={{ padding: "5px 8px", fontSize: "14px", lineHeight: 1 }}
        title="Add emoji"
      >
        😊
      </button>
      {open && (
        <div
          style={{
            position: "absolute",
            bottom: "100%",
            right: 0,
            marginBottom: 4,
            width: 240,
            maxHeight: 200,
            overflowY: "auto",
            background: "#272420",
            border: "0.5px solid #3a3733",
            borderRadius: 6,
            padding: 6,
            display: "flex",
            flexWrap: "wrap",
            gap: 2,
            zIndex: 100,
          }}
        >
          {EMOJIS.map((emoji, i) => (
            <button
              key={i}
              type="button"
              onClick={() => {
                onSelect(emoji);
                setOpen(false);
              }}
              style={{
                background: "none",
                border: "none",
                cursor: "pointer",
                fontSize: 18,
                padding: 2,
                width: 30,
                height: 30,
                display: "flex",
                alignItems: "center",
                justifyContent: "center",
                borderRadius: 4,
              }}
              onMouseEnter={(e) => {
                (e.currentTarget as HTMLElement).style.background = "#3a3733";
              }}
              onMouseLeave={(e) => {
                (e.currentTarget as HTMLElement).style.background = "none";
              }}
            >
              {emoji}
            </button>
          ))}
        </div>
      )}
    </div>
  );
}
