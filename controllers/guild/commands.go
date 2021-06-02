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
	if !c.MustGet("hasAuth").(bool) {
		return
	}
	doc, err := database.FindCommand(c.Param("id"), c.Param("name"))
	if err != nil {
		c.JSON(500, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}
	err = doc.Delete()
	if err != nil {
		c.JSON(500, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"status":  "success",
		"message": "Successfully deleted the command",
	})
}

func HandleCreateCommand(c *gin.Context) {
	if !c.MustGet("hasAuth").(bool) {
		return
	}

	cmd := database.CustomCommand{
		GuildId:  c.Param("id"),
		Name:     c.Param("name"),
		Response: c.MustGet("body").(map[string]interface{})["response"].(string),
	}
	err := cmd.Save()
	if err != nil {
		c.JSON(500, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
	}
}
