package guild

import (
	"github.com/edward-backend/utils"
	"github.com/gin-gonic/gin"
)

func HandleGetMember(ctx *gin.Context) {
	if !ctx.MustGet("hasAuth").(bool) {
		return
	}
	member, err := utils.GetDiscordInstance().GuildMember(ctx.Param("id"), ctx.Param("uid"))
	if err != nil {
		ctx.JSON(500, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}
	ctx.JSON(200, member)
}
