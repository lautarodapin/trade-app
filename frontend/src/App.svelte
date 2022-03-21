<script lang="ts">
    import "./global_styles.css";
    import type { ApiResponse } from "./types";
    import { onMount } from "svelte";
    import Login from "./components/Login.svelte";
    import { logout, token, user } from "./stores/index";
    import SymbolList from "./components/SymbolList.svelte";
    import FavSymbols from "./components/FavSymbols.svelte";
    import TradesTable from "./components/TradesTable.svelte";

    token.subscribe(async (token) => {
        if (token) {
            const response = await fetch(
                "http://localhost:8080/users/current",
                {
                    headers: {
                        Authorization: `Bearer ${token}`,
                    },
                }
            );
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
        <h1 class="text-5xl font-bold flex-row flex justify-between">
            Welcome {$user?.email}
            <button class="py-2 px-4 rounded-md" on:click={logout}
                >Logout</button
            >
        </h1>
        <div class="border border-red-700">
            <SymbolList />
        </div>
        <div class="border border-blue-700">
            <FavSymbols />
        </div>
        <div class="border border-green-600">
            <TradesTable />
        </div>
    {:else}
        <Login />
    {/if}
</main>
