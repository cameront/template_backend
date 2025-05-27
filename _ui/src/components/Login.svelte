<script lang="ts">
    import { getRpcHost } from "../lib/rpc_client";
    import { authCookie } from "../lib/stores";
    import { navigate } from "svelte5-router";

    let username = $state("meuser");
    let password = $state("pass123");
    let error = $state("");

    async function login(e: Event) {
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

<div class="flex justify-center">
    <form class="mt-10">
        <div class="mb-6">
            <label for="name-input" class="block mb-2">Username</label>
            <input
                class="text-lg p-2 border border-gray-100"
                id="name-input"
                bind:value={username}
            />
        </div>
        <div class="mb-6">
            <label for="password-input" class="block mb-2">Password</label>
            <input
                type="password"
                id="password-input"
                class="text-lg p-2 border border-gray-100"
                bind:value={password}
            />
        </div>
        <button class="w-fit" color="blue" onclick={login}>Login</button>
    </form>
</div>

{#if error}
    <div class="text-center p-5 text-red-700">{error}</div>
{/if}

<div class="text-center p-10">Hint - use "meuser" / "pass123"</div>
