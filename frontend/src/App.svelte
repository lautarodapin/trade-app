<script lang="ts">
    import type { ApiResponse } from "./types";
    import { onMount } from "svelte";
    import Login from "./components/Login.svelte";
    import { token, user } from "./stores/index";
    import TokenList from "./components/TokenList.svelte";
    import FavTokens from "./components/FavTokens.svelte";
    import TradesTable from "./components/TradesTable.svelte";

    token.subscribe(async (token) => {
        if (token) {
            const response = await fetch("http://localhost:8080/user", {
                headers: {
                    Authorization: `Bearer ${token}`,
                },
            });
            const { data, status, message }: ApiResponse =
                await response.json();
            if (status === "success") {
                user.set(data);
            } else {
                alert(message);
            }
        }
    });
</script>

<main>
    {#if $token}
        <div class="border border-red-700">
            <TokenList />
        </div>
        <div class="border border-blue-700">
            <FavTokens />
        </div>
        <div class="border border-green-600">
            <TradesTable />
        </div>
    {:else}
        <Login />
    {/if}
</main>

<style global lang="postcss">
    @tailwind base;
    @tailwind components;
    @tailwind utilities;
</style>
