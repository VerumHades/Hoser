// src/components/Navbar.tsx
import React, { useState } from "react";
import { NavLink } from "react-router-dom";
import user from "./user";

const Navbar: React.FC = () => {
  const [isOpen, setIsOpen] = useState(false);

  const links = [
    {
      name: "Explore",
      to: "/explore"
    },
    {
      name: "Rent",
      to: "/rent"
    },
    {
      name: "Develop",
      to: "/develop"
    },
  ]


  const built_links = links.map(x =>
    <NavLink
      to={x.to}
      end
      className="text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-800 hover:text-gray-900 dark:hover:text-white p-4"
      onClick={() => setIsOpen(false)}
    >
      {x.name}
    </NavLink>
  )

  const account_button =
    <NavLink to={user ? "/account" : "/login"} onClick={() => setIsOpen(false)} 
    className="flex flex-row justify-between p-5 dark:bg-gray-800 shadow-sm m-2 dark:hover:bg-gray-700 hover:text-gray-900">
      <img src="/account_icon.svg"></img>
      {user ? user.Username : "Login"}
    </NavLink>

  return (
    <div className="fixed inset-0 flex flex-col pointer-events-none">
      <nav className="bg-white dark:bg-gray-900 shadow-md w-full pointer-events-auto">
        <div className="flex justify-between h-16 pl-5 pr-5 relative">
          <div className="flex-shrink-0 flex items-center">
            <img src="/ghost_logo.svg" className="h-full pt-3 pb-3" />
            <label className="text-2xl white">gHost</label>
            <div className="h-full justify-center items-center hidden sm:flex">
              {built_links}
            </div>
          </div>

          <div className="justify-center items-center hidden sm:flex">
            {account_button}
          </div>

          <div className="flex items-center sm:h-0 h-full sm:hidden">
            <button
              onClick={() => setIsOpen(!isOpen)}
              className="p-2 focus:outline-none bg-transparent"
            >
              {isOpen ? (
                <span className="text-2xl">✕</span>
              ) : (
                <span className="text-2xl">☰</span>
              )}
            </button>
          </div>
        </div>
      </nav>

      <div className={"top-0 transition-all overflow-hidden pointer-events-auto " + (isOpen ? "flex-1" : "h-0") + 
        " sm:hidden flex flex-col justify-between bg-white dark:bg-gray-900"}>
        <div className="flex flex-col">
          {isOpen ? built_links : <></>}
        </div>
        {account_button}
      </div>
    </div>
  );
};

export default Navbar;