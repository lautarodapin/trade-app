export type ApiResponse = {
    data?: any;
    status: "success" | "error";
    message?: string
}
export type Pair = { ID: number, symbol: string }
export type FavPair = { ID: number; pair: Pair }
export enum TradeType {
    BUY = 1,
    SELL = 2,
}
export type Trade = { ID: number, type: TradeType, quantity: number, price: number, earns: number }