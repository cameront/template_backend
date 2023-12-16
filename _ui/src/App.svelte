<script lang="ts">
  import { Link, Route, Router } from "svelte-routing";
  import Counter from "./components/Counter.svelte";
  import Login from "./components/Login.svelte";
  import Sidenav from "./components/Sidenav.svelte";
  import { authCookie } from "./lib/stores";
  import { authCookieName } from "./lib/constants";

  $authCookie = document.cookie.split(";").some((c) => {
    return c.trim().startsWith(authCookieName + "=");
  });
</script>

<main>
  <Router>
    <div class="w-5 h-10 border-2 border-red-400">Hello</div>
    <Sidenav />
    <main class="ml-40">
      <Route path="/login" component={Login} />
      <Route path="/" component={Counter} />
    </main>
  </Router>
</main>
