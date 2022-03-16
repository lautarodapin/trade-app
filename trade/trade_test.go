package trade

import (
	"fmt"
	"testing"
	"trade-app/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestMakeTrade(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:memdb1?mode=memory&cache=shared"), &gorm.Config{})

	if err != nil {
		t.Fatal(err)
	}

	db.AutoMigrate(&models.Trade{})
	db.AutoMigrate(&models.User{})
	user := models.User{
		Email: "test@test.com",
	}
	db.Create(&user)
	if user.ID == 0 {
		t.Fatal("user not created")
	}

	t.Run("Buy 12 at 100", func(t *testing.T) {
		trade := createBuyTrade(db, user, 100, 12)
		if trade.ID == 0 {
			t.Fatal("trade not created")
		}
	})

	t.Run("Buy 17 at 99", func(t *testing.T) {
		trade := createBuyTrade(db, user, 99, 17)
		if trade.ID == 0 {
			t.Fatal("trade not created")
		}
	})

	t.Run("Buy 3 at 103", func(t *testing.T) {
		trade := createBuyTrade(db, user, 103, 3)
		if trade.ID == 0 {
			t.Fatal("trade not created")
		}
	})

	t.Run("Total 32 at cost 3192", func(t *testing.T) {
		var quantity, total float64
		db.Raw(`
			SELECT SUM(quantity) as quantity, SUM(price*quantity) as total
			FROM trades
			WHERE user_id = @userId AND type = @type
		`, map[string]interface{}{"userId": user.ID, "type": models.BUY}).
			Row().Scan(&quantity, &total)

		if quantity != 32 {
			t.Errorf("Expected 32, got %f", quantity)
		}
		if total != 3192 {
			t.Errorf("Expected 3192, got %f", total)
		}
	})

	t.Run("Sell 9 at 101 make trade", func(t *testing.T) {
		sale := models.Trade{
			Type:     models.SELL,
			Quantity: 9,
			Price:    101,
			Earns:    0,
			UserID:   user.ID,
		}

		buys, _ := getBuysUntilQuantity(db, user, sale.Quantity)
		buys = makeTrade(buys, &sale)
		updateBuysQuantityTrades(db, buys)

		db.Create(&sale)
		if sale.Earns != 9*(sale.Price-100) {
			t.Errorf("Expected 9*(102-100) = %f, got %f", 9*(sale.Price-100), sale.Earns)
		}
	})

	t.Run("Total 23 at cost 2292", func(t *testing.T) {
		var quantity, total float64
		db.Raw(`
			SELECT SUM(quantity) as quantity, SUM(price*quantity) as total
			FROM trades
			WHERE user_id = @userId AND type = @type
		`, map[string]interface{}{"userId": user.ID, "type": models.BUY}).
			Row().Scan(&quantity, &total)

		if quantity != 23 {
			t.Errorf("Expected 23, got %f", quantity)
		}
		if total != 2292 {
			t.Errorf("Expected 2292, got %f", total)
		}
	})

	t.Run("Sell 4 at 105 make trade", func(t *testing.T) {
		sale := models.Trade{
			Type:     models.SELL,
			Quantity: 4,
			Price:    105,
			Earns:    0,
			UserID:   user.ID,
		}
		buys, _ := getBuysUntilQuantity(db, user, sale.Quantity)
		buys = makeTrade(buys, &sale)
		updateBuysQuantityTrades(db, buys)

		db.Create(&sale)
		if sale.Earns != 21 {
			t.Errorf("Expected 21, got %f", sale.Earns)
		}
	})

	t.Run("Total 19 at cost 1893", func(t *testing.T) {
		var quantity, total float64
		db.Raw(`
			SELECT SUM(quantity) as quantity, SUM(price*quantity) as total
			FROM trades
			WHERE user_id = @userId AND type = @type
		`, map[string]interface{}{"userId": user.ID, "type": models.BUY}).
			Row().Scan(&quantity, &total)

		if quantity != 19 {
			t.Errorf("Expected 19, got %f", quantity)
		}
		if total != 1893 {
			t.Errorf("Expected 1893, got %f", total)
		}
	})

	t.Run("Get unrealized P L with market close at 99", func(t *testing.T) {
		unrealizedPL := getUnrealizedPL(db, user, float64(99))
		fmt.Println(unrealizedPL)
		if unrealizedPL != -12 {
			t.Errorf("Expected -12, got %f", unrealizedPL)
		}
	})

	t.Run("Get cumulative realized P L", func(t *testing.T) {
		cumulativePL := getCumulativePL(db, user)
		fmt.Println(cumulativePL)
		if cumulativePL != 30 {
			t.Errorf("Expected 30, got %f", cumulativePL)
		}
	})

	t.Run("Get total earns", func(t *testing.T) {
		unrealizedPL := getUnrealizedPL(db, user, float64(99))
		cumulativePL := getCumulativePL(db, user)
		totalEarns := cumulativePL + unrealizedPL
		fmt.Println(totalEarns)
		if totalEarns != 18 {
			t.Errorf("Expected 18, got %f", totalEarns)
		}
	})

	t.Run("Get pair price for BTCUSDT", func(t *testing.T) {
		symbolRequest := getSymbolPrice("BTCUSDT")
		if symbolRequest == (SymbolRequest{}) {
			t.Errorf("Expected not empty, got %v", symbolRequest)
		}
	})
}
