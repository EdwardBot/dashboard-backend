package controllers

import (
	"github.com/edward-backend/controllers/auth"
	"github.com/edward-backend/controllers/guild"
	"github.com/gin-gonic/gin"
)

func InitGuilds(r *gin.RouterGroup) {
	r.Use(auth.HasAuth).GET("/guilds/:id", guild.HandleGuild)
	r.Use(auth.HasAuth).GET("/guild/:id", guild.HandleGetGuildConfig)
	r.Use(auth.HasAuth).GET("/guild/:id/user/:uid", guild.HandleGetMember)
}
