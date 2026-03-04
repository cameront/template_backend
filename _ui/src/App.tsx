import { useEffect, useState } from "react"
import Sidenav from "./components/Sidenav";
import Counter from "./pages/Counter";
import Login from "./pages/Login";
import { userClient } from "./lib/rpc_client";

export default function App() {
  const [loading, setLoading] = useState(false);
  const [authenticated, setAuthenticated] = useState(false);

  console.log('ahhoy/?')
  useEffect(() => {
    async function load() {
      console.log('ahoy');
      setLoading(true);
      try {
        const resp = await userClient.WhoAmI({});
        console.log(resp);
        setAuthenticated(true);
      } catch (err) {
        setAuthenticated(false);
        console.error(err);
      } finally {
        setLoading(false);
      }
    }
    load();
  });

  if (loading) return <div>Loading...</div>

  return (
    <>
      <Sidenav isLoggedIn={authenticated} afterLogout={() => setAuthenticated(false)} />
      <main className="ml-40">
        {authenticated && <Counter />}
        {!authenticated && <Login afterLogin={() => setAuthenticated(true)} />}
      </main>
    </>
  );
}

