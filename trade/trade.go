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

func getPairPrice(symbol string) SymbolRequest {
	resp, err := http.Get(fmt.Sprintf("https://api.binance.com/api/v3/ticker/price?symbol=%s", symbol))
	if err != nil {
		return SymbolRequest{}
	}
	defer resp.Body.Close()
	var symbolRequest SymbolRequest
	json.NewDecoder(resp.Body).Decode(&symbolRequest)
	return symbolRequest
}

func makeBuy(db *gorm.DB, user models.User, price float64, quantity float64) models.Trade {
	trade := models.Trade{
		UserID:   user.ID,
		Type:     models.BUY,
		Quantity: quantity,
		Price:    price,
	}
	db.Create(&trade)
	return trade
}
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

func makeSale(c *gin.Context, db *gorm.DB, symbol string, amount float64) models.Trade {
	user := c.MustGet("user").(models.User)
	symbolRequest := getPairPrice(symbol)
	price, _ := strconv.ParseFloat(symbolRequest.Price, 64)
	quantity := amount / price

	buys, _ := getBuys(db, user, quantity)

	sale := models.Trade{
		UserID:   user.ID,
		Type:     models.SELL,
		Quantity: quantity,
		Price:    price,
		Earns:    0,
	}

	makeTrade(buys, &sale)
	updateBuys(db, buys)

	return sale
}
