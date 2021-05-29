package auth

import (
	"github.com/edward-backend/database"
	"github.com/gin-gonic/gin"
	"strconv"
)

func HandleUserInfo(c *gin.Context) {
	if !c.MustGet("hasAuth").(bool) {
		return
	}
	uIdS := c.MustGet("userId").(string)
	uId, _ := strconv.ParseInt(uIdS, 10, 64)
	user, e := database.FindUser(uId)
	if e != nil {
		c.JSON(404, gin.H{
			"status": "error",
			"error":  "no user found",
		})
		return
	}
	c.JSON(200, user)
}
