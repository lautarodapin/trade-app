

from sqlmodel import select, Session
from .models import Symbol

def init_symbols(engine):
    with Session(engine) as session:
        symbols = [
            "BTCUSDT",
            "ETHUSDT",
            "BNBUSDT",
            "BCCUSDT",
            "NEOUSDT",
            "LTCUSDT",
            "QTUMUSDT",
            "ADAUSDT",
            "XRPUSDT",
            "EOSUSDT",
        ]
        for symbol in symbols:
            query = select(Symbol.id).where(Symbol.name == symbol)
            exists = session.exec(query).first()
            if exists:
                continue
            s = Symbol(name=symbol)
            session.add(s)
        session.commit()