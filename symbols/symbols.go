package symbols

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"trade-app/models"
	"trade-app/schemas"
	"trade-app/trade"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ListSymbolsHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var pairList []models.Pair
		db.Debug().Find(&pairList)
		c.JSON(http.StatusOK, schemas.Response{Status: "success", Data: pairList})
	}
}

func FavListSymbolsHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("user").(models.User)
		var favPairList []models.FavPair
		db.Debug().Preload("User").Preload("Pair").Where("user_id = ?", user.ID).Find(&favPairList)

		fmt.Println("favPairList: ", favPairList)
		c.JSON(http.StatusOK, schemas.Response{Status: "success", Data: favPairList})
	}
}

func CreateFavSymbolHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var jsonData schemas.Symbol
		user := c.MustGet("user").(models.User)
		c.BindJSON(&jsonData)
		symbol := jsonData.Symbol
		var exists int64
		db.Debug().
			Model(&models.FavPair{}).
			Joins("JOIN users ON users.id = user_id").
			Joins("JOIN pairs ON pairs.id = pair_id").
			Where("pairs.symbol = ?", symbol).
			Where("user_id = ?", user.ID).
			Count(&exists)
		if exists > 0 {
			c.JSON(http.StatusBadRequest, schemas.Response{Status: "error", Message: fmt.Sprintf("%s already exists", symbol)})
			return
		}
		var count int64
		db.Debug().
			Model(&models.FavPair{}).
			Where("user_id = ?", user.ID).
			Count(&count)
		if count >= 3 {
			c.JSON(http.StatusBadRequest, schemas.Response{Status: "error", Message: "You can't add more than 3 symbols"})
			return
		}
		var pair models.Pair
		db.Debug().Where("symbol = ?", symbol).First(&pair)
		favPair := models.FavPair{UserID: user.ID, PairID: pair.ID}
		db.Debug().Create(&favPair)
		c.JSON(http.StatusCreated, schemas.Response{Status: "success", Message: fmt.Sprintf("%s added to fav list", symbol), Data: favPair})
	}
}

func DeleteFavSymbolHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		user := c.MustGet("user").(models.User)
		var favPair models.FavPair
		db.Debug().Where("user_id = ?", user.ID).Where("id = ?", id).First(&favPair)
		if favPair.ID == 0 {
			c.JSON(http.StatusNotFound, schemas.Response{Status: "error", Message: fmt.Sprintf("Fav pair with id %d not found", id)})
			return
		}
		db.Debug().Unscoped().Delete(&favPair) // unscoped to ignore softs deletes
		c.JSON(http.StatusOK, schemas.Response{Status: "success", Message: fmt.Sprintf("%d deleted from fav", favPair.ID)})
	}
}

func PricesListFavSymbolHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var symbolRequestList []trade.SymbolRequest
		user := c.MustGet("user").(models.User)
		var favPairList []models.Pair
		db.Debug().
			Joins("JOIN fav_pairs ON fav_pairs.pair_id = pairs.id").
			Where("fav_pairs.user_id = ?", user.ID).
			Find(&favPairList)

		for _, pair := range favPairList {
			resp, err := http.Get(fmt.Sprintf("https://api.binance.com/api/v3/ticker/price?symbol=%s", pair.Symbol))
			if err != nil {
				c.JSON(http.StatusInternalServerError, schemas.Response{Status: "error", Message: fmt.Sprintf("Error fetching symbol %s, got %s", pair.Symbol, err.Error())})
				return
			}
			defer resp.Body.Close()
			var symbolRequest trade.SymbolRequest
			json.NewDecoder(resp.Body).Decode(&symbolRequest)
			symbolRequestList = append(symbolRequestList, symbolRequest)
		}
		c.JSON(http.StatusOK, schemas.Response{Status: "success", Data: symbolRequestList})
	}
}
