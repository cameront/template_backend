<script lang="ts">
    import { Button, Card, Input, Label } from "flowbite-svelte";
    import { getRpcHost } from "../lib/rpc_client";
    import { authCookie } from "../lib/stores";
    import { navigate } from "svelte-routing";

    let username = "meuser";
    let password = "pass123";
    let error = "";

    async function login() {
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
                $authCookie = true;
                navigate("/");
            } else {
                error = `error: ${response.text}`;
            }
        } catch (err) {
            error = `${err}`;
        }
    }
</script>

<Card class="mx-auto">
    <form>
        <div class="mb-6">
            <Label for="name-input" class="block mb-2">Username</Label>
            <Input id="name-input" size="lg" bind:value={username} />
        </div>
        <div class="mb-6">
            <Label for="password-input" class="block mb-2">Password</Label>
            <Input
                type="password"
                id="password-input"
                size="lg"
                bind:value={password}
            />
        </div>
        <Button class="w-fit" color="blue" on:click={login}>Login</Button>
    </form>
</Card>

{#if error}
    <div class="text-center p-5 text-red-700">{error}</div>
{/if}

<div class="text-center p-10">Hint - use "meuser" / "pass123"</div>
