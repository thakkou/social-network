"use client";

import { SessionProvider } from "next-auth/react";
import { WSProvider } from "./ws-provider";
import { MessageProvider } from "./message-provider";

export function Providers({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <SessionProvider>
      <WSProvider>
        <MessageProvider>
          {children}
        </MessageProvider>
      </WSProvider>
    </SessionProvider>
  );
}