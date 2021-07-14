package controllers

import (
	"github.com/edward-backend/controllers/auth"
	"github.com/edward-backend/controllers/guild"
	"github.com/gin-gonic/gin"
)

func InitGuilds(r *gin.RouterGroup) {
	r.Use(auth.HasAuth)
	r.GET("/guilds/:id", guild.HandleGuild)
	r.GET("/guild/:id", guild.HandleGetGuildConfig)
	r.GET("/guild/:id/user/:uid", guild.HandleGetMember)
	r.GET("/guild/:id/commands", guild.HandleCommands)
	r.DELETE("/guild/:id/commands/:name", guild.HandleDeleteCommand)
	r.POST("/guild/:id/commands/:name", guild.HandleCreateCommand)
	r.GET("/guild/:id/kicks", guild.HandleKicks)
}
