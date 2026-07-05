// app/(app)/layout.tsx
// import Sidebar from '@/components/Sidebar';
// import Header from '@/components/Header';

import Sidebar from "../_components/Sidebar";
import Sidebar2 from "../_components/Sidebar2";

export default function AppLayout({ children }: Readonly<{ children: React.ReactNode }>) {
  return (
    <div className="layout">
        <Sidebar />
        {children}
        <Sidebar2 />
    </div>
  );
}

{/* <>
    <Header />
    <Sidebar />
    <main>{children}</main>
</> */}