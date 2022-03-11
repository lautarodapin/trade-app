package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"trade-app/models"
	"trade-app/schemas"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware(users []models.User) gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")

		token := c.Request.Header.Get("Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		var _user models.User
		for _, u := range users {
			if u.Token == strings.Split(token, "Bearer ")[1] {
				_user = u
				break
			}
		}
		fmt.Printf("Usuario: %s\n", _user.Email)
		c.Set("user", _user)

		c.Next()
	}
}

func readFile(path string) ([][]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return [][]string{}, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return [][]string{}, err
	}

	return records, nil
}

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	var records [][]string
	records, _ = readFile("./pair-list.csv")
	pairList := []string{}
	for _, v := range records {
		pairList = append(pairList, v[0])
	}
	records, _ = readFile("./users.csv")
	users := []models.User{}
	for _, v := range records {
		id, _ := strconv.Atoi(v[0])
		user := models.User{
			Id:        id,
			FirstName: v[1],
			LastName:  v[2],
			Email:     v[3],
			Token:     v[4],
			Password:  v[5],
		}
		users = append(users, user)
	}
	favPairList := pairList[:3]
	r.Use(CORSMiddleware(users))

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
	pair_list_route := r.Group("/pair-list")
	{
		pair_list_route.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, pairList)
		})

		pair_list_route.GET("/fav", func(c *gin.Context) {
			c.JSON(http.StatusOK, favPairList)
		})
		pair_list_route.POST("/fav", func(c *gin.Context) {
			var jsonData schemas.Symbol
			c.BindJSON(&jsonData)
			symbol := jsonData.Symbol
			for _, favPair := range favPairList {
				if favPair == symbol {
					c.JSON(http.StatusConflict, gin.H{"detail": "Symbol already exists"})
					return
				}
			}
			favPairList = append(favPairList, symbol)
			if len(favPairList) > 3 {
				favPairList = favPairList[1:]
			}
			c.JSON(http.StatusCreated, favPairList)
		})
		pair_list_route.DELETE("/fav/:index", func(c *gin.Context) {
			index, _ := strconv.Atoi(c.Param("index"))
			if index > len(favPairList) {
				c.JSON(http.StatusNotFound, gin.H{"detail": "Index not found"})
			} else {
				favPairList = append(favPairList[:index], favPairList[index+1:]...)
				c.JSON(http.StatusOK, favPairList)
			}
		})
		pair_list_route.GET("/fav/prices", func(c *gin.Context) {
			var symbolRequestList []schemas.SymbolRequest
			for _, favSymbol := range favPairList {
				resp, err := http.Get(fmt.Sprintf("https://api.binance.com/api/v3/ticker/price?symbol=%s", favSymbol))
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"detail": fmt.Sprintf("Error fetching symbol %s, got %s", favSymbol, err.Error())})
					return
				}
				defer resp.Body.Close()
				var symbolRequest schemas.SymbolRequest
				json.NewDecoder(resp.Body).Decode(&symbolRequest)
				symbolRequestList = append(symbolRequestList, symbolRequest)
			}
			c.JSON(http.StatusOK, symbolRequestList)
		})
	}
	r.Run()
}
