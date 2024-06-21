import { FaChevronUp } from "react-icons/fa";
import Logo from "../assets/logo.png";
import { useEffect, useRef, useState } from "react";
import { Message } from "../types/msg";
import { Messages } from "./Message";

const SSE_URL = "http://127.0.0.1:2323/messages";

export const Chat = ({ sidebar }: { sidebar: boolean }) => {
  const [messages, setMessages] = useState<Message[]>([]);
  const [disabled, setDisabled] = useState(false);
  const [input, setInput] = useState("");
  const [token, setToken] = useState("");
  // const [messageResponseChunk, setMessageResponseChunk] = useState("");
  // const [totalChunks, setTotalChunks] = useState<string[]>([]);

  const messageIdRef = useRef(0);

  useEffect(() => {
    const eventSource = new EventSource(SSE_URL);
    eventSource.onmessage = (event) => {
      console.log(event.data);
    };

    eventSource.addEventListener("message", (event) => {
      const resp = JSON.parse(event.data);
      const botMessage: Message = {
        id: ++messageIdRef.current,
        message: resp.data,
        isUser: false,
      };

      setMessages((prevMsgs) => [...prevMsgs, botMessage]);
    });

    eventSource.addEventListener("connected", function (event) {
      const parsedData = JSON.parse(event.data);
      setToken(parsedData.data);
      console.log(parsedData);
    });

    eventSource.onerror = () => {
      eventSource.close();
      return;
    };

    eventSource.onopen = () => {
      console.log("CONNECTION ESTABLISHED");
    };
  }, []);

  const handleChat = async (e: React.FormEvent) => {
    e.preventDefault();
    setDisabled(true);

    try {
      if (input.length === 0) return;
      const newMessage: Message = {
        id: ++messageIdRef.current,
        message: input,
        isUser: true,
      };

      setMessages((prevMsgs) => [...prevMsgs, newMessage]);
      setInput("");

      const resp = await fetch("http://localhost:2323/chat", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          prompt: input,
        }),
      });
      const data = await resp.json();
      console.log(data);
    } catch (err) {
      console.log(err);
    } finally {
      setDisabled(false);
    }
  };

  return (
    <div className="font-lora p-5 ">
      {messages.length > 0 ? (
        <Messages messages={messages} />
      ) : (
        <div className="flex flex-col items-center justify-center -top-16 gap-y-6 overflow-y-scroll h-screen relative">
          <img
            src={Logo}
            alt="logo"
            className="h-24 w-24 object-cover rounded-full "
          />
          <h1 className=" text-text md:w-1/2 text-center font-semibold sm:text-2xl font-montserrat">
            Get directions to your destination with LocateX. Enter your location
          </h1>
        </div>
      )}

      <div
        className={`absolute bottom-0 backdrop-blur-sm pb-5 transition-all ${
          sidebar ? "left-[10%]" : "left-0"
        } w-full flex items-center justify-center`}
      >
        <form
          className="flex items-center justify-center mt-5 relative w-[80%] sm:w-1/2"
          onSubmit={handleChat}
        >
          {/* <button
            type="button"
            onClick={() => {
              console.log(totalChunks.join(""));
            }}
          >
            Check
          </button> */}
          <textarea
            className="p-2 border border-gray-300 rounded-lg w-full pr-[3.6rem] outline-accent text-sm sm:text-base text-mutualText resize-none min-h-[40px] max-h-[200px] overflow-y-auto flex items-center"
            cols={30}
            placeholder="Enter your location"
            value={input}
            onChange={(e) => setInput(e.target.value)}
          ></textarea>
          <button
            type="submit"
            className={` ${
              disabled ? "bg-lightAccent" : "bg-accent"
            }  h-full px-5 text-white absolute right-0 rounded-r-lg`}
            disabled={disabled}
          >
            <FaChevronUp />
          </button>
        </form>
      </div>
    </div>
  );
};
