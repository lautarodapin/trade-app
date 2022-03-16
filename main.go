package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"trade-app/models"
	"trade-app/schemas"
	"trade-app/trade"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func CORSMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")

		authHeader := c.Request.Header.Get("Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		fmt.Printf("authHeader: %s\n", authHeader)
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"detail": "unauthorized"})
			return
		}
		var _user models.User
		token := strings.TrimPrefix(authHeader, "Bearer ")
		db.Where("token = ?", token).First(&_user)
		fmt.Printf("token: %s\n", token)
		fmt.Printf("_user: %+v\n", _user)
		c.Set("user", _user)
		c.Next()
	}
}

func main() {

	db, _ := gorm.Open(sqlite.Open("gorm.sqlite"), &gorm.Config{})
	models.AutoMigrate(db)
	models.InitUsers(db)
	models.InitPairList(db)
	models.InitFavPairList(db, false)
	r := gin.Default()
	public := r.Group("")
	private := r.Group("")
	private.Use(CORSMiddleware(db))
	public.GET("/ping", func(c *gin.Context) {
		user, exists := c.Get("user")
		if exists {
			fmt.Printf("user: %+v\n", user)
			c.JSON(200, gin.H{
				"message": "pong",
				"user":    user,
			})
			return
		}
		c.JSON(200, gin.H{
			"message": "ping",
		})
	})
	public.POST("/login", func(c *gin.Context) {
		var userLogin schemas.UserLogin
		fmt.Printf("userLogin: %+v\n", userLogin)
		err := c.BindJSON(&userLogin)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
			return
		}
		var user models.User
		db.Where("email = ?", userLogin.Email).First(&user)
		if user.Email == userLogin.Email && user.Password == userLogin.Password {
			c.JSON(http.StatusOK, gin.H{"token": user.Token})
			return
		} else if user.Email == userLogin.Email && user.Password != userLogin.Password {
			c.JSON(http.StatusUnauthorized, gin.H{"detail": "Invalid password"})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"detail": fmt.Sprintf("User %s not found", userLogin.Email)})
	})

	users_route := private.Group("/users")
	{
		users_route.GET("/current", func(c *gin.Context) {
			user, exists := c.Get("user")
			if !exists {
				c.JSON(http.StatusNotFound, gin.H{"detail": "User not found"})
				return
			}
			c.JSON(http.StatusOK, user)
		})
		users_route.GET("/:id", func(c *gin.Context) {
			id, _ := strconv.Atoi(c.Param("id"))
			var user models.User
			db.Where("id = ?", id).First(&user)
			c.JSON(http.StatusOK, user)
			// TODO: handle c.JSON(http.StatusNotFound, gin.H{"detail": fmt.Sprintf("User with id %d not found", id)})
		})
	}
	pair_list_route := private.Group("/pair-list")
	{
		pair_list_route.GET("/", func(c *gin.Context) {
			var pairList []models.Pair
			db.Find(&pairList)
			c.JSON(http.StatusOK, pairList)
		})
		pair_list_route.GET("/fav", func(c *gin.Context) {
			user := c.MustGet("user").(models.User)
			var favPairList []models.FavPair
			db.
				Joins("JOIN users ON users.id = user_id").
				Joins("JOIN pairs ON pairs.id = pair_id").
				Where("user_id = ?", user.ID).
				Find(&favPairList)
			c.JSON(http.StatusOK, favPairList)
		})
		pair_list_route.POST("/fav", func(c *gin.Context) {
			var jsonData schemas.Symbol
			user := c.MustGet("user").(models.User)
			c.BindJSON(&jsonData)
			symbol := jsonData.Symbol
			var exists int64
			db.
				Model(&models.FavPair{}).
				Joins("JOIN users ON users.id = user_id").
				Joins("JOIN pairs ON pairs.id = pair_id").
				Where("pairs.symbol = ?", symbol).
				Where("user_id = ?", user.ID).
				Count(&exists)
			if exists > 0 {
				c.JSON(http.StatusBadRequest, gin.H{"detail": fmt.Sprintf("%s already exists", symbol)})
				return
			}
			var count int64
			db.
				Model(&models.FavPair{}).
				Where("user_id = ?", user.ID).
				Count(&count)
			if count >= 3 {
				c.JSON(http.StatusBadRequest, gin.H{"detail": "You can't add more than 3 symbols"})
				return
			}
			var pair models.Pair
			db.Where("symbol = ?", symbol).First(&pair)
			favPair := models.FavPair{UserID: user.ID, PairID: pair.ID}
			db.Create(&favPair)
			c.JSON(http.StatusCreated, favPair)
		})
		pair_list_route.DELETE("/fav/:id", func(c *gin.Context) {
			id, _ := strconv.Atoi(c.Param("id"))
			user := c.MustGet("user").(models.User)
			var favPair models.FavPair
			db.Where("user_id = ?", user.ID).Where("id = ?", id).First(&favPair)
			if favPair.ID == 0 {
				c.JSON(http.StatusNotFound, gin.H{"detail": fmt.Sprintf("Fav pair with id %d not found", id)})
				return
			}
			db.Unscoped().Delete(&favPair) // unscoped to ignore softs deletes
			c.JSON(http.StatusOK, gin.H{"detail": fmt.Sprintf("%d deleted from fav", favPair.ID)})
		})
		pair_list_route.GET("/fav/prices", func(c *gin.Context) {
			var symbolRequestList []trade.SymbolRequest
			user := c.MustGet("user").(models.User)
			var favPairList []models.Pair
			db.
				Joins("JOIN fav_pairs ON fav_pairs.pair_id = pairs.id").
				Where("fav_pairs.user_id = ?", user.ID).
				Find(&favPairList)

			for _, pair := range favPairList {
				resp, err := http.Get(fmt.Sprintf("https://api.binance.com/api/v3/ticker/price?symbol=%s", pair.Symbol))
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"detail": fmt.Sprintf("Error fetching symbol %s, got %s", pair.Symbol, err.Error())})
					return
				}
				defer resp.Body.Close()
				var symbolRequest trade.SymbolRequest
				json.NewDecoder(resp.Body).Decode(&symbolRequest)
				symbolRequestList = append(symbolRequestList, symbolRequest)
			}
			c.JSON(http.StatusOK, symbolRequestList)
		})
		trades_routes := private.Group("trades")
		trades_routes.POST("/buy", trade.MakeTradeBuyHandler(db))
		trades_routes.POST("/sale", trade.MakeTradeSaleHandler(db))
		trades_routes.GET("/:symbol", trade.PnlHandler(db))
	}
	r.Run()
}
