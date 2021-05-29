package guild

import (
	"github.com/edward-backend/database"
	"github.com/gin-gonic/gin"
)

func HandleCommands(c *gin.Context) {
	if !c.MustGet("hasAuth").(bool) {
		return
	}
	c.JSON(200, database.FindCommands(c.Param("id")))
}

func HandleDeleteCommand(c *gin.Context) {

}

func HandleCreateCommand(c *gin.Context) {

}
