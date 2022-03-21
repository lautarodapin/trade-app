
from sqlite3 import IntegrityError
from typing import List, Literal, Optional
from sqlalchemy import UniqueConstraint
from sqlmodel import Field, Relationship, SQLModel, Session, Column, String, select


class Symbol(SQLModel, table=True):
    id: Optional[int] = Field(None, primary_key=True)
    name: str = Field(sa_column=Column("name", String, unique=True))
    
    trades: List["Trade"] = Relationship(back_populates="symbol")


class Trade(SQLModel, table=True):
    id: Optional[int] = Field(None, primary_key=True)
    type: Literal["buy", "sell"] = Field("buy", sa_column=Column("type", String))
    quantity: float
    price: float
    earns: float

    symbol_id: Optional[int] = Field(None, foreign_key="symbol.id")
    symbol: Optional[Symbol] = Relationship(back_populates="trades")


class CreateTrade(SQLModel):
    amount: float
    symbol_id: int