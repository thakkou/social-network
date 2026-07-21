"use client";

import {
  createContext,
  useContext,
  useState,
  useCallback,
  useEffect,
  type ReactNode,
} from "react";

// ─── Types ───

export type ToastType = "info" | "success" | "error" | "warning";

export interface Toast {
  id: string;
  type: ToastType;
  title: string;
  message?: string;
  duration?: number; // ms, default 4000
}

interface ToastContextValue {
  toasts: Toast[];
  addToast: (t: Omit<Toast, "id">) => string;
  removeToast: (id: string) => void;
}

// ─── Context ───

const ToastContext = createContext<ToastContextValue>({
  toasts: [],
  addToast: () => "",
  removeToast: () => {},
});

export function useToast() {
  return useContext(ToastContext);
}

// ─── Provider ───

export function ToastProvider({ children }: { children: ReactNode }) {
  const [toasts, setToasts] = useState<Toast[]>([]);

  const removeToast = useCallback((id: string) => {
    setToasts((prev) => prev.filter((t) => t.id !== id));
  }, []);

  const addToast = useCallback(
    (t: Omit<Toast, "id">): string => {
      const id = Math.random().toString(36).slice(2, 9);
      const toast: Toast = { ...t, id };
      setToasts((prev) => [...prev, toast]);

      const dur = t.duration ?? 4000;
      if (dur > 0) {
        setTimeout(() => removeToast(id), dur);
      }
      return id;
    },
    [removeToast]
  );

  return (
    <ToastContext.Provider value={{ toasts, addToast, removeToast }}>
      {children}

      {/* ── Toast container ── */}
      <div
        style={{
          position: "fixed",
          top: 16,
          right: 16,
          zIndex: 9999,
          display: "flex",
          flexDirection: "column",
          gap: 8,
          maxWidth: 360,
          pointerEvents: "none",
        }}
      >
        {toasts.map((toast) => (
          <ToastItem key={toast.id} toast={toast} onDismiss={removeToast} />
        ))}
      </div>
    </ToastContext.Provider>
  );
}

// ─── Single toast item ───

const BG_MAP: Record<ToastType, string> = {
  info: "#2e2b27",
  success: "#1a3a2a",
  error: "#3a1a1a",
  warning: "#3a3020",
};

const ICON_MAP: Record<ToastType, string> = {
  info: "ti ti-info-circle",
  success: "ti ti-check-circle",
  error: "ti ti-alert-circle",
  warning: "ti ti-alert-triangle",
};

function ToastItem({
  toast,
  onDismiss,
}: {
  toast: Toast;
  onDismiss: (id: string) => void;
}) {
  useEffect(() => {
    const el = document.getElementById(`toast-${toast.id}`);
    if (el) {
      requestAnimationFrame(() => {
        el.style.opacity = "1";
        el.style.transform = "translateX(0)";
      });
    }
  }, [toast.id]);

  return (
    <div
      id={`toast-${toast.id}`}
      className="card"
      style={{
        padding: "10px 14px",
        display: "flex",
        alignItems: "flex-start",
        gap: 10,
        background: BG_MAP[toast.type],
        border: `1px solid rgba(255,255,255,0.08)`,
        borderRadius: 8,
        opacity: 0,
        transform: "translateX(40px)",
        transition: "opacity 0.25s ease, transform 0.25s ease",
        pointerEvents: "auto",
        cursor: "pointer",
      }}
      onClick={() => onDismiss(toast.id)}
    >
      <i
        className={ICON_MAP[toast.type]}
        style={{ fontSize: 16, marginTop: 1, flexShrink: 0 }}
      />
      <div style={{ flex: 1, minWidth: 0 }}>
        <div style={{ fontSize: 12, fontWeight: 600, marginBottom: 2 }}>
          {toast.title}
        </div>
        {toast.message && (
          <div
            style={{
              fontSize: 11,
              color: "#a09c94",
              lineHeight: 1.4,
              overflow: "hidden",
              textOverflow: "ellipsis",
              whiteSpace: "nowrap",
            }}
          >
            {toast.message}
          </div>
        )}
      </div>
    </div>
  );
}
