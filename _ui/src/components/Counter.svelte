<script lang="ts">
  import { client } from "../lib/rpc_client";

  let counterName = "";
  let tmpCounterName = "default";
  let loading = false;
  let error = "";
  let count: number = 0;

  async function fetchCounter() {
    loading = true;
    error = "";
    try {
      const resp = await client.getValue({
        name: counterName,
      });
      count = Number(resp.response.value);
    } catch (err) {
      error = `Error fetching counter: ${err}`;
    }
    loading = false;
  }

  async function incrementCounter() {
    loading = true;
    error = "";
    try {
      const resp = await client.increment({
        name: counterName,
        amount: BigInt(1),
      });
      count = Number(resp.response.value);
    } catch (err) {
      error = `Error incrementing counter: ${err}`;
    }
    loading = false;
  }

  function setAndFetch() {
    counterName = tmpCounterName;
    fetchCounter();
  }

  setAndFetch();
</script>

<div class="flex justify-center mt-10">
  {#if error !== ""}
    <div class="max-w-[300px] error">{error}</div>
  {:else if loading}
    <div>Loading...</div>
  {:else}
    <div>
      <div class="mb-6">
        <label for="name-input" class="block mb-2">Counter Name</label>
        <div class="flex flex-row gap-4">
          <input
            id="name-input"
            class="text-lg p-2 border border-gray-100"
            bind:value={tmpCounterName}
          />
          <button
            class="border p-2 rounded-md hover:bg-gray-50"
            onclick={setAndFetch}>Set</button
          >
        </div>
      </div>
      <div class="text-center">
        <h5
          class="mb-2 text-2xl font-bold tracking-tight text-gray-900 dark:text-white"
        >
          Value: {count}
        </h5>

        <button color="green" class="w-fit" onclick={incrementCounter}
          >Increment</button
        >
      </div>
    </div>
  {/if}
</div>

<style>
  .error {
    color: red;
  }
</style>
