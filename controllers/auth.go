package controllers

import (
	"github.com/edward-backend/controllers/auth"
	"github.com/gin-gonic/gin"
)

func InitAuth(router *gin.RouterGroup) {
	router.GET("/oauth", auth.HandleOAuth)
	router.Use(auth.HasAuth).GET("/user", auth.HandleUserInfo)
	router.Use(auth.HasAuth).POST(`/refresh`, auth.HandleRefresh)
}
