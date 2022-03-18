package trade

import (
	"fmt"
	"strconv"
	"trade-app/models"

	"gorm.io/gorm"
)

type TradeResultSql struct {
	models.Trade
	Acc float64
}

// Gets all buy trades order by First in first out, until reach the desired quantity
func getBuysUntilQuantity(db *gorm.DB, user models.User, quantity float64) ([]TradeResultSql, error) {
	var buys []TradeResultSql
	err := db.Raw(`
		SELECT *
		FROM (
			SELECT id, SUM(quantity) OVER(ORDER BY id) as acc
			FROM trades
			WHERE type = @type AND user_id = @userId
		) as subquery
		JOIN trades as t ON t.id = subquery.id
		WHERE t.type = @type AND t.user_id = @userId`, // FIXME:  AND subquery.acc < @quantity for some reason it doesn't work
		map[string]interface{}{
			"userId":   user.ID,
			"type":     models.BUY,
			"quantity": quantity,
		}).Scan(&buys).Error

	return buys, err
}

// Update all the quantities of the buy trades
func updateBuysQuantityTrades(db *gorm.DB, buys []TradeResultSql) {
	for _, buy := range buys {
		db.Raw(`
			UPDATE trades
			SET quantity = @q
			WHERE id = @id
		`, map[string]interface{}{"q": buy.Quantity, "id": buy.ID}).
			Scan(&buy)
	}
}

type UnrealizedPLSqlResult struct {
	Value        float64
	CostHoldings float64
	Symbol       string
}

func queryUnrealizedPL(db *gorm.DB, user models.User) ([]UnrealizedPLSqlResult, error) {
	var results []UnrealizedPLSqlResult
	query := `
		SELECT p.symbol as symbol, SUM(t.quantity) as value, SUM(t.quantity * t.price) as cost_holdings
		FROM trades as t
		JOIN pairs p ON p.id = t.pair_id
		WHERE t.user_id = @userId AND t.type = @type
		GROUP BY p.symbol
	`
	db.Debug().Model(models.Trade{}).Raw(query, map[string]interface{}{
		"userId": user.ID,
		"type":   models.BUY,
	}).Scan(&results)
	return results, nil
}

// Calculates the unrealized PL of the user
func calculateUnrealizedPL(db *gorm.DB, user models.User) float64 {
	results, _ := queryUnrealizedPL(db, user)
	var values float64
	var costHoldings float64
	for _, result := range results {
		fmt.Printf("%+v\n", result)
		response, _ := getSymbolPrice(result.Symbol)
		price, _ := strconv.ParseFloat(response.Price, 64)
		values += result.Value * price
		costHoldings += result.CostHoldings
	}
	fmt.Printf("unrealizedPLSqlResult=%+v\n\n", results)
	fmt.Printf("values=%+v\n", values)
	return values - costHoldings
}

// Calculates de cumulative PL of the user
func calculateCumulativePL(db *gorm.DB, user models.User) float64 {
	var value float64
	query := `
		SELECT SUM(earns) as value
		FROM trades
		WHERE type = @sell AND user_id = @user_id
	`
	db.Debug().Raw(query, map[string]interface{}{
		"user_id": user.ID,
		"sell":    models.SELL,
	}).Row().Scan(&value)
	return value
}

func calculateNetPNL(unrealized float64, cumulative float64) float64 {
	return cumulative + unrealized
}

func getTrades(db *gorm.DB, user models.User) []models.Trade {
	var trades []models.Trade
	db.Debug().Preload("Pair").Where("user_id = ?", user.ID).Find(&trades)
	return trades
}
