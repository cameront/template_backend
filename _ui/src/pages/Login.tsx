import { useState } from "react";
import { getRpcHost } from '../lib/rpc_client'
import { setCookie } from "../lib/cookie";
import { useNavigate } from "react-router";

export default function Login() {
  const [error, setError] = useState('');
  const [username, setUsername] = useState('meuser');
  const [password, setPassword] = useState('pass123');
  const navigate = useNavigate();

  async function doLogin(e: React.MouseEvent) {
    e.preventDefault();
    try {
      const response = await fetch(getRpcHost() + "/login", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ username, password }),
        // this must be included to set the cookie returned from the
        // rpc host when it's on a different port than the ui host
        // (i.e. dev environment)
        credentials: "include",
      });
      if (response.status === 200) {
        setCookie();
        navigate("/");
      } else {
        setError(`error: ${response.text}`);
      }
    } catch (err) {
      setError(`${err}`);
    }
  }

  return (
    <>
      <div className="flex justify-center">
        <form className="mt-10">
          <div className="mb-6">
            <label htmlFor="name-input" className="block mb-1">Username</label>
            <input
              className="text-lg p-2 border border-gray-100"
              id="name-input"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
            />
          </div>
          <div className="mb-6">
            <label htmlFor="password-input" className="block mb-1">Password</label>
            <input
              type="password"
              id="password-input"
              value={password}
              className="text-lg p-2 border border-gray-100"
              onChange={(e) => setPassword(e.target.value)}
            />
          </div>
          <button className="w-fit border rounded-lg p-2" color="blue" onClick={(e) => doLogin(e)}>Login</button>
        </form>
      </div>

      {error &&
        <div className="text-center p-5 text-red-700">{error}</div>
      }

      <div className="text-center p-10">Hint - use "meuser" / "pass123"</div>

    </>
  )
}