package guild

import (
	"github.com/edward-backend/database"
	"github.com/gin-gonic/gin"
)

func HandleGetGuildConfig(ctx *gin.Context) {
	if !ctx.MustGet("hasAuth").(bool) {
		return
	}
	conf, err := database.FindGConf(ctx.Param("id"))
	wallet := database.FindWallet(ctx.MustGet("userId").(string), ctx.Param("id"))
	if err != nil {
		ctx.JSON(404, gin.H{
			"status":   "errorMsg",
			"errorMsg": "No guild found",
		})
		return
	}
	ctx.JSON(200, gin.H{
		"guild":  conf,
		"wallet": *wallet,
	})
}
