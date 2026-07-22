"use client";

interface ConfirmModalProps {
  open: boolean;
  title: string;
  message: string;
  confirmLabel?: string;
  cancelLabel?: string;
  confirmClass?: string;
  onConfirm: () => void;
  onCancel: () => void;
}

export default function ConfirmModal({
  open,
  title,
  message,
  confirmLabel = "delete",
  cancelLabel = "cancel",
  confirmClass = "btn-red",
  onConfirm,
  onCancel,
}: ConfirmModalProps) {
  if (!open) return null;

  return (
    <div
      style={{
        position: "fixed",
        inset: 0,
        zIndex: 9999,
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
        background: "rgba(0,0,0,0.6)",
        backdropFilter: "blur(4px)",
      }}
      onClick={onCancel}
    >
      <div
        className="card"
        style={{
          width: "90%",
          maxWidth: "380px",
          padding: "1.25rem",
          border: "0.5px solid #3a3733",
        }}
        onClick={(e) => e.stopPropagation()}
      >
        <p
          style={{
            fontSize: "13px",
            fontWeight: 600,
            color: "#e8e4dc",
            marginBottom: "8px",
          }}
        >
          {title}
        </p>
        <p
          style={{
            fontSize: "12px",
            color: "#a09c94",
            lineHeight: 1.5,
            marginBottom: "16px",
          }}
        >
          {message}
        </p>
        <div
          style={{
            display: "flex",
            gap: "6px",
            justifyContent: "flex-end",
          }}
        >
          <button
            className="btn btn-g"
            style={{ fontSize: "11px" }}
            onClick={onCancel}
          >
            {cancelLabel}
          </button>
          <button
            className={`btn ${confirmClass}`}
            style={{ fontSize: "11px" }}
            onClick={onConfirm}
          >
            {confirmLabel}
          </button>
        </div>
      </div>
    </div>
  );
}
