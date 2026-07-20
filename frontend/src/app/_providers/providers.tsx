"use client";

import { SessionProvider } from "next-auth/react";
import { WSProvider } from "./ws-provider";
import { ChatProvider } from "./chatProvider";

export function Providers({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <SessionProvider>
      <WSProvider>
        <ChatProvider>
          {children}
        </ChatProvider>
      </WSProvider>
    </SessionProvider>
  );
}