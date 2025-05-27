<script lang="ts">
    import caretUp from "../assets/caretup.svg";
    import { Link } from "svelte5-router";
    import { authCookie } from "../lib/stores";
    import { onDestroy } from "svelte";
    import { authCookieName } from "../lib/constants";
    import Fingerprint from "../icons/Fingerprint.svelte";

    let hasAuthcookie = false;

    const unsubscribe = authCookie.subscribe((cookieSet) => {
        hasAuthcookie = cookieSet;
    });
    onDestroy(unsubscribe);

    function deleteCookie(name: string) {
        document.cookie =
            authCookieName + "=" + ";expires=Thu, 01 Jan 1970 00:00:01 GMT";
        $authCookie = false;
    }
</script>

<aside
    id="default-sidebar"
    class="fixed top-0 left-0 z-80 w-40 h-screen"
    aria-label="Sidebar"
>
    <div
        class="flex flex-col h-full px-3 py-4 overflow-y-auto bg-gray-50 dark:bg-gray-800"
    >
        <div class="h-10"></div>
        <ul class="space-y-2 font-medium">
            <li>
                <Link
                    id="counter"
                    to="/"
                    class="flex items-center p-2 text-gray-900 rounded-lg dark:text-white hover:bg-gray-100 dark:hover:bg-gray-700 group"
                >
                    <img src={caretUp} alt="caret" width="20" height="20" />
                    <span class="flex-1 ms-3 whitespace-nowrap">Counter</span>
                </Link>
            </li>
            <li>
                <Link
                    id="login"
                    to="/login"
                    class="flex items-center p-2 text-gray-900 rounded-lg dark:text-white hover:bg-gray-100 dark:hover:bg-gray-700 group"
                >
                    <div class="text-gray-500">
                        <Fingerprint width="20" height="20" />
                    </div>
                    <span class="flex-1 ms-3 whitespace-nowrap">Login</span>
                </Link>
            </li>
        </ul>
        {#if $authCookie}
            <div class="flex-grow"></div>
            <ul>
                <li>
                    <Link
                        to="/login"
                        onclick={() => deleteCookie("ac")}
                        class="flex items-center p-2 text-gray-900 rounded-lg dark:text-white hover:bg-gray-100 dark:hover:bg-gray-700 group"
                    >
                        <svg
                            class="w-6 h-6 text-gray-800 dark:text-white"
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
                        <span class="flex-1 ms-3 whitespace-nowrap">Logout</span
                        >
                    </Link>
                </li>
            </ul>
        {/if}
    </div>
</aside>

<style>
</style>
