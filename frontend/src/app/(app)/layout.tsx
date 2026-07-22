import Header from "../_components/Header";
import Sidebar from "../_components/Sidebar";
import MobileNav from "../_components/MobileNav";

export default function AppLayout({ children }: Readonly<{ children: React.ReactNode }>) {
  return (
    <>
      <Header />
      <MobileNav />
      <div className="screen active">
        <div className="layout">
            <Sidebar />
            {children}
        </div>
      </div>
    </>
  );
}
