import { Link } from "react-router";
import { ChevronUp, FingerprintPattern } from 'lucide-react';
import { deleteCookie } from "../lib/cookie";

export default function Sidenav() {
  const authCookie = '';

  return (
    <aside
      id="default-sidebar"
      className="fixed top-0 left-0 z-80 w-40 h-screen"
      aria-label="Sidebar"
    >
      <div
        className="flex flex-col h-full px-3 py-4 overflow-y-auto bg-gray-50 dark:bg-gray-800"
      >
        <div className="h-10"></div>
        <ul className="space-y-2 font-medium">
          <li>
            <Link
              id="counter"
              to="/"
              className="flex items-center p-2 text-gray-900 rounded-lg dark:text-white hover:bg-gray-100 dark:hover:bg-gray-700 group"
            >
              <ChevronUp />
              <span className="flex-1 ms-3 whitespace-nowrap">Counter</span>
            </Link>
          </li>
          <li>
            <Link
              id="login"
              to="/login"
              className="flex items-center p-2 text-gray-900 rounded-lg dark:text-white hover:bg-gray-100 dark:hover:bg-gray-700 group"
            >
              <div className="text-gray-500">
                <FingerprintPattern width="20" height="20" />
              </div>
              <span className="flex-1 ms-3 whitespace-nowrap">Login</span>
            </Link>
          </li>
        </ul>
        {authCookie && <>
          <div className="flex-grow"></div>
          <ul>
            <li>
              <Link
                to="/login"
                onClick={() => deleteCookie("ac", () => { })}
                className="flex items-center p-2 text-gray-900 rounded-lg dark:text-white hover:bg-gray-100 dark:hover:bg-gray-700 group"
              >
                <svg
                  className="w-6 h-6 text-gray-800 dark:text-white"
                  aria-hidden="true"
                  xmlns="http://www.w3.org/2000/svg"
                  fill="none"
                  viewBox="0 0 16 16"
                >
                  <path
                    stroke="currentColor"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M4 8h11m0 0-4-4m4 4-4 4m-5 3H3a2 2 0 0 1-2-2V3a2 2 0 0 1 2-2h3"
                  />
                </svg>
                <span className="flex-1 ms-3 whitespace-nowrap">Logout</span>
              </Link>
            </li>
          </ul>
        </>
        }
      </div>
    </aside>
  );
}
