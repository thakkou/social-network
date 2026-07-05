// app/(public)/layout.tsx
export default function PublicLayout({ children }: Readonly<{ children: React.ReactNode }>) {
  return (
    <div style={{ display: 'grid', gridTemplateColumns:'1fr 1fr', minHeight:'560px' }}>
      {children}
    </div>
  );
}

{/* <main className="auth-layout"> */}
