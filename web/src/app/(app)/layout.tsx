import Header from "../_components/Header";
import Sidebar from "../_components/Sidebar";
import Sidebar2 from "../_components/Sidebar2";

export default function AppLayout({ children }: Readonly<{ children: React.ReactNode }>) {
  return (
    <>
      <Header />
      <div className="screen active"> {/* + screen id */}
        <div className="layout">
            <Sidebar />
            {children}
            <Sidebar2 />
        </div>
      </div>
    </>
  );
}