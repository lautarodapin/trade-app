<script lang="ts">
    import { onMount } from "svelte";
    import { fetchTrades, token, trades } from "../stores/index";
    import { TradeType, type ApiResponse } from "../types";

    let unrealized_pl: number = 0;
    let net_pnl: number = 0;
    let cumulative_pl = 0;

    const getEarns = async () => {
        const response = await fetch("http://localhost:8080/trades/earns", {
            headers: { Authorization: `Bearer ${$token}` },
        });
        const { data, status, message }: ApiResponse = await response.json();
        if (status === "success") {
            cumulative_pl = data.cumulative_pl;
            net_pnl = data.net_pnl;
            unrealized_pl = data.unrealized_pl;
        } else {
            alert(message);
        }
    };
    onMount(() => {
        fetchTrades().then(getEarns);
        let timer = setInterval(() => {
            getEarns();
        }, 1000 * 15);
        return () => clearInterval(timer);
    });
</script>

<div class="inline-flex">
    <table class="m-4">
        <thead>
            <tr>
                <th class="px-4">ID</th>
                <th class="px-4">Symbol</th>
                <th class="px-4">Type</th>
                <th class="px-4">Quantity</th>
                <th class="px-4">Price</th>
                <th class="px-4">Earns</th>
            </tr>
        </thead>
        <tbody>
            {#each $trades as trade}
                <tr>
                    <td class="text-center my-2">{trade.ID}</td>
                    <td class="text-center my-2">{trade.pair.symbol}</td>
                    <td class="text-center my-2"
                        >{trade.type === TradeType.BUY ? "Buy" : "Sell"}</td
                    >
                    <td class="text-center my-2">{trade.quantity}</td>
                    <td class="text-center my-2">{trade.price}</td>
                    <td class="text-center my-2">{trade.earns}</td>
                </tr>
            {/each}
        </tbody>
    </table>
</div>
<div class="inline-flex">
    <table>
        <thead>
            <th colspan="2">Earns</th>
        </thead>
        <tbody>
            <tr>
                <th>Unrealized PL</th>
                <td>{unrealized_pl}</td>
            </tr>
            <tr>
                <th>Cumulative PL</th>
                <td>{cumulative_pl}</td>
            </tr>
            <tr>
                <th>Net PNL</th>
                <td>{net_pnl}</td>
            </tr>
        </tbody>
    </table>
</div>
