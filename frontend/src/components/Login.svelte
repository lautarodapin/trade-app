<script lang="ts">
    import { login } from "../stores/index";
    import type { ApiResponse } from "../types";

    let email = "";
    let password = "";
    const submitForm = async () => {
        const response = await fetch("http://localhost:8080/login", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            mode: "cors",
            body: JSON.stringify({
                email,
                password,
            }),
        });
        const { data, message, status }: ApiResponse = await response.json();
        if (status === "success") {
            login(data.token);
        } else {
            alert(message);
        }
    };
</script>

<div>
    <form on:submit|preventDefault={submitForm}>
        <div>
            <label for="email">Email</label>
            <input type="text" name="email" id="email" bind:value={email} />
        </div>
        <div>
            <label for="password">Password</label>
            <input
                type="password"
                name="password"
                id="password"
                bind:value={password}
            />
        </div>
        <button type="submit">Login</button>
    </form>
</div>
