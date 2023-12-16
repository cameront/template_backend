<script lang="ts">
  import { client } from "../lib/rpc_client";
  import { Button, Card, Input, Label } from "flowbite-svelte";

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

<div class="max-w-sm mx-auto">
  {#if error !== ""}
    <div class="error">{error}</div>
  {:else if loading}
    <div>Loading...</div>
  {:else}
    <Card>
      <div class="mb-6">
        <Label for="name-input" class="block mb-2">Counter Name</Label>
        <div class="flex flex-row gap-4">
          <Input id="name-input" size="lg" bind:value={tmpCounterName} />
          <Button color="light" on:click={setAndFetch}>Set</Button>
        </div>
      </div>
      <div class="text-center">
        <h5
          class="mb-2 text-2xl font-bold tracking-tight text-gray-900 dark:text-white"
        >
          Value: {count}
        </h5>

        <Button color="green" class="w-fit" on:click={incrementCounter}
          >Increment</Button
        >
      </div>
    </Card>
  {/if}
</div>

<style>
</style>
