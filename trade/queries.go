package trade

import (
	"fmt"
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

// Calculates the unrealized PL of the user
func getUnrealizedPL(db *gorm.DB, user models.User, closePrice float64, symbol string) float64 {
	var value float64
	var costHoldings float64
	query := `
		SELECT SUM(t.quantity) * @closePrice as value, SUM(t.quantity * t.price) as cost_holdings
		FROM trades t
		JOIN pairs p ON p.id = t.pair_id
		WHERE t.user_id = @userId AND t.type = @type AND p.symbol = @symbol
	`
	db.Debug().Model(models.Trade{}).Raw(query, map[string]interface{}{
		"closePrice": closePrice,
		"userId":     user.ID,
		"type":       models.BUY,
		"symbol":     symbol,
	}).Row().Scan(&value, &costHoldings)
	fmt.Printf("value=%+v, costHoldings=%+v\n", value, costHoldings)
	return value - costHoldings
}

// Calculates de cumulative PL of the user
func getCumulativePL(db *gorm.DB, user models.User) float64 {
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

func getNetPNL(unrealized float64, cumulative float64) float64 {
	return cumulative + unrealized
}

func getTrades(db *gorm.DB, user models.User) []models.Trade {
	var trades []models.Trade
	db.Where("user_id = ?", user.ID).Find(&trades)
	return trades
}
