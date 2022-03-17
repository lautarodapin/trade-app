package users

import (
	"net/http"
	"strconv"
	"trade-app/models"
	"trade-app/schemas"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CurrentUserHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusNotFound, schemas.Response{Status: "error", Message: "User not found"})
			return
		}
		c.JSON(http.StatusOK, schemas.Response{Status: "success", Data: user})
	}
}

func UserByIdHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		var user models.User
		db.Debug().Where("id = ?", id).First(&user)
		c.JSON(http.StatusOK, schemas.Response{Status: "success", Data: user})
		// TODO: handle c.JSON(http.StatusNotFound, gin.H{"status":"error","message": fmt.Sprintf("User with id %d not found", id)})
	}
}
