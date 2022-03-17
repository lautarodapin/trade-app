<script lang="ts">
    import { onMount } from "svelte";
    import { fetchTrades, trades } from "../stores/index";
    import { TradeType } from "../types";

    onMount(() => {
        fetchTrades();
    });
</script>

<div>
    <table class="m-4">
        <tr>
            <td class="px-4">ID</td>
            <td class="px-4">Type</td>
            <td class="px-4">Quantity</td>
            <td class="px-4">Price</td>
            <td class="px-4">Earns</td>
        </tr>
        <tbody>
            {#each $trades as trade}
                <tr>
                    <td class="text-center my-2">{trade.ID}</td>
                    <td class="text-center my-2"
                        >{trade.type === TradeType.BUY ? "Buy" : "Sell"}</td
                    >
                    <td class="text-center my-2">{trade.quantity.toFixed(4)}</td
                    >
                    <td class="text-center my-2">{trade.price.toFixed(4)}</td>
                    <td class="text-center my-2">{trade.earns.toFixed(4)}</td>
                </tr>
            {/each}
        </tbody>
    </table>
</div>
