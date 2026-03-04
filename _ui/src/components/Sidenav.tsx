import { Link } from "react-router";
import { ChevronUp, FingerprintPattern } from 'lucide-react';
import { userClient } from "../lib/rpc_client";

export default function Sidenav({ isLoggedIn, afterLogout }: { isLoggedIn: boolean, afterLogout: () => void }) {
  async function logout() {
    try {
      await userClient.Logout({});
      afterLogout();
    } catch (err) {
      console.error(err);
    }
  }


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
          {isLoggedIn && <><li>
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
                id="logout"
                to="/"
                onClick={() => logout()}
                className="flex items-center p-2 text-gray-900 rounded-lg dark:text-white hover:bg-gray-100 dark:hover:bg-gray-700 group"
              >
                <div className="text-gray-500">
                  <FingerprintPattern width="20" height="20" />
                </div>
                <span className="flex-1 ms-3 whitespace-nowrap">Logout</span>
              </Link>
            </li>
          </>
          }
          {!isLoggedIn &&
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
            </li>}
        </ul>
      </div>
    </aside>
  );
}
