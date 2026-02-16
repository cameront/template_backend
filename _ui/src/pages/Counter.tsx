import { useEffect, useState } from "react";
import { counterClient } from "../lib/rpc_client";

export default function Counter() {
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const [count, setCount] = useState(0);
  const [counterName, setCounterName] = useState('');
  const [tmpCounterName, setTmpCounterName] = useState('');

  useEffect(() => {
    async function load() {
      setLoading(true);
      try {
        const val = await counterClient.GetValue({ name: counterName });
        setCount(Number(val.value));
      } catch (err) {
        setError(`error: ${err}`);
      } finally {
        setLoading(false);
      }
    }
    load();
  }, [counterName])

  async function updateCounter() {
    setCounterName(tmpCounterName);
  }

  async function incrementCounter() {
    try {
      const res = await counterClient.Increment({ name: counterName, amount: BigInt(1) });
      setCount(Number(res.value));
    } catch (err) {
      setError(`error: ${err}`);
    }
  }

  return (
    <div className="flex justify-center mt-10">
      {error !== "" &&
        <div className="max-w-[300px] error">{error}</div>
      }
      {loading &&
        <div>Loading...</div>
      }
      {!error && !loading &&
        <div>
          <div className="mb-6">
            <label htmlFor="name-input" className="block mb-2">Counter Name</label>
            <div className="flex flex-row gap-4">
              <input
                id="name-input"
                className="text-lg p-2 border border-gray-100"
                value={tmpCounterName}
                onChange={(e) => setTmpCounterName(e.target.value)}
              />
              <button
                className="border p-2 rounded-md hover:bg-gray-50"
                onClick={() => updateCounter()}>Set</button>
            </div>
          </div>
          <div className="text-center">
            <h5
              className="mb-2 text-2xl font-bold tracking-tight text-gray-900 dark:text-white"
            >
              Value: {count}
            </h5>

            <button color="green" className="w-fit" onClick={() => incrementCounter()}
            >Increment</button>
          </div>
        </div>
      }
    </div>

  );
}

