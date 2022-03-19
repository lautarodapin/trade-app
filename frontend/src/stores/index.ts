import { writable } from 'svelte/store';
import type { User, ApiResponse, FavPair, Trade } from '../types';

// const fakeToken = "9ae4:9c47:a59f:9427:bc36:f6ec:536f:3c83"
const fakeToken = undefined
export const token = writable(fakeToken || localStorage.getItem('token') || '');

export const login = (newToken: string) => {
    localStorage.setItem('token', newToken);
    token.set(newToken);
}

export const logout = () => {
    localStorage.removeItem('token');
    token.set('');
}

export const user = writable<User | null>(null)
export const favSymbols = writable<FavPair[]>([]);
export const fetchFavSymbols = async () => {
    const token = localStorage.getItem('token');
    const response = await fetch("http://localhost:8080/pair-list/fav", {
        headers: {
            Authorization: `Bearer ${token}`,
        },
    });
    const { data, status, message }: ApiResponse = await response.json();
    if (status === "success") {
        favSymbols.set(data);
    } else {
        alert(message);
        favSymbols.set([]);
    }
}

export const trades = writable<Trade[]>([]);
export const fetchTrades = async () => {
    const token = localStorage.getItem('token');
    const response = await fetch("http://localhost:8080/trades/", {
        headers: { Authorization: `Bearer ${token}` },
    });
    const { data, status, message }: ApiResponse = await response.json();
    if (status === "success") {
        trades.set(data);
    }
}
