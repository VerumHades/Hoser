import React from "react";

interface DropdownProps {
    children: React.ReactNode
    title: string
}


export default function HoverDropdown({ children, title }: DropdownProps) {
  return (
    <div className="relative inline-block text-left group">
      {/* Trigger Button */}
      <button className="inline-flex items-center justify-center w-full rounded-md bg-white dark:bg-gray-800 text-gray-800 dark:text-gray-100 font-medium px-4 py-2 shadow-sm ring-1 ring-gray-300 dark:ring-gray-700 hover:bg-gray-50 dark:hover:bg-gray-700 transition">
        {title}
        <svg
          className="w-4 h-4 ml-2"
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
        >
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
        </svg>
      </button>

      {/* Dropdown Menu */}
      <div className="absolute left-0 mt-2 w-40 origin-top-left scale-95 opacity-0 group-hover:opacity-100 group-hover:scale-100 transform transition-all duration-200 ease-out z-10">
        <div className="rounded-md shadow-lg bg-white dark:bg-gray-800 ring-1 ring-black/5 dark:ring-gray-700">
          <ul className="py-1 text-sm text-gray-700 dark:text-gray-200">
            {children ? React.Children.toArray(children).map(child => {
                return <li className="block px-4 py-2 hover:bg-gray-100 dark:hover:bg-gray-700 transition">
                {child}
                </li>
            }) : <></>}
          </ul>
        </div>
      </div>
    </div>
  );
};