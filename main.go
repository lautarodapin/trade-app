package main

import (
	"fmt"
	"net/http"
	"trade-app/models"
	"trade-app/schemas"
	"trade-app/symbols"
	"trade-app/trade"
	"trade-app/users"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {

	db, _ := gorm.Open(sqlite.Open("gorm.sqlite"), &gorm.Config{})
	models.AutoMigrate(db)
	models.InitUsers(db)
	models.InitPairList(db)
	models.InitFavPairList(db, false)
	r := gin.Default()
	r.Use(CORSMiddleware())
	public := r.Group("")
	private := r.Group("")
	private.Use(AuthMiddleware(db))
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
			c.JSON(http.StatusBadRequest, schemas.Response{Status: "error", Message: err.Error()})
			return
		}
		var user models.User
		db.Debug().Where("email = ?", userLogin.Email).First(&user)
		if user.Email == userLogin.Email && models.VerifyPassword(user.Password, userLogin.Password) {
			c.JSON(http.StatusOK, schemas.Response{Status: "success", Data: gin.H{"token": user.Token}})
			return
		}
		fmt.Println("userLogin: ", userLogin)
		fmt.Println("user: ", user)
		if user.Email == userLogin.Email && !models.VerifyPassword(user.Password, userLogin.Password) {
			c.JSON(http.StatusUnauthorized, schemas.Response{Status: "error", Message: "Invalid password"})
			return
		}
		c.JSON(http.StatusNotFound, schemas.Response{Status: "error", Message: fmt.Sprintf("User %s not found", userLogin.Email)})
	})

	users_routes := private.Group("/users")
	{
		users_routes.GET("/current", users.CurrentUserHandler(db))
		users_routes.GET("/:id", users.UserByIdHandler(db))
	}
	symbols_routes := private.Group("/pair-list")
	{
		symbols_routes.GET("/", symbols.ListSymbolsHandler(db))
		symbols_routes.GET("/fav", symbols.FavListSymbolsHandler(db))
		symbols_routes.POST("/fav", symbols.CreateFavSymbolHandler(db))
		symbols_routes.DELETE("/fav/:id", symbols.DeleteFavSymbolHandler(db))
		symbols_routes.GET("/fav/prices", symbols.PricesListFavSymbolHandler(db))
	}
	trades_routes := private.Group("trades")
	{
		trades_routes.POST("/buy", trade.MakeTradeBuyHandler(db))
		trades_routes.POST("/sale", trade.MakeTradeSaleHandler(db))
		trades_routes.GET("/:symbol", trade.PnlHandler(db))
	}
	r.Run()
}
