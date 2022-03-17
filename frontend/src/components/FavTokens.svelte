<script lang="ts">
    import { onMount } from "svelte";
    import { favSymbols, fetchFavSymbols, token } from "../stores/index";
    import type { ApiResponse } from "../types";

    onMount(() => {
        fetchFavSymbols();
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

<div>
    <ul>
        {#each $favSymbols as favSymbol}
            <li>
                {favSymbol.pair.symbol}
                <button
                    on:click={() => {
                        removeFromFav(favSymbol.ID);
                    }}
                >
                    Remove from fav {favSymbol.ID}
                </button>
            </li>
        {/each}
    </ul>
</div>
