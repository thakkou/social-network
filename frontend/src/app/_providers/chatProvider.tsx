"use client";

import {
  createContext,
  useContext,
  useState,
  ReactNode,
} from "react";

type ChatType = "user" | "group";

type SelectedChat = {
  id: string;
  type: ChatType;
} | null;

type ChatContextType = {
  selectedChat: SelectedChat;
  selectChat: (chat: SelectedChat) => void;
  clearChat: () => void;
};

const ChatContext = createContext<ChatContextType>({
  selectedChat: null,
  selectChat: () => {},
  clearChat: () => {},
});

export function ChatProvider({
  children,
}: {
  children: ReactNode;
}) {
  const [selectedChat, setSelectedChat] =
    useState<SelectedChat>(null);

  function selectChat(chat: SelectedChat) {
    setSelectedChat(chat);
  }

  function clearChat() {
    setSelectedChat(null);
  }

  return (
    <ChatContext.Provider
      value={{
        selectedChat,
        selectChat,
        clearChat,
      }}
    >
      {children}
    </ChatContext.Provider>
  );
}

export function useChat() {
  return useContext(ChatContext);
}