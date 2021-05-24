package guild

import (
	"github.com/edward-backend/database"
	"github.com/edward-backend/utils"
	"github.com/gin-gonic/gin"
)

func HandleGetGuildConfig(ctx *gin.Context) {
	conf, err := database.FindGConf(ctx.Param("id"))
	tokenData := utils.ParseToken(ctx)
	wallet := database.FindWallet(tokenData["sub"].(string), ctx.Param("id"))
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
