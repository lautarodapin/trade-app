package trade

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"trade-app/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SymbolRequest struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

// Gets the price of a symbol from binance
func getSymbolPrice(symbol string) (SymbolRequest, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.binance.com/api/v3/ticker/price?symbol=%s", symbol))
	if err != nil {
		return SymbolRequest{}, err
	}
	defer resp.Body.Close()
	var symbolRequest SymbolRequest
	json.NewDecoder(resp.Body).Decode(&symbolRequest)
	return symbolRequest, err
}

// Creates a buy trade for certain user, with some price and quantity
func createBuyTrade(db *gorm.DB, user models.User, price float64, quantity float64) models.Trade {
	trade := models.Trade{
		UserID:   user.ID,
		Type:     models.BUY,
		Quantity: quantity,
		Price:    price,
	}
	db.Create(&trade)
	return trade
}

// Loops through all the buys trades until reach the desired quantity and calculates the earns of the sale
func makeTrade(buys []TradeResultSql, sale *models.Trade) []TradeResultSql {
	quantity := sale.Quantity
	for i, buy := range buys {
		if quantity == 0 || buys[i].Quantity == 0 {
			continue
		}
		diff := buys[i].Quantity - quantity
		if diff < 0 {
			sale.Earns += (quantity - math.Abs(diff)) * (sale.Price - buy.Price)
			quantity = math.Abs(diff)
			buys[i].Quantity = 0
		} else {
			sale.Earns += quantity * (sale.Price - buy.Price)
			quantity = 0
			buys[i].Quantity = diff
		}
	}
	return buys
}

// Creates a sale based on the symbol and the amount requested
func makeSaleTrade(c *gin.Context, db *gorm.DB, symbol string, amount float64) models.Trade {
	user := c.MustGet("user").(models.User)
	symbolRequest, _ := getSymbolPrice(symbol)
	price, _ := strconv.ParseFloat(symbolRequest.Price, 64)
	quantity := amount / price

	buys, _ := getBuysUntilQuantity(db, user, quantity)

	sale := models.Trade{
		UserID:   user.ID,
		Type:     models.SELL,
		Quantity: quantity,
		Price:    price,
		Earns:    0,
	}

	makeTrade(buys, &sale)
	updateBuysQuantityTrades(db, buys)

	db.Create(&sale)

	return sale
}
