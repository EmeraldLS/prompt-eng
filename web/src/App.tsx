import { useState } from "react";
import { Chat } from "./components/Chat";
import { Sidebar } from "./components/Sidebar";
import { FaChevronRight } from "react-icons/fa6";

function App() {
  const [showSidebar, setShowSidebar] = useState(false);
  const handleToggler = () => {
    setShowSidebar(!showSidebar);
  };

  return (
    <>
      <div className="flex">
        <div
          className={`${
            showSidebar ? "w-1/5" : "w-0"
          } sticky hidden sm:flex justify-between transition-all duration-300 ease-linear overflow-hidden `}
        >
          <Sidebar />
        </div>

        <button
          className="h-screen hidden sm:flex items-center justify-center transition-all ease-in-out duration-300 hover:bg-gray-50"
          onClick={handleToggler}
        >
          <FaChevronRight
            className={`transition-all ease-in-out duration-300 ${
              showSidebar ? "rotate-180" : ""
            }`}
          />
        </button>

        <div className="flex-1 h-screen overflow-y-hidden">
          <Chat sidebar={showSidebar} />
        </div>
      </div>
    </>
  );
}

export default App;
