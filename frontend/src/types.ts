export type ApiResponse = {
    data?: any;
    status: "success" | "error";
    message?: string
}
export type Pair = { ID: number, symbol: string }
export type FavPair = { ID: number; pair: Pair }