export default function PublicLayout({ children }: Readonly<{ children: React.ReactNode }>) {
  return (
    <div className="screen active"> {/* + screen id */}
      <div style={{ display: 'grid', gridTemplateColumns:'1fr 1fr', minHeight:'560px' }}>
        {children}
      </div>
    </div>
  );
}