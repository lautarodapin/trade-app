from sqlmodel import Session, select, alias, func
from database import Symbol, Trade


def create_buy_trade(price: float, quantity: float, symbol_id: str, session:Session):
    trade = Trade(
        symbol_id=symbol_id,
        earns=0,
        price=price,
        quantity=quantity,
        type="buy",
    )
    session.add(trade)
    session.commit()
    session.refresh(trade)
    return trade

def create_sale_trade(
    price: float,
    quantity: float,
    symbol_id: int,
    session: Session
):
    sale = Trade(
        symbol_id=symbol_id,
        type="sell",
        quantity=quantity,
        price=price,
        earns=0,
    )
    trades = trades_until_quantity(quantity, session)
    trades, sale = make_trade(trades, sale)

    return sale
    
    
def trades_until_quantity(
    quantity_required: float,
    session: Session,
) -> list[Trade]:
    subquery = (
        select(Trade.id, func.sum(Trade.quantity).over(order_by=Trade.id).label("ac"))
        .where(Trade.type == "buy")
        .subquery()
    )
    # sq=alias(Trade, subquery)

    query = (
        # select(Trade, subquery.c.ac)
        select(Trade)
        .where(Trade.type == "buy")
        # .where(subquery.c.ac <= quantity_required)
        .where(subquery.c.ac.between(quantity_required, quantity_required))
        # .where(Trade.id == subquery.c.id)
    )
    results = session.exec(query).fetchall()
    print(results)
    return results


def make_trade(
    trades: list[Trade],
    sale: Trade,
): 
    quantity = sale.quantity
    for trade in trades:
        if quantity == 0 or trade.quantity == 0:
            continue
        diff = trade.quantity - quantity
        if diff < 0:
            sale.earns += (quantity - abs(diff)) * (sale.price - trade.price)
            quantity = abs(diff)
            trade.quantity = 0
            continue
        sale.earns += quantity * (sale.price - trade.price)
        quantity = 0
        trade.quantity = diff
        break
    return trades, sale
