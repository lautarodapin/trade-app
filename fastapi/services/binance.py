from decimal import Decimal
from typing import Tuple
from pydantic import BaseModel
import requests


class BinanceService:
    class Response(BaseModel):
        symbol: str
        price: Decimal

    def get_symbol_price(self, symbol: str) -> Tuple[Response, int]:
        resp = requests.get(f"https://api.binance.com/api/v3/ticker/price?symbol={symbol}")
        if resp.status_code != 200:
            return None, resp.status_code
        symbol_request = resp.json()
        return self.Response(**symbol_request), resp.status_code

    def get_symbols_prices(self, symbols: list[str]) -> Tuple[list[Response], int]:
        resp = requests.get("https://api.binance.com/api/v3/ticker/price?symbols=[{}]".format(",".join([f"\"{s}\"" for s in symbols])))
        if resp.status_code != 200:
            return None, resp.status_code
        data = resp.json()
        return [self.Response(**d) for d in data], resp.status_code


def get_binance_service():
    yield BinanceService()