import React, { useEffect, useRef } from "react";
import type { Message } from "../types/msg";
import Logo from "../assets/logo.png";

const UserMessage: React.FC<Message> = ({ message }) => {
  return (
    <div className="bg-gray-200 p-2 text-sm sm:text-lg rounded-xl max-w-[70%] my-2 relative font-roboto ml-auto rouned-lg">
      {message}
    </div>
  );
};

const BotMessage: React.FC<Message> = React.memo(({ message }) => {
  return (
    <div className=" bg-transparent p-2 rounded-lg max-w-[70%] my-2 relative flex gap-x-2 items-center font-lora mr-auto">
      <img src={Logo} className="h-8 w-8 rounded-full object-cover" alt="" />
      <span>{message}</span>
    </div>
  );
});

interface msgsProps {
  messages: Message[];
}

export const Messages: React.FC<msgsProps> = ({ messages }) => {
  const chatSectionRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (chatSectionRef.current) {
      chatSectionRef.current.scrollTop = chatSectionRef.current.scrollHeight;
    }
  }, [messages]);

  return (
    <div className=" w-full  flex flex-col items-center">
      <div className=" flex justify-center mb-3">
        <img
          src={Logo}
          alt="logo"
          className="h-24 w-24 rounded-full object-cover"
        />
      </div>
      <div
        ref={chatSectionRef}
        className="flex flex-col relative h-[80vh] overflow-y-scroll pb-[5rem] scroll-smooth w-full  md:w-[100%] lg:w-[50%] sm:px-10"
      >
        {messages.map((msg, i) =>
          msg.isUser ? (
            <UserMessage
              id={msg.id}
              message={msg.message}
              isUser={msg.isUser}
              key={`user-${i}`}
            />
          ) : (
            <BotMessage
              id={msg.id}
              message={msg.message}
              isUser={msg.isUser}
              key={`bot-${i}`}
            />
          )
        )}
      </div>
    </div>
  );
};
