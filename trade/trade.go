package trade

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"trade-app/models"
	"trade-app/schemas"

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
func createBuyTrade(db *gorm.DB, user models.User, price float64, quantity float64, symbol string) models.Trade {
	var pair models.Pair
	db.Where("symbol = ?", symbol).First(&pair)
	trade := models.Trade{
		UserID:   user.ID,
		Type:     models.BUY,
		Quantity: quantity,
		Price:    price,
		PairID:   pair.ID,
	}
	db.Create(&trade)
	return trade
}

// Loops through all the buys trades until reach the desired quantity and calculates the earns of the sale
func makeTrade(buys []TradeResultSql, sale *models.Trade) ([]TradeResultSql, error) {
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
	if quantity > 0 {
		return buys, errors.New("not enought to sale")
	}
	fmt.Print("quantity=", quantity, " \n\n")
	fmt.Print("sale=", sale, " \n\n")

	return buys, nil
}

// Creates a sale based on the symbol and the amount requested
func makeSaleTrade(db *gorm.DB, user models.User, symbol string, amount float64) (models.Trade, error) {
	var pair models.Pair
	db.Where("symbol = ?", symbol).First(&pair)
	symbolRequest, _ := getSymbolPrice(symbol)
	price, _ := strconv.ParseFloat(symbolRequest.Price, 64)
	quantity := amount / price

	buys, _ := getBuysUntilQuantity(db, user, quantity, pair.Symbol)
	sale := models.Trade{
		UserID:   user.ID,
		Type:     models.SELL,
		Quantity: quantity,
		Price:    price,
		Earns:    0,
		PairID:   pair.ID,
	}

	buys, err := makeTrade(buys, &sale)
	if err != nil {
		return models.Trade{}, err
	}
	updateBuysQuantityTrades(db, buys)

	db.Create(&sale)

	return sale, nil
}

func MakeTradeBuyHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var postData schemas.BuyTradePost
		user := c.MustGet("user").(models.User)
		c.BindJSON(&postData)
		symbol := postData.Symbol
		amount := postData.Amount
		symbolRequest, _ := getSymbolPrice(symbol)
		price, _ := strconv.ParseFloat(symbolRequest.Price, 64)
		quantity := amount / price
		buy := createBuyTrade(db, user, price, quantity, symbol)
		c.JSON(http.StatusOK, schemas.Response{
			Status:  "success",
			Message: "Buy trade created",
			Data:    buy,
		})
	}
}

func MakeTradeSaleHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var postData schemas.BuyTradePost
		user := c.MustGet("user").(models.User)
		c.BindJSON(&postData)
		symbol := postData.Symbol
		amount := postData.Amount
		sale, err := makeSaleTrade(db, user, symbol, amount)
		if err != nil {
			c.JSON(http.StatusBadRequest, schemas.Response{
				Status:  "error",
				Message: err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, schemas.Response{
			Status:  "success",
			Message: "Sale trade created",
			Data:    sale,
		})

	}
}

func PnlHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("user").(models.User)
		unrealizedPL := calculateUnrealizedPL(db, user)
		cumulativePL := calculateCumulativePL(db, user)
		netPNL := calculateNetPNL(unrealizedPL, cumulativePL)
		c.JSON(http.StatusOK, schemas.Response{
			Status:  "success",
			Message: "PnL",
			Data: schemas.PnlResponse{
				UnrealizedPL: unrealizedPL,
				CumulativePL: cumulativePL,
				NetPNL:       netPNL,
			},
		})
	}
}

func ListTradesHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("user").(models.User)
		trades := getTrades(db, user)
		c.JSON(http.StatusOK, schemas.Response{
			Status:  "success",
			Message: "Trades",
			Data:    trades,
		})
	}
}
