package guild

import (
	"github.com/edward-backend/database"
	"github.com/gin-gonic/gin"
	"strconv"
)

var (
	errorMsg = gin.H{
		"status": "errorMsg",
		"errorMsg": "Invalid id provided",
	}
)

func HandleGuild(c *gin.Context)()  {
	i, e := strconv.ParseInt(c.Param("id"), 10, 64)
	if e != nil {
		c.JSON(500, errorMsg)
		return
	}
	g, _ := database.FindGuilds(i)
	c.JSON(200, g)
}