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

	users_route := r.Group("/users")
	{
		users_route.POST("/login", func(c *gin.Context) {
			var userLogin schemas.UserLogin
			err := c.BindJSON(&userLogin)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
				return
			}
			for _, user := range users {
				if user.Email == userLogin.Email && user.Password == userLogin.Password {
					c.JSON(http.StatusOK, gin.H{"token": user.Token})
					return
				} else if user.Email == userLogin.Email && user.Password != userLogin.Password {
					c.JSON(http.StatusUnauthorized, gin.H{"detail": "Invalid password"})
					return
				}
			}
			c.JSON(http.StatusNotFound, gin.H{"detail": fmt.Sprintf("User %s not found", userLogin.Email)})
		})
		users_route.GET("/current", func(c *gin.Context) {
			var authHeader schemas.AuthHeader
			err := c.BindHeader(&authHeader)
			if err != nil {
				c.JSON(http.StatusForbidden, gin.H{"detail": "No Authorization header"})
				return
			}
			for _, user := range users {
				if fmt.Sprintf("Bearer %s", user.Token) == authHeader.Authorization {
					c.JSON(http.StatusOK, gin.H{"user": user})
					return
				}
			}
			c.JSON(http.StatusUnauthorized, gin.H{"detail": "Invalid token"})
		})
		users_route.GET("/:id", func(c *gin.Context) {
			id, _ := strconv.Atoi(c.Param("id"))
			for _, user := range users {
				if user.Id == id {
					c.JSON(http.StatusOK, user)
					return
				}
			}
			c.JSON(http.StatusNotFound, gin.H{"detail": fmt.Sprintf("User with id %d not found", id)})
		})
	}
	r.Run()
}
