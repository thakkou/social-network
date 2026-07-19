import "~/styles/globals.css";

import { type Metadata } from "next";
import { Geist } from "next/font/google";

import { Providers } from "./_providers/providers";

export const metadata: Metadata = {
  title: "01Social",
  description: "Social network webapp",
  icons: [{ rel: "icon", url: "/favicon.ico" }],
};

const geist = Geist({
  subsets: ["latin"],
  variable: "--font-geist-sans",
});

export default function RootLayout({
  children,
}: Readonly<{ children: React.ReactNode }>) {
  return (
    <html lang="en" className={`${geist.variable}`}>
      <body>
        <main /*style={{height: "100%"}}*/>
          <Providers>
            {children}
          </Providers>
        </main>
      </body>
    </html>
  );
}
