<script lang="ts">
    import { onMount } from "svelte";
    import { favSymbols, fetchFavSymbols, token } from "../stores/index";
    import type { ApiResponse } from "../types";

    let prices: { symbol: string; price: string }[] = [];

    const getFavPrices = async () => {
        if ($favSymbols.length === 0) return;
        const response = await fetch(
            "http://localhost:8080/pair-list/fav/prices",
            {
                headers: {
                    Authorization: `Bearer ${$token}`,
                },
            }
        );
        const { data, status, message }: ApiResponse = await response.json();
        if (status === "success") {
            prices = data;
        } else {
            alert(message);
        }
    };

    onMount(() => {
        const interval = setInterval(() => {
            getFavPrices();
        }, 5 * 1000);
        fetchFavSymbols();
        return () => clearInterval(interval);
    });

    const removeFromFav = async (id: number) => {
        const response = await fetch(
            "http://localhost:8080/pair-list/fav/" + id.toString(),
            {
                method: "DELETE",
                headers: {
                    Authorization: `Bearer ${$token}`,
                },
            }
        );
        const { status, message }: ApiResponse = await response.json();
        if (status === "success") {
            alert(message);
            fetchFavSymbols();
        } else {
            alert(message);
        }
    };
</script>

<div class="m-10">
    <h1 class="">Fav symbols</h1>
    <ul>
        {#each $favSymbols as favSymbol}
            <li class="my-2">
                {favSymbol.pair.symbol}
                <button
                    class="ml-2 rounded-sm py-1 px-2"
                    on:click={() => {
                        removeFromFav(favSymbol.ID);
                    }}
                >
                    Remove from fav {favSymbol.ID}
                </button>
                Price: {prices.find(
                    ({ symbol }) => symbol === favSymbol.pair.symbol
                )?.price || "N/A"}
            </li>
        {/each}
    </ul>
</div>
