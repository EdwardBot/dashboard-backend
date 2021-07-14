package guild

import (
	"errors"
	"github.com/edward-backend/database"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
)

var (
	errorMsg = gin.H{
		"status":   "errorMsg",
		"errorMsg": "Invalid id provided",
	}
	roles = gin.H{}
)

func HandleGuild(c *gin.Context) {
	if !c.MustGet("hasAuth").(bool) {
		return
	}
	i, e := strconv.ParseInt(c.Param("id"), 10, 64)
	if e != nil {
		c.JSON(500, errorMsg)
		return
	}
	var user database.User
	r := database.Conn.First(&user, i)
	if errors.Is(r.Error, gorm.ErrRecordNotFound) {
		c.JSON(404, errorMsg)
		return
	}

	var guilds []database.Guild
	database.Conn.Model(&database.Guild{}).Where("gid = any(?::bigint[])", user.Guilds).Find(&guilds)

	for g := range guilds {
		guilds[g].ID = strconv.FormatInt(int64(guilds[g].GuildID), 10)
	}

	c.JSON(200, guilds)
}
