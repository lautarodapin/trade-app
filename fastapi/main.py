from fastapi import Depends, FastAPI, Query
from sqlmodel import SQLModel, Session, create_engine, select
from services import BinanceService,  get_binance_service, create_buy_trade
from database import Symbol ,init_symbols ,Trade, CreateTrade

app = FastAPI(debug=True)

sqlite_file_name = "database.sqlite"
sqlite_url = f"sqlite:///{sqlite_file_name}"

engine = create_engine(sqlite_url, echo=True)
SQLModel.metadata.create_all(engine)


def get_session():
    with Session(engine) as session:
        yield session


@app.on_event("startup")
def startup():
    init_symbols(engine)

@app.get("/symbols")
def symbols_list(session: Session = Depends(get_session)):
    return session.exec(select(Symbol).all())


@app.get("/symbols/{symbol_name}/price", response_model=BinanceService.Response)
def get_symbol_price(symbol_name: str, service: BinanceService= Depends(get_binance_service)):
    return service.get_symbol_price(symbol_name)[0]

@app.get("symbols/price", response_model=list[BinanceService.Response])
def get_symbols_prices(symbols: list[str] = Query(None), service: BinanceService= Depends(get_binance_service)):
    return service.get_symbols_prices(symbols)[0]



@app.post("/symbols/buy", response_model=Trade)
def buy_symbol(
    data: CreateTrade,
    service: BinanceService= Depends(get_binance_service),
    session: Session = Depends(get_session),
):
    id, name = session.exec(select(Symbol.i, Symbol.name).where(Symbol.id == data.symbol_id)).first()
    resp, _ = service.get_symbol_price(name)
    quantity = data.amount / float(resp.price)
    trade = create_buy_trade(
        price=resp.price,
        quantity=quantity,
        symbol_id=id,
        session=session,
    )

    return trade

@app.post("/symbols/sell", response_model=Trade)
def sell_symbol(
    data: CreateTrade,
    service: BinanceService= Depends(get_binance_service),
    session: Session = Depends(get_session),
):
    id, name = session.exec(select(Symbol.i, Symbol.name).where(Symbol.id == data.symbol_id)).first()
    resp, _ = service.get_symbol_price(name)
    quantity = data.amount / float(resp.price)
    trade = create_buy_trade(
        price=resp.price,
        quantity=quantity,
        symbol_id=id,
        session=session,
    )

    return trade