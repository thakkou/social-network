"use client";

import {
  createContext,
  useContext,
  useState,
  ReactNode,
} from "react";


type Message = {
  id?: string;
  senderId: string;
  receiverId: string;
  text: string;
  createdAt?: Date;
};


type MessageContextType = {
  messages: Message[];
  addMessage: (message: Message) => void;
};


const MessageContext = createContext<MessageContextType>({
  messages: [],
  addMessage: () => {},
});

//here i will check each one is user or group and each id or group/user data

export function MessageProvider({
  children,
}: {
  children: ReactNode;
}) {

  const [messages, setMessages] = useState<Message[]>([]);


  function addMessage(message: Message) {
    setMessages((prev) => [
      ...prev,
      message,
    ]);
  }


  return (
    <MessageContext.Provider
      value={{
        messages,
        addMessage,
      }}
    >
      {children}
    </MessageContext.Provider>
  );
}


export function useMessages() {
  return useContext(MessageContext);
}