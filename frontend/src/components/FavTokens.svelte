<script lang="ts">
    import { onMount } from "svelte";
    import {
        favSymbols,
        fetchFavSymbols,
        fetchTrades,
        token,
    } from "../stores/index";
    import type { ApiResponse } from "../types";

    let prices: { symbol: string; price: string }[] = [];
    let symbol: string;
    let amount: number;
    let buying = false;
    const buySymbols = async () => {
        buying = true;
        const response = await fetch("http://localhost:8080/trades/buy", {
            method: "POST",
            headers: { Authorization: `Bearer ${$token}` },
            body: JSON.stringify({ symbol, amount }),
        });
        const { data, status, message }: ApiResponse = await response.json();
        if (status === "success") {
            alert(
                `Congratz you buy ${data.quantity.toFixed(
                    4
                )} of ${symbol} for ${data.price.toFixed(4)}$`
            );
            symbol = undefined;
            amount = undefined;
            fetchTrades();
        } else {
            alert(message);
        }
        buying = false;
    };
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
        fetchFavSymbols().then(getFavPrices);
        const interval = setInterval(() => {
            getFavPrices();
        }, 5 * 1000);
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
    <form on:submit|preventDefault={buySymbols}>
        <fieldset disabled={buying}>
            <div>
                <label for="symbol">Symbol to buy</label>
                <select name="symbol" id="symbol" bind:value={symbol}>
                    <option value="">--- SELECT ---</option>
                    {#each $favSymbols as symbol}
                        <option value={symbol.pair.symbol}
                            >{symbol.pair.symbol}</option
                        >
                    {/each}
                </select>
            </div>
            <div>
                <label for="amount">Amount</label>
                <input
                    type="number"
                    name="amount"
                    id="amount"
                    bind:value={amount}
                />
            </div>
            <div>
                <button class="rounded-sm py-1 px-2" type="submit">BUY!</button>
            </div>
        </fieldset>
    </form>
</div>
