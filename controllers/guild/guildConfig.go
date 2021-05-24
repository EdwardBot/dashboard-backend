package guild

import (
	"github.com/edward-backend/database"
	"github.com/gin-gonic/gin"
)

func HandleGetGuildConfig(ctx *gin.Context) {
	conf, err := database.FindGConf(ctx.Param("id"))
	if err != nil {
		ctx.JSON(404, gin.H{
			"status": "errorMsg",
			"errorMsg": "No guild found",
		})
		return
	}
	ctx.JSON(200, conf)
}
