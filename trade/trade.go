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
