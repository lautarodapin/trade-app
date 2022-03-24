<script lang="ts">
    import { API_URL } from "./../utils";
    import { onMount } from "svelte";
    import { favSymbols, fetchFavSymbols, token } from "../stores/index";
    import type { ApiResponse } from "../types";

    let symbols: { id: number; symbol: string }[] = [];

    onMount(async () => {
        const response = await fetch(`${API_URL}/pair-list/`, {
            headers: {
                Authorization: `Bearer ${$token}`,
            },
        });
        const { data, status, message }: ApiResponse = await response.json();
        symbols = data;
    });
    const addToFav = async (symbol: string) => {
        const response = await fetch(`${API_URL}/pair-list/fav`, {
            method: "POST",
            headers: {
                Authorization: `Bearer ${$token}`,
            },
            body: JSON.stringify({ symbol }),
        });
        const { status, message }: ApiResponse = await response.json();
        if (status === "success") {
            fetchFavSymbols();
        } else {
            alert(message);
        }
    };
</script>

<div>
    <ul>
        {#each symbols as symbol}
            {#if !$favSymbols.some((favSymbol) => favSymbol.pair.symbol === symbol.symbol)}
                <span class="my-2 w-16 h-64">
                    <button
                        on:click={() => addToFav(symbol.symbol)}
                        class="px-4 py-1"
                    >
                        {symbol.symbol}
                        (ADD TO FAV)
                    </button>
                </span>
            {/if}
        {/each}
    </ul>
</div>
