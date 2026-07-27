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

  // Track container position to decide dropdown alignment
  const [alignRight, setAlignRight] = useState(false);
  useEffect(() => {
    function updateAlignment() {
      if (!ref.current) return;
      const rect = ref.current.getBoundingClientRect();
      // If button is in the right half of the viewport, align dropdown to the right
      setAlignRight(rect.left > window.innerWidth / 2);
    }
    if (open) {
      updateAlignment();
      window.addEventListener("resize", updateAlignment);
      return () => window.removeEventListener("resize", updateAlignment);
    }
  }, [open]);

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
            marginBottom: 6,
            width: 320,
            maxWidth: "calc(100vw - 24px)",
            maxHeight: 280,
            overflowY: "auto",
            background: "#272420",
            border: "0.5px solid #3a3733",
            borderRadius: 8,
            padding: 8,
            display: "flex",
            flexWrap: "wrap",
            gap: 2,
            zIndex: 200,
            boxShadow: "0 8px 24px rgba(0,0,0,0.4)",
            // Align based on button position
            ...(alignRight
              ? { right: 0 }
              : { left: 0 }
            ),
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
                fontSize: 20,
                padding: 3,
                width: 36,
                height: 36,
                display: "flex",
                alignItems: "center",
                justifyContent: "center",
                borderRadius: 6,
                transition: "background 0.1s ease",
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
