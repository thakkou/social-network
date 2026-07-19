"use client";

import {
  createContext,
  useContext,
  useEffect,
  useState,
  ReactNode,
} from "react";

type WSContextType = {
  socket: WebSocket | null;
  connected: boolean;
};

const WSContext = createContext<WSContextType>({
  socket: null,
  connected: false,
});

export function WSProvider({ children }: { children: ReactNode }) {
  const [socket, setSocket] = useState<WebSocket | null>(null);
  const [connected, setConnected] = useState(false);

  useEffect(() => {
    const ws = new WebSocket("ws://localhost:8080");

    ws.onopen = () => {
      console.log("WS connected");
      setConnected(true);
    };

    ws.onclose = () => {
      console.log("WS disconnected");
      setConnected(false);
    };

    setSocket(ws);

    return () => {
      ws.close();
    };
  }, []);

  return (
    <WSContext.Provider
      value={{
        socket,
        connected,
      }}
    >
      {children}
    </WSContext.Provider>
  );
}


export function useWS() {
  return useContext(WSContext);
}