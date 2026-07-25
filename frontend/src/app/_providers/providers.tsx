"use client";

import { SessionProvider } from "next-auth/react";
import   {ToastProvider} from "~/app/_components/Toast"
import { WSProvider } from "./ws-provider";
import { ChatProvider } from "./chatProvider";

export function Providers({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <SessionProvider>
      <ToastProvider>
        <WSProvider>
          <ChatProvider>
            {children}
          </ChatProvider>
        </WSProvider>
      </ToastProvider>
    </SessionProvider>
  );
}