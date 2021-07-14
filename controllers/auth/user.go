package auth

import (
	"errors"
	"github.com/edward-backend/database"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
)

func HandleUserInfo(c *gin.Context) {
	if !c.MustGet("hasAuth").(bool) {
		return
	}
	uIdS := c.MustGet("userId").(string)
	uId, _ := strconv.ParseInt(uIdS, 10, 64)
	var user database.User
	r := database.Conn.First(&user, uId)
	if errors.Is(r.Error, gorm.ErrRecordNotFound) {
		c.JSON(404, gin.H{
			"status": "error",
			"error":  "no user found",
		})
		return
	}
	user.UID = strconv.FormatInt(user.UserID, 10)
	c.JSON(200, user)
}
