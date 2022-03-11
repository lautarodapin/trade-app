package main

import (
	"fmt"
	"net/http"
	"strconv"
	"trade-app/models"
	"trade-app/schemas"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	pairList := []string{
		"BTCUSDT",
		"ETHUSDT",
		"BNBUSDT",
		"BCCUSDT",
		"NEOUSDT",
		"LTCUSDT",
		"QTUMUSDT",
		"ADAUSDT",
		"XRPUSDT",
		"EOSUSDT",
	}
	favPairList := pairList[:3]
	users := [3]models.User{
		{Id: 1, FirstName: "John", LastName: "Doe", Email: "user1@test.com", Token: "d73e:9666:2dec:2ed8:073f:7f52:1ffc:5b9d", Password: "password"},
		{Id: 2, FirstName: "Koko", LastName: "Doe", Email: "user2@test.com", Token: "9ae4:9c47:a59f:9427:bc36:f6ec:536f:3c83", Password: "password"},
		{Id: 3, FirstName: "Francis", LastName: "Sunday", Email: "user3@test.com", Token: "f049:fc4e:eb2a:2d50:2962:5ab7:f5c7:6b96", Password: "password"},
	}
	r.Run()
}
